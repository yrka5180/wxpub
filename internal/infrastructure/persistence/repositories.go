package persistence

import (
	"fmt"
	"public-platform-manager/internal/infrastructure/pkg/kafka"
	redis2 "public-platform-manager/internal/infrastructure/pkg/redis"
	"time"

	"github.com/Shopify/sarama"

	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

type Repositories struct {
	MQ    *kafka.MQ
	DB    *gorm.DB
	Redis *redis.UniversalClient
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
