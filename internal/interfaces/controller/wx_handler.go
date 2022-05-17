package controller

import (
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/application"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/errors"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/httputil"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/middleware"
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
	ctx := middleware.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil.DefaultResponse()
	defer httputil.HTTPResponse(ctx, c, &resp)
	var param entity.WXCheckReq
	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf("validate WXCheckReq ShouldBindQuery failed, traceID:%s, err:%v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInvalidParams, "Invalid query provided")
		return
	}
	// wx开放平台验证
	ok := a.wx.GetWXCheckSign(param.Signature, param.TimeStamp, param.Nonce, consts.Token)
	if !ok {
		log.Infof("wx public platform access failed!")
		return
	}
	// 原样返回
	httputil.SetSuccessfulResponse(&resp, errors.CodeOK, param.EchoStr)
}

func (a *WX) GetEventXML(c *gin.Context) {
	ctx := middleware.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil.DefaultResponse()
	defer httputil.HTTPResponse(ctx, c, &resp)
	var param entity.WXCheckReq
	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf("validate GetEventXML ShouldBindQuery failed, traceID:%s, err:%v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInvalidParams, "Invalid query provided")
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
		log.Errorf("validate GetEventXML ShouldBindXML failed, traceID:%s, err:%v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInvalidParams, "Invalid xml body provided")
		return
	}
	// 事件xml返回
	respBody, err := a.wx.GetEventXML(ctx, reqBody)
	if err != nil {
		log.Errorf("wx public platform GetEventXML access failed,traceID:%s,err:%v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInternalServerError, "wx public platform GetEventXML access failed!")
		return
	}
	// 原样返回
	httputil.SetSuccessfulResponse(&resp, errors.CodeOK, string(respBody))
}
