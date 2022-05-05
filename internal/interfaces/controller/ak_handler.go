package controller

import (
	"public-platform-manager/internal/application"
	"public-platform-manager/internal/domain/entity"
	"public-platform-manager/internal/interfaces/errors"
	"public-platform-manager/internal/interfaces/httputil"
	"public-platform-manager/internal/interfaces/middleware"
	"public-platform-manager/internal/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type AccessToken struct {
	ak application.AccessTokenInterface
}

func NewAccessTokenController(akApp application.AccessTokenInterface) *AccessToken {
	return &AccessToken{
		ak: akApp,
	}
}

func (a *AccessToken) GetAccessToken(c *gin.Context) {
	ctx := middleware.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil.DefaultResponse()
	defer httputil.HTTPJSONResponse(ctx, c, &resp)
	ak, err := a.ak.GetAccessToken(ctx)
	if err != nil {
		log.Errorf("GetAccessToken AccessTokenInterface get accss token failed,traceID:%s,err:%v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInternalServerError, err.Error())
		return
	}
	httputil.SetSuccessfulResponse(&resp, errors.CodeOK, entity.GetAccessTokenResp{
		AccessToken: ak,
	})
}

func (a *AccessToken) FreshAccessToken(c *gin.Context) {
	ctx := middleware.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil.DefaultResponse()
	defer httputil.HTTPJSONResponse(ctx, c, &resp)
	ak, err := a.ak.FreshAccessToken(ctx)
	if err != nil {
		log.Errorf("FreshAccessToken AccessTokenInterface fresh accss token failed,traceID:%s,err:%v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInternalServerError, err.Error())
		return
	}
	httputil.SetSuccessfulResponse(&resp, errors.CodeOK, entity.GetAccessTokenResp{
		AccessToken: ak,
	})
}
