package controller

import (
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/application"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	errors2 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/errorx"
	httputil2 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/ginx/httputil"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type WX struct {
	wx application.WXInterface
}

func NewWXController(awApp application.WXInterface) *WX {
	return &WX{
		wx: awApp,
	}
}

func (a *WX) GetWXCheckSign(c *gin.Context) {
	ctx := httputil2.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil2.DefaultResponse()
	defer httputil2.HTTPResponse(ctx, c, &resp)
	var param entity.WXCheckReq
	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf("validate WXCheckReq ShouldBindQuery failed, traceID:%s, err:%+v", traceID, err)
		httputil2.SetErrorResponse(&resp, errors2.CodeInvalidParams, "Invalid query provided")
		return
	}
	// wx开放平台验证
	ok := a.wx.GetWXCheckSign(param.Signature, param.TimeStamp, param.Nonce, consts.Token)
	if !ok {
		log.Infof("wx public platform access failed!")
		return
	}
	// 原样返回
	httputil2.SetSuccessfulResponse(&resp, errors2.CodeOK, param.EchoStr)
}

func (a *WX) HandleEventXML(c *gin.Context) {
	ctx := httputil2.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil2.DefaultResponse()
	defer httputil2.HTTPResponse(ctx, c, &resp)
	var param entity.WXCheckReq
	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf("validate HandleEventXML ShouldBindQuery failed, traceID:%s, err:%+v", traceID, err)
		httputil2.SetErrorResponse(&resp, errors2.CodeInvalidParams, "Invalid query provided")
		return
	}
	// wx开放平台验证
	ok := a.wx.GetWXCheckSign(param.Signature, param.TimeStamp, param.Nonce, consts.Token)
	if !ok {
		log.Infof("wx public platform access failed!")
		return
	}
	var reqBody *entity.TextRequestBody
	if err := c.ShouldBindXML(&reqBody); err != nil {
		log.Errorf("validate HandleEventXML ShouldBindXML failed, traceID:%s, err:%+v", traceID, err)
		httputil2.SetErrorResponse(&resp, errors2.CodeInvalidParams, "Invalid xml body provided")
		return
	}
	// 事件xml返回
	respBody, err := a.wx.HandleEventXML(ctx, reqBody)
	if err != nil {
		log.Errorf("wx public platform HandleEventXML access failed,traceID:%s,err:%+v", traceID, err)
		httputil2.SetErrorResponse(&resp, errors2.CodeInternalServerError, "wx public platform GetEventXML access failed!")
		return
	}
	// 原样返回
	httputil2.SetSuccessfulResponse(&resp, errors2.CodeOK, string(respBody))
}
