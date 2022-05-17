package application

import (
	"context"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/repository"
)

type messageApp struct {
	message repository.MessageRepository
}

// messageApp implements the MessageInterface
var _ MessageInterface = &messageApp{}

type MessageInterface interface {
	GetMissingUsers(ctx context.Context, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, []entity.User, error)
	SendTmplMsg(ctx context.Context, users []entity.User, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, error)
}

func (u *messageApp) GetMissingUsers(ctx context.Context, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, []entity.User, error) {
	return u.message.GetMissingUsers(ctx, param)
}

func (u *messageApp) SendTmplMsg(ctx context.Context, users []entity.User, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, error) {
	return u.message.SendTmplMsg(ctx, users, param)
}
