package middleware

import (
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/httputil"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/application"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/errors"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Passport struct {
	passportApp application.PassportAppInterface
}

func NewPassportMiddleware(passportApp application.PassportAppInterface) *Passport {
	return &Passport{
		passportApp: passportApp,
	}
}

func (p *Passport) VerifyToken(c *gin.Context) {
	ctx := DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	auth := c.GetHeader(consts.Authorization)
	if len(auth) <= 0 {
		log.Errorf("verifyToken oidc introspect failed,traceID:%s", traceID)
		httputil.Abort(c, errors.CodeUnauthorized, errors.GetErrorMessage(errors.CodeUnauthorized), traceID)
		return
	}
	err := p.passportApp.GetAuthN(ctx, auth)
	if err != nil {
		log.Errorf("VerifyToken GetAuthN failed,traceID:%s,err:%v", traceID, err)
		httputil.Abort(c, err.(errors.CustomError).ErrorCode, err.(errors.CustomError).ErrorMsg, traceID)
		return
	}
	c.Next()
}
