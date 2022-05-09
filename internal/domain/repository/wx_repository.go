package repository

import (
	"context"
	"encoding/xml"
	"fmt"
	"public-platform-manager/internal/consts"
	"public-platform-manager/internal/domain/entity"
	"public-platform-manager/internal/infrastructure/persistence"
	"public-platform-manager/internal/utils"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type WXRepository struct {
	wx   *persistence.WxRepo
	user *persistence.UserRepo
}

func NewWXRepository(wx *persistence.WxRepo, user *persistence.UserRepo) *WXRepository {
	return &WXRepository{
		wx:   wx,
		user: user,
	}
}

func (a *WXRepository) GetWXCheckSign(signature, timestamp, nonce, token string) bool {
	// 本地计算signature
	si := []string{token, timestamp, nonce}
	// 字典序排序
	sort.Strings(si)
	n := len(timestamp) + len(nonce) + len(token)
	var b strings.Builder
	b.Grow(n)
	for _, v := range si {
		b.WriteString(v)
	}
	return utils.Sha1(b.String()) == signature
}

func (a *WXRepository) GetEventXml(ctx context.Context, reqBody *entity.TextRequestBody) (respBody []byte, err error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("GetEventXml traceID:%s", traceID)
	if reqBody == nil {
		return nil, fmt.Errorf("xml request body is empty")
	}
	responseTextBody, err := a.handlerEvent(ctx, reqBody)
	if err != nil {
		log.Errorf("GetEventXml handlerEvent failed traceID:%s,err:%v", traceID, err)
		return nil, err
	}
	return responseTextBody, nil
}

func (a *WXRepository) handlerEvent(ctx context.Context, reqBody *entity.TextRequestBody) ([]byte, error) {
	var respContent string
	var err error
	// 事件类型
	switch reqBody.Event {
	// 关注订阅
	case consts.SubscribeEvent:
		if respContent, err = a.handlerSubscribeEvent(ctx, reqBody); err != nil {
			return nil, err
		}
	case consts.UnsubscribeEvent:
		if respContent, err = a.handlerUnSubscribeEvent(ctx, reqBody); err != nil {
			return nil, err
		}
	case consts.TEMPLATESENDJOBFINISHEvent:
	}
	return a.makeTextResponseBody(reqBody.ToUserName, reqBody.FromUserName, respContent)
}

func (a *WXRepository) handlerSubscribeEvent(ctx context.Context, reqBody *entity.TextRequestBody) (string, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("handlerSubscribeEvent traceID:%s", traceID)
	// 判断是否存在该消息id,用 FromUserName+CreateTime 去重
	msgID := fmt.Sprintf("%s%d", reqBody.FromUserName, reqBody.CreateTime)
	exist, err := a.isExistUserMsgID(ctx, msgID, reqBody.FromUserName, reqBody.CreateTime)
	if exist {
		return "", nil
	}
	// 持久化保存
	u := entity.User{
		OpenID:     reqBody.FromUserName,
		CreateTime: reqBody.CreateTime,
	}
	err = a.user.SaveUser(ctx, u)
	if err != nil {
		log.Errorf("handlerSubscribeEvent WXRepository wx repo SaveUser traceID:%s,err:%v", traceID, err)
		return "", err
	}
	err = a.wx.SetMsgIDToRedis(ctx, msgID)
	if err != nil {
		log.Errorf("handlerSubscribeEvent WXRepository wx repo set msg id to redis failed,traceID:%s,err:%v", traceID, err)
	}
	return consts.SubscribeRespContent, nil
}

func (a *WXRepository) handlerUnSubscribeEvent(ctx context.Context, reqBody *entity.TextRequestBody) (string, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("handlerUnSubscribeEvent traceID:%s", traceID)
	// 判断是否存在该消息id,用 FromUserName+CreateTime 去重
	msgID := fmt.Sprintf("%s%d", reqBody.FromUserName, reqBody.CreateTime)
	exist, err := a.isExistUserMsgID(ctx, msgID, reqBody.FromUserName, reqBody.CreateTime)
	if exist {
		return "", nil
	}
	// 删除用户信息
	u := entity.User{
		OpenID: reqBody.FromUserName,
	}
	err = a.user.DelUser(ctx, u)
	if err != nil {
		log.Errorf("handlerUnSubscribeEvent WXRepository wx repo SaveUser traceID:%s,err:%v", traceID, err)
		return "", err
	}
	err = a.wx.SetMsgIDToRedis(ctx, msgID)
	if err != nil {
		log.Errorf("handlerUnSubscribeEvent WXRepository wx repo set msg id to redis failed,traceID:%s,err:%v", traceID, err)
	}
	return consts.UnSubscribeRespContent, nil
}

func (a *WXRepository) handlerTEMPLATESENDJOBFINISHEvent(ctx context.Context, reqBody *entity.TextRequestBody) (string, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("handlerTEMPLATESENDJOBFINISHEvent traceID:%s", traceID)
	// 对事件推送由于其他原因发送失败的消息进行重发
	// todo: 引入消息队列重发，定时任务扫描 发送中 状态的消息，定时任务补偿3次，还没成功状态设为失败，人工介入，消费者根据消息id做幂等处理

	return "", nil
}

func (a *WXRepository) isExistUserMsgID(ctx context.Context, msgID string, fromUserName string, createTime int64) (bool, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("IsExistUserMsgID traceID:%s", traceID)
	exist, err := a.wx.IsExistMsgIDFromRedis(ctx, msgID)
	if err != nil {
		log.Errorf("handlerSubscribeEvent IsExistMsgIDFromRedis failed,traceID:%s,err:%v", traceID, err)
		return false, err
	}
	// 若存在返回空串,不存在则持久化存储,并保存msg id 到 redis
	if exist {
		return true, nil
	}
	// 从db上找，存在则返回空串
	exist, err = a.user.IsExistUserFromDB(ctx, fromUserName, createTime)
	if err != nil {
		log.Errorf("handlerSubscribeEvent IsExistUserFromDB failed,traceID:%s,err:%v", traceID, err)
		return false, err
	}
	if exist {
		// 回写到redis中
		err = a.wx.SetMsgIDToRedis(ctx, msgID)
		if err != nil {
			log.Errorf("handlerSubscribeEvent WXRepository wx repo set msg id to redis failed,traceID:%s,err:%v", traceID, err)
		}
		return true, nil
	}
	return false, nil
}

func (a *WXRepository) makeTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	textResponseBody := &entity.TextResponseBody{}
	textResponseBody.FromUserName = a.value2CDATA(fromUserName)
	textResponseBody.ToUserName = a.value2CDATA(toUserName)
	textResponseBody.MsgType = a.value2CDATA("text")
	textResponseBody.Content = a.value2CDATA(content)
	textResponseBody.CreateTime = time.Now().Unix()
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}

func (a *WXRepository) value2CDATA(v string) entity.CDATAText {
	return entity.CDATAText{Text: "<![CDATA[" + v + "]]>"}
}
