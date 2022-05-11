package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"public-platform-manager/internal/config"
	"public-platform-manager/internal/consts"
	"public-platform-manager/internal/domain/entity"
	"public-platform-manager/internal/domain/repository"
	"public-platform-manager/internal/infrastructure/persistence"
	"public-platform-manager/internal/tasks/consumer"
	"public-platform-manager/internal/tasks/g"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

type CurrentTopicManger struct {
	TopicIndex int
}

var (
	topicMgr     = &CurrentTopicManger{}
	taskCtx      context.Context
	MaxMsgChan   = make(chan string, 100)
	msgRepo      *persistence.MessageRepo
	akRepository *repository.AccessTokenRepository
)

func ConsumerTask(ctx context.Context) {
	taskCtx = ctx
	msgRepo = persistence.NewMessageRepo(config.KafkaTopics)
	akRepository = repository.NewAccessTokenRepository(persistence.NewAkRepo())
	// 消息任务处理
	g.Add(1)
	go listenConsumer(ctx)
	// 业务处理
	g.Add(1)
	go handleMsg(ctx)
}

func (tm *CurrentTopicManger) GetTopic() string {
	topic := config.KafkaTopics[tm.TopicIndex]
	tm.TopicIndex++
	if tm.TopicIndex >= len(config.KafkaTopics)-1 {
		tm.TopicIndex = 0
	}
	return topic
}

func consume(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-taskCtx.Done():
			log.Infof("consume exit..., err:%v", taskCtx.Err())
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			log.Debugf("receive topic[%s], msg: %s", msg.Topic, string(msg.Value))
			MaxMsgChan <- string(msg.Value)
			session.MarkMessage(msg, "")
		}
	}
}

func listenConsumer(ctx context.Context) {
	defer g.Done()
	queue := persistence.CommonRepositories.MQ
	customConsumer := consumer.NewCustomConsumer(consume)
	for {
		select {
		case <-ctx.Done():
			log.Infof("listenConsumer exit..., err:%v", ctx.Err())
			return
		default:
			if err := queue.ConsumerGroup.Consume(ctx, config.KafkaTopics, customConsumer); err != nil {
				log.Errorf("mq consumer handle topics wx-public proxy msg failed, err:%v", err)
				return
			}
		}
	}
}

// 处理消息，状态为发送中，针对消息失败记录失败次数，失败3次，状态改为发送失败
func handleMsg(ctx context.Context) {
	defer g.Done()
	for {
		select {
		case <-ctx.Done():
			allCount, successCount := len(MaxMsgChan), 0
			log.Infof("handleMsg get context done, still has data len[%d]", allCount)
			if allCount > 0 {
				for msg := range MaxMsgChan {
					err := msgRepo.SendTmplMsgToMQ(context.TODO(), topicMgr.GetTopic(), msg)
					if err != nil {
						log.Errorf("rewrite the msg to queue failed, msg:%v, err:%v", msg, err)
					}
					successCount++
					if len(MaxMsgChan) == 0 {
						break
					}
				}
			}
			log.Infof("handleMsg exit...,  data len:%d, rewrite count:%d, err:%v", allCount, successCount, ctx.Err())
			return
		case msg, ok := <-MaxMsgChan:
			if !ok {
				return
			}
			log.Debugf("recv msg is %s", msg)
			// 失败次数判断，状态重试
			item, err := validKafkaTmplMsg(msg)
			if err != nil {
				continue
			}
			g.Add(1)
			go func() {
				defer g.Done()
				// 业务处理
				var resp entity.SendTmplMsgRemoteResp
				var failureMsg entity.FailureMsgLog
				failureMsg = item.TransferSendRetryMsgLog("")
				// 获取access token
				var ak string
				ak, err = akRepository.GetAccessToken(ctx)
				item.SendTmplMsgRemoteReq.AccessToken = ak
				resp, err = msgRepo.SendTmplMsgFromRequest(ctx, item.SendTmplMsgRemoteReq)
				if err != nil {
					// 记录当前错误状态为重试中
					failureMsg = item.TransferSendRetryMsgLog(err.Error())
					failureMsg.Count = item.FailureCount
					// 回写
					retryToQueue(&item)
				}
				failureMsg.MsgID = resp.MsgID
				log.Debugf("resp msg id is %v", resp.MsgID)
				err = msgRepo.SaveFailureMsgLog(context.TODO(), failureMsg)
				if err != nil {
					log.Errorf("handleMsg SaveFailureMsgLog failed,err:%v", err)
				}
			}()
		}
	}
}

// 检验操作日志格式以及是否超过重试上限
func validKafkaTmplMsg(m string) (entity.KafkaTmplMsg, error) {
	var msg entity.KafkaTmplMsg
	err := json.Unmarshal([]byte(m), &msg)
	if err != nil {
		log.Errorf("validKafkaTmplMsg failed by json unmarshal, message:%v, err:%v", msg, err)
		return entity.KafkaTmplMsg{}, err
	}

	// 超过重试次数丢弃
	if msg.FailureCount > consts.MaxRetryCount {
		log.Errorf("msg has exceed max retry count[%d], will discard it, message:%v", consts.MaxRetryCount, msg)
		return entity.KafkaTmplMsg{}, fmt.Errorf("msg has exceed max retry count[%d]", consts.MaxExpireTime)
	}

	// 超过重试时间丢弃
	if msg.AcceptedTime > 0 && time.Now().Unix()-msg.AcceptedTime > consts.MaxExpireTime {
		log.Errorf("msg has exceed max expire time, will discard it, MaxExpireTime:%d, message:%v", consts.MaxExpireTime, msg)
		return entity.KafkaTmplMsg{}, fmt.Errorf("msg has exceed max expire time[%d]", consts.MaxExpireTime)
	}
	return msg, nil
}

// 消费失败重写回队列
func retryToQueue(msg *entity.KafkaTmplMsg) {
	msg.FailureCount++
	msg.AcceptedTime = time.Now().Unix()
	body, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("retry to queue failed by json marshal, msg:%v, err:%v", msg, err)
		return
	}
	err = msgRepo.SendTmplMsgToMQ(context.TODO(), topicMgr.GetTopic(), string(body))
	if err != nil {
		log.Errorf("retryToQueue failed, msg:%v, err:%v", msg, err)
	}
}
