package errorx

import (
	"fmt"
	"net/http"

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

func (err CustomError) String() string {
	return err.ErrorMsg
}

func BombErr(code int, format string, p ...interface{}) {
	panic(NewCustomError(nil, code, fmt.Sprintf(format, p...)))
}

func CustomErr(v interface{}, code ...int) {
	if v == nil {
		return
	}
	c := http.StatusOK
	if len(code) > 0 {
		c = code[0]
	}

	switch t := v.(type) {
	case string:
		if t != "" {
			panic(NewCustomError(nil, c, t))
		}
	case error:
		panic(NewCustomError(t, c, t.Error()))
	}
}
