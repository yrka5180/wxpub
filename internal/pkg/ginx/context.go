package ginx

import (
	"context"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/consts"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

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
