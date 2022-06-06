package main

import (
	"context"

	"git.nova.net.cn/nova/misc/wx-public/proxy/src/interfaces/webapi"
)

var (
	globalCtx    context.Context
	globalCancel context.CancelFunc
)

func main() {
	globalCtx, globalCancel = context.WithCancel(context.Background())
	webapi.Run(globalCtx, globalCancel)
}
