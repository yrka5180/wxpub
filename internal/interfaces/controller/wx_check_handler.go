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

type WXCheckSignature struct {
	wx application.WXCheckSignatureInterface
}

func NewWXCheckSignController(awApp application.WXCheckSignatureInterface) *WXCheckSignature {
	return &WXCheckSignature{
		wx: awApp,
	}
}

func (a *WXCheckSignature) GetWXCheckSign(c *gin.Context) {
	ctx := middleware.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil.DefaultResponse()
	defer httputil.HTTPJSONResponse(ctx, c, &resp)
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
