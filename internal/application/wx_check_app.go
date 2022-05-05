package application

import (
	"public-platform-manager/internal/domain/repository"
)

type wxCheckSignApp struct {
	wx repository.WXCheckSignRepository
}

// wxCheckSignApp implements the WXCheckSignatureInterface
var _ WXCheckSignatureInterface = &wxCheckSignApp{}

type WXCheckSignatureInterface interface {
	GetWXCheckSign(signature, timestamp, nonce, token string) bool
}

func (w *wxCheckSignApp) GetWXCheckSign(signature, timestamp, nonce, token string) bool {
	return w.wx.GetWXCheckSign(signature, timestamp, nonce, token)
}
