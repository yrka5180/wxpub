package controller

import (
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/application"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"
	errors2 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/errorx"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/ginx"
	httputil2 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/ginx/httputil"
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
	ctx := httputil2.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil2.DefaultResponse()
	defer httputil2.HTTPJSONResponse(ctx, c, &resp)

	var param entity.SendTmplMsgReq
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Errorf("SendTmplMessage validate sendmsg req ShouldBindJSON failed, traceID:%s, err:%+v", traceID, err)
		httputil2.SetErrorResponse(&resp, errors2.CodeInvalidParams, "Invalid json provided")
		return
	}
	// ginx.BindJSON(c,param)
	errMsg := param.Validate()
	if len(errMsg) > 0 {
		log.Errorf("SendTmplMessage validate sendmsg req param failed, traceID:%s, errMsg:%s", traceID, errMsg)
		httputil2.SetErrorResponse(&resp, errors2.CodeInvalidParams, errMsg)
		return
	}
	var msgResp entity.SendTmplMsgResp
	msgResp, err := a.message.SendTmplMsg(ctx, param)
	if err != nil {
		log.Errorf("SendTmplMessage MessageInterface send msg failed,traceID:%s,err:%+v", traceID, err)
		httputil2.SetErrorResponse(&resp, errors2.CodeInternalServerError, err.Error())
		return
	}
	httputil2.SetSuccessfulResponse(&resp, errors2.CodeOK, msgResp)
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
	ctx := httputil2.DefaultTodoNovaContext(c)
	traceID := utils.ShouldGetTraceID(ctx)
	log.Debugf("%s", traceID)

	resp := httputil2.DefaultResponse()
	defer httputil2.HTTPJSONResponse(ctx, c, &resp)

	var param entity.TmplMsgStatusReq
	param.RequestID = ginx.URLParamStr(c, "id")
	errMsg := param.Validate()
	if len(errMsg) > 0 {
		log.Errorf("TmplMsgStatus validate req param failed, traceID:%s, errMsg:%s", traceID, errMsg)
		httputil2.SetErrorResponse(&resp, errors2.CodeInvalidParams, errMsg)
		return
	}
	var msgStatusResp entity.TmplMsgStatusResp
	msgStatusResp, err := a.message.TmplMsgStatus(ctx, param.RequestID)
	if err != nil {
		log.Errorf("TmplMsgStatus MessageInterface tmpl msg status failed,traceID:%s,err:%+v", traceID, err)
		httputil2.SetErrorResponse(&resp, errors2.CodeInternalServerError, err.Error())
		return
	}
	httputil2.SetSuccessfulResponse(&resp, errors2.CodeOK, msgStatusResp)
}
