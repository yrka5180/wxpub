package repository

import (
	"context"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/persistence"
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
	return t.msg.SendTmplMsgFromRequest(ctx, param)
}
