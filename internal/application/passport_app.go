package application

import (
	"context"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/repository"
)

type passportApp struct {
	passport repository.PassportRepository
}

// passportApp implements the PassportAppInterface
var _ PassportAppInterface = &passportApp{}

type PassportAppInterface interface {
	GetAuthN(ctx context.Context, auth string) error
}

func (p *passportApp) GetAuthN(ctx context.Context, auth string) error {
	return p.passport.GetAuthN(ctx, auth)
}
