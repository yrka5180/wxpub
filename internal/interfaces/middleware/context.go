package middleware

import (
	"context"
	"public-platform-manager/internal/consts"
	"public-platform-manager/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func NovaContext(ctx *gin.Context) {
	traceID := ctx.GetHeader(consts.HTTPTraceIDHeader)
	timeoutStr := ctx.GetHeader(consts.HTTPTimeoutHeader)
	if !validTraceID(traceID) {
		var err error
		log.Infof("Request %s doesn't input a trace id", ctx.Request.URL.Path)
		traceID, err = utils.GetUUID()
		if err != nil {
			log.Errorf("Request %s new uuid failed, err:%s", ctx.Request.URL.Path, err.Error())
		}
	}
	timeoutSec, _ := strconv.Atoi(timeoutStr)
	if timeoutSec < 1 || timeoutSec > consts.DefaultHTTPTimeOut {
		log.Infof("Request %s doesn't input a timeout argument or it's invalid: %s", ctx.Request.URL.Path, timeoutStr)
		timeoutSec = consts.DefaultHTTPTimeOut
	}

	c, cancelF := context.WithTimeout(context.WithValue(context.Background(), consts.ContextTraceID, traceID), time.Second*time.Duration(timeoutSec))
	defer cancelF()
	ctx.Set(consts.GinContextContext, c)
	ctx.Next()
}

func defaultNovaContext(defaultContext context.Context, ctx *gin.Context) (c context.Context) {
	c = defaultContext
	it, b := ctx.Get(consts.GinContextContext)
	if !b {
		log.Warnln("nova-context doesn't exists")
		return
	}
	if c, b = it.(context.Context); !b {
		log.Warnln("invalid nova-context value type")
		c = defaultContext
		return
	}
	return
}

func DefaultTodoNovaContext(ctx *gin.Context) context.Context {
	return defaultNovaContext(context.TODO(), ctx)
}

func validTraceID(id string) bool {
	return len(id) > 0
}
