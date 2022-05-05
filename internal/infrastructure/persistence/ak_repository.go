package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"public-platform-manager/internal/config"
	"public-platform-manager/internal/consts"
	"public-platform-manager/internal/domain/entity"
	redis2 "public-platform-manager/internal/infrastructure/pkg/redis"
	"public-platform-manager/internal/interfaces/httputil"
	"public-platform-manager/internal/utils"
	"time"

	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
)

type AkRepo struct {
	Redis *redis.UniversalClient
}

func NewAkRepo() *AkRepo {
	return &AkRepo{
		Redis: CommonRepositories.Redis,
	}
}

func (a *AkRepo) GetAccessTokenFromRequest(ctx context.Context) (entity.AccessTokenResp, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("getAccessTokenFromRequest traceID:%s", traceID)
	// 请求wx access token
	requestProperty := httputil.GetRequestProperty(http.MethodGet, config.WXAccessTokenURL+fmt.Sprintf("?grant_type=%s&appid=%s&secret=%s", consts.Credential, config.AppID, config.AppSecret),
		nil, make(map[string]string))
	statusCode, body, _, err := httputil.RequestWithContextAndRepeat(ctx, requestProperty, traceID)
	if err != nil {
		log.Errorf("request wx access token failed, traceID:%s, error:%v", traceID, err)
		return entity.AccessTokenResp{}, err
	}
	if statusCode != http.StatusOK {
		log.Errorf("request wx access token failed, statusCode:%d,traceID:%s, error:%v", statusCode, traceID, err)
		return entity.AccessTokenResp{}, err
	}
	var akResp entity.AccessTokenResp
	err = json.Unmarshal(body, &akResp)
	if err != nil {
		log.Errorf("get wx access token failed by unmarshal, resp:%s, traceID:%s, err:%v", string(body), traceID, err)
		return entity.AccessTokenResp{}, err
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
		oldAk, err = redis2.RGet(consts.RedisKeyAccessToken)
		if err != nil && err != redis.Nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		break
	}
	if err == nil || err == redis.Nil {
		return string(oldAk), nil
	}
	log.Errorf("GetAccessTokenFromRedis get wx access token from redis failed by redis,traceID:%s, err:%v", traceID, err)
	return "", err
}
