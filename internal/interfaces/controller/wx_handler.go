package controller

import (
	"public-platform-manager/internal/application"
	"public-platform-manager/internal/consts"
	"public-platform-manager/internal/domain/entity"
	"public-platform-manager/internal/interfaces/errors"
	"public-platform-manager/internal/interfaces/httputil"
	"public-platform-manager/internal/interfaces/middleware"
	"public-platform-manager/internal/utils"

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
	log.Infof("wx public platform access successfully!")
}

func (a *WX) GetEventXml(c *gin.Context) {
	ctx := middleware.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil.DefaultResponse()
	defer httputil.HTTPResponse(ctx, c, &resp)
	var param entity.WXCheckReq
	if err := c.ShouldBindQuery(&param); err != nil {
		log.Errorf("validate GetEventXml ShouldBindQuery failed, traceID:%s, err:%v", traceID, err)
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
		log.Errorf("validate GetEventXml ShouldBindXml failed, traceID:%s, err:%v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInvalidParams, "Invalid xml body provided")
		return
	}
	// 事件xml返回
	respBody, err := a.wx.GetEventXml(ctx, reqBody)
	if err != nil {
		log.Errorf("wx public platform GetEventXml access failed,traceID:%s,err:%v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInternalServerError, "wx public platform GetEventXml access failed!")
		return
	}
	// 原样返回
	httputil.SetSuccessfulResponse(&resp, errors.CodeOK, string(respBody))
	log.Infof("wx public platform access successfully!")
}
