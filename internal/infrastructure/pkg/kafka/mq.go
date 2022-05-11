package kafka

import (
	"context"
	osLog "log"
	"os"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

type MQ struct {
	Producer      sarama.SyncProducer
	ConsumerGroup sarama.ConsumerGroup
}

var (
	logger = osLog.New(os.Stderr, "", osLog.LstdFlags)
)

func InitProducer(config *sarama.Config, kafkaBrokers []string, kafkaVersion string, debugMode bool) (sarama.SyncProducer, error) {
	var err error
	if config == nil {
		config = sarama.NewConfig()
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
		config.Version, err = sarama.ParseKafkaVersion(kafkaVersion)
		if err != nil {
			log.Errorf("kafka Producer ParseKafkaVersion failed, err:%v", err)
			return nil, err
		}
	}

	// 调试kafka的时候可以打开日志
	if debugMode {
		sarama.Logger = logger
	}

	producer, err := sarama.NewSyncProducer(kafkaBrokers, config)
	if err != nil {
		log.Errorf("kafka Producer NewSyncProducer failed, err:%v", err)
		return nil, err
	}
	return producer, nil
}

func InitKafkaConsumerGroup(config *sarama.Config, kafkaBrokers []string, groupID string, kafkaVersion string, isSync bool) (sarama.ConsumerGroup, error) {
	var err error
	if config == nil {
		config = sarama.NewConfig()
		config.Consumer.Return.Errors = true
		config.Version, err = sarama.ParseKafkaVersion(kafkaVersion)
		if isSync {
			config.Consumer.Offsets.AutoCommit.Enable = false
		}
		if err != nil {
			log.Errorf("kafka Consumer ParseKafkaVersion failed, err:%v", err)
			return nil, err
		}
	}

	consumerGroup, err := sarama.NewConsumerGroup(kafkaBrokers, groupID, config)
	if err != nil {
		log.Errorf("kafka Consumer NewConsumer failed, err:%v", err)
		return nil, err
	}

	go func() {
		for err := range consumerGroup.Errors() {
			log.Errorf("consumer consume, groupID:%s, error: %s", groupID, err.Error())
		}
	}()

	return consumerGroup, nil
}

func (mq *MQ) SendMessage(ctx context.Context, topic, message string) error {
	select {
	case <-ctx.Done():
		log.Infof("kafka producer send message process exit, topic:%s, ctx err:%v", topic, ctx.Err())
		return ctx.Err()
	default:
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(message),
		}
		// 发送消息
		partition, offset, err := mq.Producer.SendMessage(msg)
		if err != nil {
			log.Errorf("kafka producer send message failed, topic:%s, err:%v", topic, err)
			return err
		}
		log.Debugf("kafka producer send message successful, topic:%s, partition:%d, offset:%v, message:%s", topic, partition, offset, message)
		return nil
	}
}
