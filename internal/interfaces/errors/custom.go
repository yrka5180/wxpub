package errors

import (
	"github.com/pkg/errors"
)

const (
	XCode    = "x-code"
	XMsg     = "x-msg"
	XTraceID = "x-nova-trace-id"
)

type CustomError struct {
	Cause     error  // 真正错误原因
	ErrorCode int    // 错误码
	ErrorMsg  string // 错误信息
}

func NewCustomError(err error, code int, msg string) error {
	return CustomError{
		Cause:     errors.Wrap(err, msg),
		ErrorCode: code,
		ErrorMsg:  msg,
	}
}

func (err CustomError) Error() string {
	if err.Cause != nil {
		return err.Cause.Error()
	}
	return err.ErrorMsg
}
