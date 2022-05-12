package repository

import (
	"context"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/persistence"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"
	log "github.com/sirupsen/logrus"
)

type PassportRepository struct {
	passport *persistence.PassportRepo
}

func NewPassportRepository(passport *persistence.PassportRepo) *PassportRepository {
	return &PassportRepository{
		passport: passport,
	}
}

func (p *PassportRepository) GetAuthN(ctx context.Context, auth string) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("GetAuthN traceID:%s", traceID)
	var e error
	var authN *entity.AuthN
	// 尝试从缓存中取，过期时间为5分钟
	authN, e = p.passport.GetAuthFromRedis(ctx, auth)
	if e != nil {
		log.Errorf("GetAuthN PassportRepository getAuthNFromRedis failed,traceID:%s,err:%v", traceID, e)
	}
	if authN != nil {
		return nil
	}
	// 请求passport request
	_, err := p.passport.GetAuthFromRequest(ctx, auth)
	if err != nil {
		log.Errorf("GetAuthN getAuthNFromRequest failed,traceID:%s,err:%v", traceID, err)
		return err
	}
	e = p.passport.SetAuthN2Redis(ctx, authN, auth)
	if e != nil {
		log.Errorf("GetAuthN setAuthN2Redis failed,traceID:%s,err:%v", traceID, e)
	}
	return nil
}
