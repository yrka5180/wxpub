package persistence

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/config"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	redis2 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/pkg/redis"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/errors"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/httputil"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

type PassportRepo struct {
	Redis *redis.UniversalClient
}

var defaultPassportRepo *PassportRepo

func NewPassportRepo() *PassportRepo {
	if defaultPassportRepo == nil {
		defaultPassportRepo = &PassportRepo{
			Redis: CommonRepositories.Redis,
		}
	}
	return defaultPassportRepo
}

func NewDefaultPassportRepo() *PassportRepo {
	return defaultPassportRepo
}

func (p *PassportRepo) GetAuthFromRedis(ctx context.Context, auth string) (authN *entity.AuthN, err error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("GetAuthFromRedis traceID:%s", traceID)
	var value []byte
	for i := 0; i < 3; i++ {
		value, err = redis2.RGet(consts.RedisKeyAuthN + auth)
		if err != nil && err != redis.Nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		break
	}
	if err != nil {
		log.Errorf("GetAuthFromRedis get auth from redis failed,traceID:%s, err:%v", traceID, err)
		return nil, err
	}
	err = json.Unmarshal(value, &authN)
	if err != nil {
		log.Errorf("GetAuthFromRedis unmarshal failed,traceID:%s,err:%v", traceID, err)
		return nil, err
	}
	return
}

func (p *PassportRepo) GetAuthFromRequest(ctx context.Context, auth string) (*entity.AuthN, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("GetAuthFromRequest traceID:%s", traceID)
	header := map[string]string{
		consts.Authorization: auth,
	}
	code, body, _, err := httputil.RequestWithRepeat(traceID, http.MethodPost, config.PassportOIDCIntrospectURL, nil, header)
	if err != nil {
		log.Errorf("getAuthNFromRequest oidc introspect failed,traceID:%s,err:%v", traceID, err)
		err = errors.NewCustomError(err, errors.CodeInternalServerError, "verifyToken oidc introspect failed")
		return nil, err
	}
	if code != http.StatusOK {
		log.Errorf("getAuthNFromRequest oidc introspect token invalid,traceID:%s,err:%v", traceID, err)
		err = errors.NewCustomError(nil, errors.CodeUnknownError, "verifyToken authorization oidc introspect token invalid")
		return nil, err
	}
	var authN entity.AuthN
	err = json.Unmarshal(body, &authN)
	if err != nil {
		log.Errorf("getAuthNFromRequest json unmarshal auth information failed,traceID:%s,err:%v", traceID, err)
		err = errors.NewCustomError(err, errors.CodeInternalServerError, "verifyToken json unmarshal auth information failed")
		return nil, err
	}
	if !authN.Active {
		log.Errorf("getAuthNFromRequest oidc authorization token is not active,traceID:%s", traceID)
		err = errors.NewCustomError(nil, errors.CodeForbidden, "verifyToken authorization oidc introspect token failed,token is not active")
		return nil, err
	}
	return &authN, nil
}

func (p *PassportRepo) SetAuthN2Redis(ctx context.Context, authN *entity.AuthN, auth string) (err error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SetAuthN2Redis traceID:%s", traceID)
	bs, err := json.Marshal(authN)
	if err != nil {
		log.Errorf("setAuthN2Redis marshal failed,traceID:%s,auth:%s,err:%v", traceID, auth, err)
		return
	}
	err = redis2.RSet(consts.RedisKeyAuthN+auth, bs, consts.RedisAuthTTL)
	if err != nil {
		log.Errorf("setAuthN2Redis redis set failed,traceID:%s,err:%v", traceID, err)
	}
	return
}
