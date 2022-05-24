package persistence

import (
	"fmt"
	redis3 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/redis"
	oslog "log"
	"os"

	"gorm.io/gorm/logger"

	"time"

	"gorm.io/gorm"

	smsPb "git.nova.net.cn/nova/notify/sms-xuanwu/pkg/grpcIFace"
	captchaPb "git.nova.net.cn/nova/shared/captcha/pkg/grpcIFace"
	"github.com/go-redis/redis/v7"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
)

type Repositories struct {
	DB                *gorm.DB
	Redis             *redis.UniversalClient
	SmsGRPCClient     smsPb.SenderClient
	CaptchaGRPCClient captchaPb.CaptchaServiceClient
}

type DBConfig struct {
	DBDriver, DBUser, DBPassword, DBHost, DBName string
	MaxIdleConn, MaxOpenConn                     int
}

var CommonRepositories Repositories

func NewRepositories(DBConfig DBConfig, redisAddresses []string, smsRPCAddr, captchaRPCAddr string, debugMode bool) error {
	err := NewDBRepositories(DBConfig, debugMode)
	if err != nil {
		return err
	}
	err = NewRedisRepositories(redisAddresses)
	if err != nil {
		return err
	}
	err = NewSmsGRPCClientRepositories(smsRPCAddr)
	if err != nil {
		return err
	}
	err = NewCaptchaGRPCClientRepositories(captchaRPCAddr)
	if err != nil {
		return err
	}

	// persistence repo init
	NewAkRepo()
	NewMessageRepo()
	NewUserRepo()
	NewWxRepo()
	NewPhoneVerifyRepo()
	return nil
}

func NewDBRepositories(config DBConfig, debugMode bool) error {
	dbUser, dbPassword, dbHost, dbName := config.DBUser, config.DBPassword, config.DBHost, config.DBName
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&interpolateParams=true", dbUser, dbPassword, dbHost, dbName)

	db, err := gorm.Open(mysql.Open(dataSource), &gorm.Config{})
	if err != nil {
		return err
	}
	sqlDB, _ := db.DB()
	if config.MaxIdleConn > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	}
	if config.MaxOpenConn > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Hour) // 设置最大连接超时
	}
	if debugMode {
		newLogger := logger.New(
			oslog.New(os.Stdout, "\r\n", oslog.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,   // Slow SQL threshold
				LogLevel:                  logger.Silent, // Log level
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,         // Disable color
			},
		)
		db.Logger = newLogger
	}

	CommonRepositories.DB = db

	return nil
}

func NewRedisRepositories(addresses []string) error {
	redisClient := redis3.NewRedisClient(addresses)
	err := redisClient.Ping().Err()
	if err != nil {
		return err
	}
	CommonRepositories.Redis = &redisClient
	log.Info("redis client init success")
	return nil
}

func NewSmsGRPCClientRepositories(smsRPCAddr string) error {
	smsConn, err := grpc.Dial(smsRPCAddr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("failed to dial sms grpc server: %+v", err)
		return err
	}
	smsClient := smsPb.NewSenderClient(smsConn)
	CommonRepositories.SmsGRPCClient = smsClient
	return nil
}

func NewCaptchaGRPCClientRepositories(captchaRPCAddr string) error {
	captchaConn, err := grpc.Dial(captchaRPCAddr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("failed to dial captcha grpc server: %+v", err)
		return err
	}

	captchaClient := captchaPb.NewCaptchaServiceClient(captchaConn)
	CommonRepositories.CaptchaGRPCClient = captchaClient
	return nil
}
