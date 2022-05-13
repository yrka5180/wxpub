package controller

import (
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/application"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/errors"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/httputil"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/middleware"
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

// swagger:route POST /message/push/tmpl 消息推送 SendTmplMessage
//
// description: 模板消息推送
//
// responses:
//   200: APISendTmplMessage
//   400: badRequest
//   401: unauthorized
//   403: forbidden
//   404: notfound
//   409: conflict
//   500: serverError
func (a *Message) SendTmplMessage(c *gin.Context) {
	ctx := middleware.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil.DefaultResponse()
	defer httputil.HTTPJSONResponse(ctx, c, &resp)

	var param entity.SendTmplMsgReq
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Errorf("validate sendmsg req ShouldBindJSON failed, traceID:%s, err:%v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInvalidParams, "Invalid json provided")
		return
	}
	errMsg := param.Validate()
	if len(errMsg) > 0 {
		log.Errorf("validate sendmsg req param failed, traceID:%s, errMsg:%s", traceID, errMsg)
		httputil.SetErrorResponse(&resp, errors.CodeInvalidParams, errMsg)
		return
	}
	msgResp, err := a.message.SendTmplMsg(ctx, param)
	if err != nil {
		log.Errorf("SendMessage MessageInterface send msg failed,traceID:%s,err:%v", traceID, err)
		httputil.SetErrorResponse(&resp, errors.CodeInternalServerError, err.Error())
		return
	}
	httputil.SetSuccessfulResponse(&resp, errors.CodeOK, msgResp)
}
