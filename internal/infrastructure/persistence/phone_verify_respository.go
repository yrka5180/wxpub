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
	log "github.com/sirupsen/logrus"
)

func NewPhoneVerifyRepo(smsClient smsPb.SenderClient) *PhoneVerifyRepo {
	return &PhoneVerifyRepo{
		smsGRPCClient: smsClient,
	}
}

type PhoneVerifyRepo struct {
	smsGRPCClient smsPb.SenderClient
}

func (r *PhoneVerifyRepo) SendSms(ctx context.Context, content string, sender string, phone string) (err error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SendSms traceID:%s", traceID)

	// 发短信是调用的第三方的服务，计费使用
	_, err = r.smsGRPCClient.SendMessage(ctx, &smsPb.SendMsgRequest{
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

func (r *PhoneVerifyRepo) SetVerifyCodeSmsStorage(ctx context.Context, challenge string, verifyCodeID string, verifyCodeAnswer string) (err error) {
	var verifyCodeSmsRedisValue entity.VerifyCodeRedisValue

	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SetVerifyCodeSmsStorage traceID:%s", traceID)

	smsCreateTime := time.Now().UnixNano()
	verifyCodeSmsRedisValue.VerifyCodeCreateTime = smsCreateTime
	verifyCodeSmsRedisValue.VerifyCodeAnswer = verifyCodeAnswer

	smsRedisValue, _ := json.Marshal(verifyCodeSmsRedisValue)

	// 使用txpipeline进行原子操作
	pipe := redis.RClient.TxPipeline()
	// redis存放verifyCodeID:{verifyCodeAnswer,smsCreateTime}到相应的challenge的hashset上
	err = pipe.HSet(consts.RedisKeyPrefixChallenge+challenge, consts.RedisKeyPrefixSms+verifyCodeID, smsRedisValue).Err()
	if err != nil {
		log.Errorf("failed to do redis HSet, error: %+v, traceID: %s", err, traceID)
		return
	}

	// 更新当前challenge的key过期时间为30分钟，30分钟不执行验证短信验证码操作就无法绑定手机号，即为30分钟内可以重发短信验证码
	err = pipe.Expire(consts.RedisKeyPrefixChallenge+challenge, consts.VerifyCodeSmsChallengeTTL).Err()
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
