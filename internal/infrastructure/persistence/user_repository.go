package persistence

import (
	"context"
	"time"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	"github.com/go-redis/redis/v7"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type UserRepo struct {
	DB    *gorm.DB
	Redis *redis.UniversalClient
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		DB:    CommonRepositories.DB,
		Redis: CommonRepositories.Redis,
	}
}

func (a *UserRepo) IsExistUserMsgFromDB(ctx context.Context, fromUserName string, createTime int64) (bool, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("IsExistUserMsgFromDB traceID:%s", traceID)
	var user entity.User
	err := a.DB.Where("open_id = ? AND create_time = ?", fromUserName, createTime).First(&user).Error
	if err != nil {
		// 不存在记录
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a *UserRepo) IsExistUserFromDB(ctx context.Context, fromUserName string) (bool, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("IsExistUserFromDB traceID:%s", traceID)
	var user entity.User
	err := a.DB.Where("open_id = ?", fromUserName).First(&user).Error
	if err != nil {
		// 不存在记录
		if gorm.IsRecordNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (a *UserRepo) SaveUser(ctx context.Context, user entity.User) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SaveUser traceID:%s", traceID)
	// 先查看是否有这用户，如果没有则创建，否则将创建时间和删除时间更新
	exist, err := a.IsExistUserFromDB(ctx, user.OpenID)
	if err != nil {
		log.Errorf("SaveUser IsExistUserFromDB failed,traceID:%s,err:%v", traceID, err)
		return err
	}
	if !exist {
		if err = a.DB.Create(&user).Error; err != nil {
			log.Errorf("SaveUser create user failed,traceID:%s,err:%v", traceID, err)
			return err
		}
	}
	if err = a.UpdateUser(ctx, user); err != nil {
		log.Errorf("SaveUser UpdateUser failed,traceID:%s,err:%v", traceID, err)
		return err
	}
	return nil
}

func (a *UserRepo) DelUser(ctx context.Context, user entity.User) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("DelUser traceID:%s", traceID)
	user.DeleteTime = time.Now().Unix()
	if err := a.DB.Model(&entity.User{}).Where("open_id = ?", user.OpenID).Update(&user).Error; err != nil {
		log.Errorf("DelUser delete user failed,traceID:%s,err:%v", traceID, err)
		return err
	}
	return nil
}

func (a *UserRepo) UpdateUser(ctx context.Context, user entity.User) error {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("UpdateUser traceID:%s", traceID)
	if err := a.DB.Model(&entity.User{}).Where("open_id = ?", user.OpenID).Updates(map[string]interface{}{
		"create_time": user.CreateTime,
		"delete_time": user.DeleteTime,
	}).Error; err != nil {
		log.Errorf("UpdateUser update user failed,traceID:%s,err:%v", traceID, err)
		return err
	}
	return nil
}

func (a *UserRepo) ListUser(ctx context.Context) (users []entity.User, err error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("ListUser traceID:%s", traceID)
	if err = a.DB.Where("delete_time = 0").Find(&users).Error; err != nil {
		log.Errorf("ListUser find list users failed,traceID:%s,err:%v", traceID, err)
		return nil, err
	}
	return users, nil
}

func (a *UserRepo) GetUserByID(ctx context.Context, id int) (user entity.User, err error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("GetUserByID traceID:%s", traceID)
	if err = a.DB.Where("id = ?", id).First(&user).Error; err != nil {
		log.Errorf("GetUserByID get user by id failed,traceID:%s,err:%v", traceID, err)
		return
	}
	return
}

func (a *UserRepo) ListUserByPhones(ctx context.Context, phones []string) (user []entity.User, err error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("ListUserByPhones traceID:%s", traceID)
	if err = a.DB.Where("id IN (?)", phones).First(&user).Error; err != nil {
		log.Errorf("ListUserByPhones get list user by phones failed,traceID:%s,err:%v", traceID, err)
		return
	}
	return
}
