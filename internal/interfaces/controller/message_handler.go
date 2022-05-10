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

type Message struct {
	message application.MessageInterface
}

func NewMessageController(msg application.MessageInterface) *Message {
	return &Message{
		message: msg,
	}
}

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
