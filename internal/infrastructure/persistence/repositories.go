package persistence

import (
	"fmt"
	"time"

	redis2 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/pkg/redis"

	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql" // for gorm
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type Repositories struct {
	DB    *gorm.DB
	Redis *redis.UniversalClient
}

type DBConfig struct {
	DBDriver, DBUser, DBPassword, DBHost, DBName string
	MaxIdleConn, MaxOpenConn                     int
}

const (
	mysqlType = "mysql"
)

var CommonRepositories Repositories

func NewDBRepositories(config DBConfig, debugMode bool) error {
	dbUser, dbPassword, dbHost, dbName := config.DBUser, config.DBPassword, config.DBHost, config.DBName

	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4", dbUser, dbPassword, dbHost, dbName)
	db, err := gorm.Open(mysqlType, dataSource)
	if err != nil {
		return err
	}
	if config.MaxIdleConn > 0 {
		db.DB().SetMaxIdleConns(config.MaxIdleConn)
	}
	if config.MaxOpenConn > 0 {
		db.DB().SetMaxOpenConns(config.MaxOpenConn)
		db.DB().SetConnMaxLifetime(time.Hour) // 设置最大连接超时
	}
	if debugMode {
		db.LogMode(true)
	}

	CommonRepositories.DB = db

	return nil
}

func NewRedisRepositories(addresses []string) error {
	redisClient := redis2.NewRedisClient(addresses)
	err := redisClient.Ping().Err()
	if err != nil {
		return err
	}
	CommonRepositories.Redis = &redisClient
	log.Info("redis client init success")
	return nil
}
