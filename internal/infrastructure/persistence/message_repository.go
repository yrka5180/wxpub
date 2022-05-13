package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/pkg/kafka"
	"github.com/jinzhu/gorm"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/config"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/errors"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/httputil"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	log "github.com/sirupsen/logrus"
)

type MessageRepo struct {
	kafkaTopics       []string
	currentTopicIndex int
	DB                *gorm.DB
	MQ                *kafka.MQ
}

var defaultMessageRepo *MessageRepo

func NewMessageRepo(topics []string) {
	if defaultMessageRepo == nil {
		defaultMessageRepo = &MessageRepo{
			kafkaTopics: topics,
			DB:          CommonRepositories.DB,
			MQ:          CommonRepositories.MQ,
		}
	}
}

func DefaultMessageRepo() *MessageRepo {
	return defaultMessageRepo
}

// GetTopic 获取Topic
func (a *MessageRepo) GetTopic() string {
	topic := a.kafkaTopics[a.currentTopicIndex]
	a.currentTopicIndex++
	if a.currentTopicIndex >= len(a.kafkaTopics)-1 {
		a.currentTopicIndex = 0
	}
	return topic
}

func (a *MessageRepo) SendTmplMsgFromRequest(ctx context.Context, param entity.SendTmplMsgRemoteReq) (entity.SendTmplMsgRemoteResp, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SendTmplMsgFromRequest traceID:%s", traceID)
	// 请求wx msg send
	bs, err := json.Marshal(param)
	if err != nil {
		log.Errorf("SendTmplMsgFromRequest json marshal send msg req failed,traceID:%s,err:%v", traceID, err)
		return entity.SendTmplMsgRemoteResp{}, err
	}
	requestProperty := httputil.GetRequestProperty(http.MethodPost, config.WXMsgTmplSendURL+fmt.Sprintf("?access_token=%s", param.AccessToken),
		bs, make(map[string]string))
	statusCode, body, _, err := httputil.RequestWithContextAndRepeat(ctx, requestProperty, traceID)
	if err != nil {
		log.Errorf("SendTmplMsgFromRequest request wx msg send failed, traceID:%s, error:%v", traceID, err)
		return entity.SendTmplMsgRemoteResp{}, err
	}
	if statusCode != http.StatusOK {
		log.Errorf("SendTmplMsgFromRequest request wx msg send failed, statusCode:%d,traceID:%s, error:%v", statusCode, traceID, err)
		return entity.SendTmplMsgRemoteResp{}, err
	}
	var msgResp entity.SendTmplMsgRemoteResp
	err = json.Unmarshal(body, &msgResp)
	if err != nil {
		log.Errorf("SendTmplMsgFromRequest get wx msg send failed by unmarshal, resp:%s, traceID:%s, err:%v", string(body), traceID, err)
		return entity.SendTmplMsgRemoteResp{}, err
	}
	// token过期
	if msgResp.ErrCode == errors.CodeRIDExpired {
		err = errors.NewCustomError(nil, errors.CodeTokenExpire, errors.GetErrorMessage(errors.CodeTokenExpire))
		return entity.SendTmplMsgRemoteResp{}, err
	}
	// 获取失败
	if msgResp.ErrCode != errors.CodeOK {
		log.Errorf("SendTmplMsgFromRequest get wx msg send failed,resp:%s,traceID:%s,errMsg:%s", string(body), traceID, msgResp.ErrMsg)
		return entity.SendTmplMsgRemoteResp{}, fmt.Errorf("get wx msg send failed,errMsg:%s", msgResp.ErrMsg)
	}
	return msgResp, nil
}

func (a *MessageRepo) SaveFailureMsgLog(ctx context.Context, param entity.FailureMsgLog) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SaveFailureMsgLog traceID:%s", traceID)
	if err := a.DB.Create(&param).Error; err != nil {
		log.Errorf("SaveFailureMsgLog create failure msg log failed,traceID:%s,err:%v", traceID, err)
		return err
	}
	return nil
}

func (a *MessageRepo) SendTmplMsgToMQ(ctx context.Context, topic string, message string) error {
	var err error
	for i := 0; i < 3; i++ {
		err = a.MQ.SendMessage(ctx, topic, message)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
	if err != nil {
		log.Errorf("SendTmplMsgToMQ failed by send message to MQ, message:%s,err:%v", message, err)
		return err
	}
	return nil
}

func (a *MessageRepo) UpdateFailureMsg(ctx context.Context, failure entity.FailureMsgLog) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("UpdateFailureMsgStatus traceID:%s", traceID)
	err := a.DB.Model(&entity.FailureMsgLog{}).Updates(failure).Error
	if err != nil {
		log.Errorf("UpdateFailureMsgStatus update failed,traceID:%s,err:%v", traceID, err)
		return err
	}
	return nil
}

func (a *MessageRepo) GetMaxCountFailureMsgByMsgID(ctx context.Context, msgID int64) (entity.FailureMsgLog, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("GetMaxCountFailureMsgByMsgID traceID:%s", traceID)
	var m entity.FailureMsgLog
	err := a.DB.Where("msg_id = ?", msgID).Order("`count` DESC").First(&m).Error
	if err != nil {
		log.Errorf("GetMaxCountFailureMsgByMsgID failed to get max count,traceID:%s,err:%v", traceID, err)
		return entity.FailureMsgLog{}, err
	}
	return m, err
}

func (a *MessageRepo) IsExistFailureMsgFromDB(ctx context.Context, fromUserName string, createTime int64) (bool, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("IsExistFailureMsgFromDB traceID:%s", traceID)
	var failureMsgLog entity.FailureMsgLog
	err := a.DB.Where("to_user = ? AND create_time = ?", fromUserName, createTime).First(&failureMsgLog).Error
	if err != nil {
		// 不存在记录
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
