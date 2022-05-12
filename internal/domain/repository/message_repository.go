package repository

import (
	"context"

	"encoding/json"
	"sync"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/persistence"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	log "github.com/sirupsen/logrus"
)

type MessageRepository struct {
	msg  *persistence.MessageRepo
	user *persistence.UserRepo
}

func NewMessageRepository(msg *persistence.MessageRepo, user *persistence.UserRepo) *MessageRepository {
	return &MessageRepository{
		msg:  msg,
		user: user,
	}
}

func (t *MessageRepository) SendTmplMsg(ctx context.Context, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SendTmplMsg traceID:%s", traceID)
	// 先拿到接收者的open_id列表
	users, err := t.user.ListUserByPhones(ctx, param.ToUsersPhone)
	if err != nil {
		log.Errorf("SendTmplMsg UserRepo ListUserByPhones failed,traceID:%s,err:%v", traceID, err)
		return entity.SendTmplMsgResp{}, err
	}
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
			bs, err = json.Marshal(param.TransferPerSendTmplMsg(users[idx].OpenID).TransferKafkaTmplReq())
			if err != nil {
				log.Errorf("handlerTEMPLATESENDJOBFINISHEvent json marshal tmpl msg failed,traceID:%s,err:%v", traceID, err)
				return
			}
			err = t.msg.SendTmplMsgToMQ(ctx, t.msg.GetTopic(), string(bs))
			if err != nil {
				log.Errorf("SendTmplMsg SendTmplMsgToMQ failed,param is %s,traceID:%s,err:%v", string(bs), traceID, err)
				return
			}
			log.Info("send msg success")
		}(idx)
	}
	wg.Wait()
	return entity.SendTmplMsgResp{Msg: "success"}, nil
}
