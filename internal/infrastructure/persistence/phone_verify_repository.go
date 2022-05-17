package persistence

import (
	"context"
	"encoding/json"
	"time"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/pkg/redis"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	"git.nova.net.cn/nova/go-common/uuid"
	smsPb "git.nova.net.cn/nova/notify/sms-xuanwu/pkg/grpcIFace"
	captchaPb "git.nova.net.cn/nova/shared/captcha/pkg/grpcIFace"
	log "github.com/sirupsen/logrus"
)

type PhoneVerifyRepo struct {
	smsGRPCClient    smsPb.SenderClient
	captchaRPCClient captchaPb.CaptchaServiceClient
}

var defaultPhoneVerifyRepo *PhoneVerifyRepo

func NewPhoneVerifyRepo() {
	if defaultPhoneVerifyRepo == nil {
		defaultPhoneVerifyRepo = &PhoneVerifyRepo{
			smsGRPCClient:    CommonRepositories.SmsGRPCClient,
			captchaRPCClient: CommonRepositories.CaptchaGRPCClient,
		}
	}
}

func DefaultPhoneVerifyRepo() *PhoneVerifyRepo {
	return defaultPhoneVerifyRepo
}

func (r *PhoneVerifyRepo) GenCaptcha(ctx context.Context, width int32, height int32) (captchaID, captchaBase64Value string, err error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("GenCaptcha traceID:%s", traceID)

	c := utils.ToOutGoingContext(ctx)
	rpcResp, err := r.captchaRPCClient.Get(c, &captchaPb.GetCaptchaRequest{
		Width:           width,
		Height:          height,
		NoiseCount:      10,
		ShowLineOptions: 2,
	})
	if err != nil {
		log.Errorf("GenCaptcha get captcha error: %v, traceID: %s", err, traceID)
		return
	}

	captchaID = rpcResp.GetID()
	captchaBase64Value = rpcResp.GetBase64Value()
	return
}

func (r *PhoneVerifyRepo) VerifyCaptcha(ctx context.Context, captchaID string, captchaAnswer string) (ok bool, err error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("VerifyCaptcha traceID:%s", traceID)

	c := utils.ToOutGoingContext(ctx)
	rpcResp, err := r.captchaRPCClient.Verify(c, &captchaPb.VerifyCaptchaRequest{
		ID:     captchaID,
		Answer: captchaAnswer,
	})
	if err != nil {
		log.Errorf("VerifyCaptcha Verify error: %v, traceID: %s", err, traceID)
		return
	}

	ok = rpcResp.GetData()
	return
}

func (r *PhoneVerifyRepo) SendSms(ctx context.Context, content string, sender string, phone string) (err error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SendSms traceID:%s", traceID)

	c := utils.ToOutGoingContext(ctx)
	// 发短信是调用的第三方的服务，计费使用
	_, err = r.smsGRPCClient.SendMessage(c, &smsPb.SendMsgRequest{
		Content: content,
		Sender:  sender,
		Items: []*smsPb.SendMsgRequest_Item{
			{
				To:        phone,
				MessageID: uuid.Get(), // 不需要查询，可以忽略
			},
		},
	})
	if err != nil {
		log.Errorf("send sms message error: %+v, traceID: %s", err, traceID)
	}

	return
}

func (r *PhoneVerifyRepo) SetVerifyCodeSmsStorage(ctx context.Context, openID string, verifyCodeID string, verifyCodeAnswer string) (err error) {
	var verifyCodeSmsRedisValue entity.VerifyCodeRedisValue

	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SetVerifyCodeSmsStorage traceID:%s", traceID)

	smsCreateTime := time.Now().UnixNano()
	verifyCodeSmsRedisValue.VerifyCodeCreateTime = smsCreateTime
	verifyCodeSmsRedisValue.VerifyCodeAnswer = verifyCodeAnswer

	smsRedisValue, _ := json.Marshal(verifyCodeSmsRedisValue)

	// 使用pipeline进行原子操作
	pipe := redis.RClient.TxPipeline()
	// redis存放verifyCodeID:{verifyCodeAnswer,smsCreateTime}到相应的challenge的hashset上
	err = pipe.HSet(consts.RedisKeyPrefixChallenge+openID, consts.RedisKeyPrefixSms+verifyCodeID, smsRedisValue).Err()
	if err != nil {
		log.Errorf("failed to do redis HSet, error: %+v, traceID: %s", err, traceID)
		return
	}

	// 更新当前challenge的key过期时间为30分钟，30分钟不执行验证短信验证码操作就无法绑定手机号，即为30分钟内可以重发短信验证码
	err = pipe.Expire(consts.RedisKeyPrefixChallenge+openID, consts.VerifyCodeSmsChallengeTTL).Err()
	if err != nil {
		log.Errorf("failed to do redis SetExpireTime, error: %+v, traceID: %s", err, traceID)
		return
	}

	_, err = pipe.Exec()
	if err != nil {
		log.Errorf("failed to exec redis pipeline, error: %+v, traceID: %s", err, traceID)
	}

	return
}

func (r *PhoneVerifyRepo) VerifySmsCode(ctx context.Context, openID string, verifyCodeID, verifyCodeAnswer string, ttl int64) (ok, isExpire bool, err error) {
	var value []byte
	var verifyCodeValue entity.VerifyCodeRedisValue
	ok = false
	isExpire = false
	now := time.Now().UnixNano()

	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("VerifySmsCode traceID:%s", traceID)

	value, err = redis.RClient.HGet(consts.RedisKeyPrefixChallenge+openID, consts.RedisKeyPrefixSms+verifyCodeID).Bytes()
	if err != nil {
		log.Errorf("failed to do redis HGet, error: %+v, traceID: %s", err, traceID)
		return
	}

	err = json.Unmarshal(value, &verifyCodeValue)
	if err != nil {
		log.Errorf("VerifySmsCode json unmarshal failed, error: %v, traceID: %s", err, traceID)
		return
	}

	// 是否过期
	if (now-verifyCodeValue.VerifyCodeCreateTime)/1e9 > ttl {
		isExpire = true
	}

	// 检查验证码
	if verifyCodeValue.VerifyCodeAnswer == verifyCodeAnswer {
		ok = true
	}

	return
}
