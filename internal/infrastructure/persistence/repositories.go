package persistence

import (
	"fmt"
	"time"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/pkg/kafka"
	redis2 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/infrastructure/pkg/redis"

	config2 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/config"
	smsPb "git.nova.net.cn/nova/notify/sms-xuanwu/pkg/grpcIFace"
	captchaPb "git.nova.net.cn/nova/shared/captcha/pkg/grpcIFace"
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql" // for gorm
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Repositories struct {
	MQ                *kafka.MQ
	DB                *gorm.DB
	Redis             *redis.UniversalClient
	SmsGRPCClient     smsPb.SenderClient
	CaptchaGRPCClient captchaPb.CaptchaServiceClient
}

type KafkaConfig struct {
	Config          *sarama.Config
	Brokers         []string
	ConsumerGroupID string
	Topics          []string
	KafkaVersion    string
}

type DBConfig struct {
	DBDriver, DBUser, DBPassword, DBHost, DBName string
	MaxIdleConn, MaxOpenConn                     int
}

const (
	mysqlType = "mysql"
)

var CommonRepositories Repositories

func NewRepositories(KafkaConfig KafkaConfig, DBConfig DBConfig, redisAddresses []string, smsRPCAddr, captchaRPCAddr string, debugMode bool) error {
	err := NewDBRepositories(DBConfig, debugMode)
	if err != nil {
		return err
	}
	err = NewMQRepositories(KafkaConfig, debugMode)
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
	NewMessageRepo(config2.KafkaTopics)
	NewPassportRepo()
	NewUserRepo()
	NewWxRepo()
	NewPhoneVerifyRepo()
	return nil
}

func NewMQRepositories(conf KafkaConfig, debugMode bool) error {
	config, brokers, consumerGroupID, kafkaVersion := conf.Config, conf.Brokers, conf.ConsumerGroupID, conf.KafkaVersion
	// 生产者
	producer, err := kafka.InitProducer(config, brokers, kafkaVersion, debugMode)
	if err != nil {
		return err
	}

	// 正常业务消费者
	consumer, err := kafka.InitKafkaConsumerGroup(config, brokers, consumerGroupID, kafkaVersion, false)
	if err != nil {
		return err
	}

	CommonRepositories.MQ = &kafka.MQ{
		Producer:      producer,
		ConsumerGroup: consumer,
	}

	return nil
}

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

func NewSmsGRPCClientRepositories(smsRPCAddr string) error {
	smsConn, err := grpc.Dial(smsRPCAddr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("failed to dial sms grpc server: %v", err)
		return err
	}
	smsClient := smsPb.NewSenderClient(smsConn)
	CommonRepositories.SmsGRPCClient = smsClient
	return nil
}

func NewCaptchaGRPCClientRepositories(captchaRPCAddr string) error {
	captchaConn, err := grpc.Dial(captchaRPCAddr, grpc.WithInsecure())
	if err != nil {
		log.Errorf("failed to dial captcha grpc server: %v", err)
		return err
	}

	captchaClient := captchaPb.NewCaptchaServiceClient(captchaConn)
	CommonRepositories.CaptchaGRPCClient = captchaClient
	return nil
}
