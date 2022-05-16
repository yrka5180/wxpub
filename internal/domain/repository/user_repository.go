package repository

import (
	"context"
	"fmt"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/config"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/persistence"
)

type UserRepository struct {
	user        *persistence.UserRepo
	phoneVerify *persistence.PhoneVerifyRepo
}

var defaultUserRepository = &UserRepository{}

func NewUserRepository(user *persistence.UserRepo, phoneVerify *persistence.PhoneVerifyRepo) {
	if defaultUserRepository.user == nil {
		defaultUserRepository.user = user
	}
	if defaultUserRepository.phoneVerify == nil {
		defaultUserRepository.phoneVerify = phoneVerify
	}
}

func DefaultUserRepository() *UserRepository {
	return defaultUserRepository
}

func (a *UserRepository) ListUser(ctx context.Context) ([]entity.User, error) {
	return a.user.ListUser(ctx)
}

func (a *UserRepository) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	return a.user.GetUserByID(ctx, id)
}

func (a *UserRepository) GetUserByOpenID(ctx context.Context, openID string) (entity.User, error) {
	return a.user.GetUserByOpenID(ctx, openID)
}

func (a *UserRepository) SaveUser(ctx context.Context, user entity.User) error {
	return a.user.SaveUser(ctx, user)
}

func (a *UserRepository) SendSms(ctx context.Context, req entity.SendSmsReq) error {
	verifyCodeID, verifyCodeAnswer := utils.GenVerifySmsCode()
	err := a.phoneVerify.SetVerifyCodeSmsStorage(ctx, req.OpenID, verifyCodeID, verifyCodeAnswer)
	if err != nil {
		return err
	}

	content := fmt.Sprintf(config.SmsContentTemplateCN, verifyCodeAnswer)
	sender := consts.SmsSender
	return a.phoneVerify.SendSms(ctx, content, sender, req.Phone)
}

func (a *UserRepository) VerifySmsCode(ctx context.Context, req entity.VerifyCodeReq) (bool, bool, error) {
	return a.phoneVerify.VerifySmsCode(ctx, req.OpenID, consts.RedisKeyVerifyCodeSmsID, req.VerifyCode, consts.RedisAuthTTL)
}
