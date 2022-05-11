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
	msg *persistence.MessageRepo
}

func NewMessageRepository(msg *persistence.MessageRepo) *MessageRepository {
	return &MessageRepository{
		msg: msg,
	}
}

func (t *MessageRepository) SendTmplMsg(ctx context.Context, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SendTmplMsg traceID:%s", traceID)
	wg := new(sync.WaitGroup)
	// 批量写入到kafka做消息推送
	ch := make(chan struct{}, 100)
	defer close(ch)
	for idx := range param.ToUsers {
		ch <- struct{}{}
		wg.Add(1)
		go func(idx int) {
			defer func() {
				wg.Done()
				<-ch
			}()
			bs, err := json.Marshal(param.TransferPerSendTmplMsg(idx).TransferKafkaTmplReq())
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

func (t *MessageRepository) SaveFailureMsg(ctx context.Context, param entity.FailureMsgLog) (err error) {
	return t.msg.SaveFailureMsgLog(ctx, param)
}
