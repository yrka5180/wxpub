package application

import (
	"context"
	"public-platform-manager/internal/domain/entity"
	"public-platform-manager/internal/domain/repository"
)

type wxApp struct {
	wx repository.WXRepository
}

// wxApp implements the WXInterface
var _ WXInterface = &wxApp{}

type WXInterface interface {
	GetWXCheckSign(signature, timestamp, nonce, token string) bool
	GetEventXml(ctx context.Context, reqBody *entity.TextRequestBody) (respBody []byte, err error)
}

func (w *wxApp) GetWXCheckSign(signature, timestamp, nonce, token string) bool {
	return w.wx.GetWXCheckSign(signature, timestamp, nonce, token)
}

func (w *wxApp) GetEventXml(ctx context.Context, reqBody *entity.TextRequestBody) (respBody []byte, err error) {
	return w.wx.GetEventXml(ctx, reqBody)
}
