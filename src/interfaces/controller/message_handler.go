package controller

import (
	"git.nova.net.cn/nova/misc/wx-public/proxy/src/application"
	"git.nova.net.cn/nova/misc/wx-public/proxy/src/domain/entity"
	errors2 "git.nova.net.cn/nova/misc/wx-public/proxy/src/pkg/errorx"
	"git.nova.net.cn/nova/misc/wx-public/proxy/src/pkg/ginx"
	"git.nova.net.cn/nova/misc/wx-public/proxy/src/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Message struct {
	message application.MessageInterface
}

func NewMessageController(msg application.MessageInterface) *Message {
	return &Message{
		message: msg,
	}
}

func (a *Message) SendTmplMessage(c *gin.Context) {
	ctx := ginx.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	var param entity.SendTmplMsgReq
	ginx.BindJSON(c, &param)
	errMsg := param.Validate()
	if len(errMsg) > 0 {
		log.Errorf("SendTmplMessage validate sendmsg req param failed, traceID:%s, errMsg:%s", traceID, errMsg)
		ginx.BombErr(errors2.CodeInvalidParams, errMsg)
	}
	var msgResp entity.SendTmplMsgResp
	msgResp, err := a.message.SendTmplMsg(ctx, param)
	ginx.NewRender(c).Data(msgResp, err)
}

func (a *Message) TmplMsgStatus(c *gin.Context) {
	ctx := ginx.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	var param entity.TmplMsgStatusReq
	param.RequestID = ginx.URLParamStr(c, "id")
	errMsg := param.Validate()
	if len(errMsg) > 0 {
		log.Errorf("TmplMsgStatus validate req param failed, traceID:%s, errMsg:%s", traceID, errMsg)
		ginx.BombErr(errors2.CodeInvalidParams, errMsg)
	}
	var msgStatusResp entity.TmplMsgStatusResp
	msgStatusResp, err := a.message.TmplMsgStatus(ctx, param.RequestID)
	ginx.NewRender(c).Data(msgStatusResp, err)
}