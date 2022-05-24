package middleware

import (
	"context"
	"strconv"
	"time"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"
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

func validTraceID(id string) bool {
	return len(id) > 0
}
