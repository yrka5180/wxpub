package httputil

import (
	"context"
	"fmt"
	"strconv"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/errorx"

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

	c.Header(errorx.XCode, strconv.Itoa(resp.ErrMsg.XCode))
	c.Header(errorx.XMsg, resp.ErrMsg.XMsg)
	c.Header(errorx.XTraceID, utils.ShouldGetTraceID(ctx))

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
	if customErr, ok := err.(errorx.CustomError); ok {
		setResponse(resp, customErr.ErrorCode, customErr.ErrorMsg, nil)
		return
	}
	setResponse(resp, errorx.CodeInternalServerError, errorx.GetErrorMessage(errorx.CodeInternalServerError), nil)
}

func DefaultResponse() GenericResponse {
	return GenericResponse{
		ErrMsg: &ErrorMessage{
			XCode: errorx.CodeOK,
			XMsg:  errorx.GetErrorMessage(errorx.CodeOK),
		},
	}
}

func setResponse(resp *GenericResponse, errCode int, errMsg string, body interface{}) {
	resp.StatusCode = errorx.GetStatusCode(errCode)
	if body != nil {
		resp.Body = body
	} else {
		resp.ErrMsg.XCode = errCode
		if len(errMsg) == 0 {
			errMsg = errorx.GetErrorMessage(errCode)
		}
		resp.ErrMsg.XMsg = errMsg
	}
}

func Abort(ctx *gin.Context, errCode int, errMsg string, traceID string) {
	ctx.Header(errorx.XCode, fmt.Sprintf("%d", errCode))
	ctx.Header(errorx.XMsg, errMsg)
	ctx.Header(errorx.XTraceID, traceID)
	ctx.AbortWithStatus(errorx.GetStatusCode(errCode))
}
