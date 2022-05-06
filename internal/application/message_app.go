package application

import (
	"context"
	"public-platform-manager/internal/domain/entity"
	"public-platform-manager/internal/domain/repository"
)

type messageApp struct {
	message repository.MessageRepository
}

// messageApp implements the MessageInterface
var _ MessageInterface = &messageApp{}

type MessageInterface interface {
	SendTmplMsg(ctx context.Context, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, error)
}

func (u *messageApp) SendTmplMsg(ctx context.Context, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, error) {
	return u.message.SendTmplMsg(ctx, param)
}
