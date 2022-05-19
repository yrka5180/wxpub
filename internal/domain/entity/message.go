package entity

import (
	"encoding/json"
	"time"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/config"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
)

type SendTmplMsgReq struct {
	// 接收者手机号
	ToUsersPhone []string `json:"tousers_phone"`
	// 模板数据
	Data json.RawMessage `json:"data"`
}

type SendTmplMsgResp struct {
	// 发送消息id
	SendMsgID string `json:"send_msg_id"`
	// 发送失败手机号
	FailureSendPhones []string `json:"failure_send_phones"`
}

type SendTmplMsgRemoteReq struct {
	// 获取到的凭证
	AccessToken string `json:"access_token" form:"access_token"`
	// 接收者openid
	ToUser string `json:"touser"`
	// 模板ID
	TemplateID string `json:"template_id"`
	// 模板数据
	Data json.RawMessage `json:"data"`
}

type SendTmplMsgRemoteResp struct {
	MsgID int64 `json:"msgid"`
	ErrorInfo
}

type KafkaTmplMsg struct {
	SendTmplMsgRemoteReq
	// 发送消息id
	SendMsgID string `json:"send_msg_id"`
	// 接收消息时间
	AcceptedTime int64 `json:"accepted_time"`
	// 处理失败次数
	FailureCount int `json:"failure_count"`
}

// FailureMsgLog 消息发送失败日志表
type FailureMsgLog struct {
	ID int `json:"id" gorm:"id"`
	// 发送消息id
	SendMsgID string `json:"send_msg_id" gorm:"send_msg_id"`
	// 微信回调消息id
	MsgID int64 `json:"msg_id" gorm:"msg_id"`
	// 接收者openid
	ToUser string `json:"to_user" gorm:"to_user"`
	// 模板id
	TemplateID string `json:"template_id" gorm:"template_id"`
	// 模板内容
	Content json.RawMessage `json:"content" gorm:"content"`
	// 失败原因
	Cause string `json:"cause" gorm:"cause"`
	// 发送状态，1为正常，2为重试中，3为失败
	Status int `json:"status" gorm:"status"`
	// 发送次数
	Count int `json:"count" gorm:"count"`
	// 创建时间
	CreateTime int64 `json:"create_time" gorm:"create_time"`
	// 回调更新时间
	UpdateTime int64 `json:"update_time" gorm:"update_time"`
}

func (f FailureMsgLog) TableName() string {
	return "failure_msg_log"
}

func (r *SendTmplMsgReq) Validate() (errorMsg string) {
	if len(r.ToUsersPhone) <= 0 {
		errorMsg = "toUsersPhone is empty"
		return
	}
	// 去重
	r.ToUsersPhone = utils.RemoveStringRepeated(r.ToUsersPhone)
	return
}

func (r *SendTmplMsgReq) TransferPerSendTmplMsg(toUser string) SendTmplMsgRemoteReq {
	return SendTmplMsgRemoteReq{
		ToUser:     toUser,
		TemplateID: config.TmplMsgID,
		Data:       r.Data,
	}
}

func (f *FailureMsgLog) TransferSendTmplMsgRemoteReq() SendTmplMsgRemoteReq {
	return SendTmplMsgRemoteReq{
		ToUser:     f.ToUser,
		TemplateID: f.TemplateID,
		Data:       f.Content,
	}
}

func (r SendTmplMsgRemoteReq) TransferKafkaTmplReq(sendMsgID string) KafkaTmplMsg {
	return KafkaTmplMsg{
		SendTmplMsgRemoteReq: r,
		SendMsgID:            sendMsgID,
		AcceptedTime:         time.Now().Unix(),
		FailureCount:         1,
	}
}

func (k KafkaTmplMsg) TransferFailureMsgLog(errMsg string, sendCreateTime int64) FailureMsgLog {
	return FailureMsgLog{
		SendMsgID:  k.SendMsgID,
		ToUser:     k.ToUser,
		TemplateID: k.TemplateID,
		Content:    k.Data,
		Cause:      errMsg,
		Status:     consts.SendRetry,
		Count:      k.FailureCount,
		CreateTime: sendCreateTime,
	}
}
