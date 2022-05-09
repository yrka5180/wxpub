package repository

import (
	"context"
	"public-platform-manager/internal/domain/entity"
	"public-platform-manager/internal/infrastructure/persistence"
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
