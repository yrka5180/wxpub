package persistence

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/config"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/errors"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/httputil"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	log "github.com/sirupsen/logrus"
)

type MessageRepo struct {
}

func NewMessageRepo() *MessageRepo {
	return &MessageRepo{}
}

func (a *MessageRepo) SendTmplMsgFromRequest(ctx context.Context, param entity.SendTmplMsgReq) (entity.SendTmplMsgResp, error) {
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("SendMsgFromRequest traceID:%s", traceID)
	// 请求wx msg send
	bs, err := json.Marshal(param)
	if err != nil {
		log.Errorf("SendMsgFromRequest json marshal send msg req failed,traceID:%s,err:%v", traceID, err)
		return entity.SendTmplMsgResp{}, err
	}
	requestProperty := httputil.GetRequestProperty(http.MethodPost, config.WXMsgTmplSendURL+fmt.Sprintf("?access_token=%s", param.AccessToken),
		bs, make(map[string]string))
	statusCode, body, _, err := httputil.RequestWithContextAndRepeat(ctx, requestProperty, traceID)
	if err != nil {
		log.Errorf("request wx msg send failed, traceID:%s, error:%v", traceID, err)
		return entity.SendTmplMsgResp{}, err
	}
	if statusCode != http.StatusOK {
		log.Errorf("request wx msg send failed, statusCode:%d,traceID:%s, error:%v", statusCode, traceID, err)
		return entity.SendTmplMsgResp{}, err
	}
	var msgResp entity.SendTmplMsgResp
	err = json.Unmarshal(body, &msgResp)
	if err != nil {
		log.Errorf("get wx msg send failed by unmarshal, resp:%s, traceID:%s, err:%v", string(body), traceID, err)
		return entity.SendTmplMsgResp{}, err
	}
	// 获取失败
	if msgResp.ErrCode != errors.CodeOK {
		log.Errorf("get wx msg send failed,resp:%s,traceID:%s,errMsg:%s", string(body), traceID, msgResp.ErrMsg)
		return entity.SendTmplMsgResp{}, fmt.Errorf("get wx msg send failed,errMsg:%s", msgResp.ErrMsg)
	}
	return msgResp, nil
}
