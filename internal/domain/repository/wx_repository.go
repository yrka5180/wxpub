package repository

import (
	"context"
	"encoding/xml"
	"fmt"
	"public-platform-manager/internal/domain/entity"
	"public-platform-manager/internal/utils"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type WXRepository struct {
}

func NewWXRepository() *WXRepository {
	return &WXRepository{}
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
	log.Debugf("ListTemplateFromRequest traceID:%s", traceID)
	if reqBody == nil {
		return nil, fmt.Errorf("xml request body is empty")
	}
	responseTextBody, err := a.makeTextResponseBody(reqBody.ToUserName,
		reqBody.FromUserName,
		"Hello, "+reqBody.FromUserName)
	if err != nil {
		log.Errorf("Wechat Service: makeTextResponseBody traceID:%s,err:%v", traceID, err)
		return nil, nil
	}
	return responseTextBody, nil
}

func (a *WXRepository) value2CDATA(v string) entity.CDATAText {
	return entity.CDATAText{Text: "<![CDATA[" + v + "]]>"}
}

func (a *WXRepository) makeTextResponseBody(fromUserName, toUserName, content string) ([]byte, error) {
	textResponseBody := &entity.TextResponseBody{}
	textResponseBody.FromUserName = a.value2CDATA(fromUserName)
	textResponseBody.ToUserName = a.value2CDATA(toUserName)
	textResponseBody.MsgType = a.value2CDATA("text")
	textResponseBody.Content = a.value2CDATA(content)
	textResponseBody.CreateTime = time.Duration(time.Now().Unix())
	return xml.MarshalIndent(textResponseBody, " ", "  ")
}
