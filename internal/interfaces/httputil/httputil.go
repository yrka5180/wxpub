package httputil

import (
	"context"
	"fmt"
	"strconv"

	error2 "git.nova.net.cn/nova/misc/wx-public/proxy/internal/interfaces/errors"
	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type ErrorMessage struct {
	XCode int    `json:"x-code"` // x-code
	XMsg  string `json:"x-msg"`  // x-msg
}

type GenericResponse struct {
	StatusCode int
	ErrMsg     *ErrorMessage
	Body       interface{}
}

func (resp *GenericResponse) String() string {
	str := fmt.Sprintf("StatusCode:%d, ", resp.StatusCode)
	if resp.ErrMsg != nil {
		str += fmt.Sprintf("XCode:%+v, XMsg:%s, ", resp.ErrMsg.XCode, resp.ErrMsg.XMsg)
	}
	str += fmt.Sprintf("Body:%+v", resp.Body)
	return str
}

func HTTPJSONResponse(ctx context.Context, c *gin.Context, resp *GenericResponse) {
	log.Debugf("%s, resp: %+v", utils.ShouldGetTraceID(ctx), resp)

	c.Header(error2.XCode, strconv.Itoa(resp.ErrMsg.XCode))
	c.Header(error2.XMsg, resp.ErrMsg.XMsg)
	c.Header(error2.XTraceID, utils.ShouldGetTraceID(ctx))

	if resp.Body != nil {
		c.JSON(resp.StatusCode, resp.Body)
	} else {
		c.Status(resp.StatusCode)
	}
}

func HTTPResponse(ctx context.Context, c *gin.Context, resp *GenericResponse) {
	log.Debugf("%s, resp: %+v", utils.ShouldGetTraceID(ctx), resp)
	if resp.Body != nil {
		_, err := c.Writer.WriteString(resp.Body.(string))
		if err != nil {
			log.Errorf("HTTPResponse writer failed")
			return
		}
	} else {
		c.Status(resp.StatusCode)
	}
}

func SetSuccessfulResponse(resp *GenericResponse, code int, body interface{}) {
	setResponse(resp, code, "", body)
}

func SetErrorResponse(resp *GenericResponse, errCode int, errMsg string) {
	setResponse(resp, errCode, errMsg, nil)
}

func SetErrorResponseWithError(resp *GenericResponse, err error) {
	if customErr, ok := err.(error2.CustomError); ok {
		setResponse(resp, customErr.ErrorCode, customErr.ErrorMsg, nil)
		return
	}
	setResponse(resp, error2.CodeInternalServerError, error2.GetErrorMessage(error2.CodeInternalServerError), nil)
}

func DefaultResponse() GenericResponse {
	return GenericResponse{
		ErrMsg: &ErrorMessage{
			XCode: error2.CodeOK,
			XMsg:  error2.GetErrorMessage(error2.CodeOK),
		},
	}
}

func setResponse(resp *GenericResponse, errCode int, errMsg string, body interface{}) {
	resp.StatusCode = error2.GetStatusCode(errCode)
	if body != nil {
		resp.Body = body
	} else {
		resp.ErrMsg.XCode = errCode
		if len(errMsg) == 0 {
			errMsg = error2.GetErrorMessage(errCode)
		}
		resp.ErrMsg.XMsg = errMsg
	}
}

func Abort(ctx *gin.Context, errCode int, errMsg string, traceID string) {
	ctx.Header(error2.XCode, fmt.Sprintf("%d", errCode))
	ctx.Header(error2.XMsg, errMsg)
	ctx.Header(error2.XTraceID, traceID)
	ctx.AbortWithStatus(error2.GetStatusCode(errCode))
}
