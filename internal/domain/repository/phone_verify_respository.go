package repository

import (
	"context"
	"fmt"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/config"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/persistence"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"
)

type PhoneVerifyRepository struct {
	phone *persistence.PhoneVerifyRepo
}

var defaultPhoneVerifyRepository PhoneVerifyRepository

func InitDefaultPhoneVerifyRepository(phone *persistence.PhoneVerifyRepo) {
	if defaultPhoneVerifyRepository.phone == nil {
		defaultPhoneVerifyRepository.phone = phone
	}
}

func DefaultPhoneVerifyRepository() *PhoneVerifyRepository {
	return &defaultPhoneVerifyRepository
}

func (a *PhoneVerifyRepository) SendSms(ctx context.Context, req entity.SendSmsReq) error {
	verifyCodeID, verifyCodeAnswer := utils.GenVerifySmsCode()
	err := a.phone.SetVerifyCodeSmsStorage(ctx, req.OpenID, verifyCodeID, verifyCodeAnswer)
	if err != nil {
		return err
	}

	content := fmt.Sprintf(config.SmsContentTemplateCN, verifyCodeAnswer)
	sender := consts.SmsSender
	return a.phone.SendSms(ctx, content, sender, req.Phone)
}
