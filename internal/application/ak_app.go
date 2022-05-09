package application

import (
	"context"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/repository"
)

type akApp struct {
	ak repository.AccessTokenRepository
}

// akApp implements the AccessTokenInterface
var _ AccessTokenInterface = &akApp{}

type AccessTokenInterface interface {
	GetAccessToken(ctx context.Context) (string, error)
	FreshAccessToken(ctx context.Context) (string, error)
}

func (a *akApp) GetAccessToken(ctx context.Context) (string, error) {
	return a.ak.GetAccessToken(ctx)
}

func (a *akApp) FreshAccessToken(ctx context.Context) (string, error) {
	return a.ak.FreshAccessToken(ctx)
}
