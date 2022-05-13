package application

import (
	"context"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/repository"
)

type userApp struct {
	user repository.UserRepository
}

// userApp implements the UserInterface
var _ UserInterface = &userApp{}

type UserInterface interface {
	ListUser(ctx context.Context) ([]entity.User, error)
	GetUserByID(ctx context.Context, id int) (entity.User, error)
	SendSms(ctx context.Context, req entity.SendSmsReq) error
}

func (u *userApp) ListUser(ctx context.Context) ([]entity.User, error) {
	return u.user.ListUser(ctx)
}

func (u *userApp) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	return u.user.GetUserByID(ctx, id)
}

func (u *userApp) SendSms(ctx context.Context, req entity.SendSmsReq) error {
	return u.user.SendSms(ctx, req)
}
