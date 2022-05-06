package entity

import (
	"encoding/xml"
	"time"
)

type WXCheckReq struct {
	Signature string `json:"signature" form:"signature"`
	TimeStamp string `json:"time_stamp" form:"timestamp"`
	Nonce     string `json:"nonce" form:"nonce"`
	EchoStr   string `json:"echo_str" form:"echostr"`
}

type TextRequestBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   time.Duration
	MsgType      string
	Content      string
	MsgId        int
}

type TextResponseBody struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   CDATAText
	FromUserName CDATAText
	CreateTime   time.Duration
	MsgType      CDATAText
	Content      CDATAText
}

type CDATAText struct {
	Text string `xml:",innerxml"`
}

func (u *WXCheckReq) Validate() (errorMessage string) {
	return
}
