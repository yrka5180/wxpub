package persistence

import (
	"context"
	"time"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
	redis2 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/pkg/redis"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type WxRepo struct {
	DB    *gorm.DB
	Redis *redis.UniversalClient
}

func NewWxRepo() *WxRepo {
	return &WxRepo{
		DB:    CommonRepositories.DB,
		Redis: CommonRepositories.Redis,
	}
}

func (a *WxRepo) SetMsgIDToRedis(ctx context.Context, msgID string) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SetMsgIDToRedis traceID:%s", traceID)
	var err error
	for i := 0; i < 3; i++ {
		err = redis2.RSet(consts.RedisKeyMsgID+msgID, "", consts.RedisMsgIDTTL)
		if err != nil {
			log.Errorf("SetMsgIDToRedis WxRepo redis set msg id failed,traceID:%s,err:%v", traceID, err)
			continue
		}
		break
	}
	return err
}

func (a *WxRepo) IsExistMsgIDFromRedis(ctx context.Context, msgID string) (bool, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("IsExistMsgIDFromRedis traceID:%s", traceID)
	var err error
	for i := 0; i < 3; i++ {
		_, err = redis2.RGet(consts.RedisKeyMsgID + msgID)
		if err != nil {
			if err == redis.Nil {
				return false, nil
			} else {
				time.Sleep(10 * time.Millisecond)
				continue
			}
		}
		break
	}
	if err != nil {
		log.Errorf("IsExistMsgIDFromRedis get wx msg id from redis failed,traceID:%s, err:%v", traceID, err)
		return false, err
	}
	return true, nil
}