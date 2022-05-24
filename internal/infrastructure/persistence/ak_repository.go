package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	errors2 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/errorx"
	httputil2 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/ginx/httputil"
	redis3 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/redis"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/config"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

type AkRepo struct {
	Redis *redis.UniversalClient
}

var defaultAkRepo *AkRepo

func NewAkRepo() {
	if defaultAkRepo == nil {
		defaultAkRepo = &AkRepo{
			Redis: CommonRepositories.Redis,
		}
	}
}

func DefaultAkRepo() *AkRepo {
	return defaultAkRepo
}

func (a *AkRepo) GetAccessTokenFromRequest(ctx context.Context) (entity.AccessTokenResp, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("getAccessTokenFromRequest traceID:%s", traceID)
	// 请求wx access token
	requestProperty := httputil2.GetRequestProperty(http.MethodGet, config.WXAccessTokenURL+fmt.Sprintf("?grant_type=%s&appid=%s&secret=%s", consts.Credential, config.AppID, config.AppSecret),
		nil, make(map[string]string))
	statusCode, body, _, err := httputil2.RequestWithContextAndRepeat(ctx, requestProperty, traceID)
	if err != nil {
		log.Errorf("request wx access token failed, traceID:%s, error:%+v", traceID, err)
		return entity.AccessTokenResp{}, err
	}
	if statusCode != http.StatusOK {
		log.Errorf("request wx access token failed, statusCode:%d,traceID:%s, error:%+v", statusCode, traceID, err)
		return entity.AccessTokenResp{}, err
	}
	var akResp entity.AccessTokenResp
	err = json.Unmarshal(body, &akResp)
	if err != nil {
		log.Errorf("get wx access token failed by unmarshal, resp:%s, traceID:%s, err:%+v", string(body), traceID, err)
		return entity.AccessTokenResp{}, err
	}
	// 获取失败
	if akResp.ErrCode != errors2.CodeOK {
		log.Errorf("get wx access token failed,resp:%s,traceID:%s,errMsg:%s", string(body), traceID, akResp.ErrMsg)
		return entity.AccessTokenResp{}, fmt.Errorf("get wx ak failed,errMsg:%s", akResp.ErrMsg)
	}
	return akResp, nil
}

func (a *AkRepo) GetAccessTokenFromRedis(ctx context.Context) (string, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("GetAccessTokenFromRedis traceID:%s", traceID)
	// 先从redis中获取老accessToken
	var oldAk []byte
	var err error
	for i := 0; i < 3; i++ {
		oldAk, err = redis3.RGet(consts.RedisKeyAccessToken)
		if err != nil && err != redis.Nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		break
	}
	if err == nil || err == redis.Nil {
		return string(oldAk), nil
	}
	log.Errorf("GetAccessTokenFromRedis get wx access token from redis failed by redis,traceID:%s, err:%+v", traceID, err)
	return "", err
}

func (a *AkRepo) SetAccessTokenToRedis(ctx context.Context, accessToken string, expiresIn int) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SetAccessTokenToRedis traceID:%s", traceID)
	err := redis3.RSet(consts.RedisKeyAccessToken, accessToken, expiresIn)
	if err != nil {
		log.Errorf("SetAccessTokenToRedis AkRepo redis set new ak failed,traceID:%s,err:%+v", traceID, err)
		return err
	}
	return nil
}
