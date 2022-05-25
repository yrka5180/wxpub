package controller

import (
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/application"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	errors2 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/errorx"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/ginx"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

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

// swagger:route POST /message/tmpl-push 消息推送 SendTmplMessage
//
// description: 模板消息推送
//
// responses:
//   200: APISendTmplMessage
//   400: badRequest
//   401: unauthorized
//   403: forbidden
//   404: notfound
//   409: APISendTmplMessage
//   500: serverError
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

// swagger:route GET /message/status/:id 消息推送 TmplMsgStatus
//
// description: 查看消息发送状态,资源id为request_id
//
// responses:
//   200: APITmplMsgStatusResp
//   400: badRequest
//   401: unauthorized
//   403: forbidden
//   404: notfound
//   409: APISendTmplMessage
//   500: serverError
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
