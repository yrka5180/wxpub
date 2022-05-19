package repository

import (
	"context"
	"encoding/json"
	"sync"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/errors"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/persistence"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	log "github.com/sirupsen/logrus"
)

type MessageRepository struct {
	msg  *persistence.MessageRepo
	user *persistence.UserRepo
}

var defaultMessageRepository = &MessageRepository{}

func NewMessageRepository(msg *persistence.MessageRepo, user *persistence.UserRepo) {
	if defaultMessageRepository.msg == nil {
		defaultMessageRepository.msg = msg
	}
	if defaultMessageRepository.user == nil {
		defaultMessageRepository.user = user
	}
}

func DefaultMessageRepository() *MessageRepository {
	return defaultMessageRepository
}

func (t *MessageRepository) GetMissingUsers(ctx context.Context, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, []entity.User, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("GetMissingUsers traceID:%s", traceID)
	var resp entity.SendTmplMsgResp
	resp.FailureSendPhones = make([]string, 0)
	// 先拿到接收者的open_id列表
	users, err := t.user.ListUserByPhones(ctx, param.ToUsersPhone)
	if err != nil {
		log.Errorf("SendTmplMsg UserRepo ListUserByPhones failed,traceID:%s,err:%v", traceID, err)
		return entity.SendTmplMsgResp{}, nil, err
	}
	userPhoneMap := make(map[string]struct{})
	for _, user := range users {
		userPhoneMap[user.Phone] = struct{}{}
	}
	// 判断手机号是否不存在
	for _, phone := range param.ToUsersPhone {
		// 手机号不存在记录
		if _, ok := userPhoneMap[phone]; !ok {
			resp.FailureSendPhones = append(resp.FailureSendPhones, phone)
		}
	}
	if len(resp.FailureSendPhones) > 0 {
		return resp, nil, errors.NewCustomError(nil, errors.CodeResourcesPartialNotFound, errors.GetErrorMessage(errors.CodeResourcesPartialNotFound))
	}
	return resp, users, nil
}

func (t *MessageRepository) SendTmplMsg(ctx context.Context, users []entity.User, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SendTmplMsg traceID:%s", traceID)
	var resp entity.SendTmplMsgResp
	var err error
	sendMsgID, err := utils.GetUUID()
	if err != nil {
		log.Errorf("SendTmplMsg MessageRepository GetUUID failed,traceID:%s,err:%v", traceID, err)
		return resp, err
	}
	resp.SendMsgID = sendMsgID
	wg := new(sync.WaitGroup)
	// 批量写入到kafka做消息推送
	ch := make(chan struct{}, 100)
	defer close(ch)
	for idx := range users {
		ch <- struct{}{}
		wg.Add(1)
		go func(idx int) {
			defer func() {
				wg.Done()
				<-ch
			}()
			var bs []byte
			bs, err = json.Marshal(param.TransferPerSendTmplMsg(users[idx].OpenID).TransferKafkaTmplReq(sendMsgID))
			if err != nil {
				log.Errorf("handlerTEMPLATESENDJOBFINISHEvent json marshal tmpl msg failed,traceID:%s,err:%v", traceID, err)
				return
			}
			err = t.msg.SendTmplMsgToMQ(ctx, t.msg.GetTopic(), string(bs))
			if err != nil {
				log.Errorf("SendTmplMsg SendTmplMsgToMQ failed,param is %s,traceID:%s,err:%v", string(bs), traceID, err)
				return
			}
			log.Debugf("send msg success,msg is %v", string(bs))
		}(idx)
	}
	wg.Wait()
	return resp, nil
}
