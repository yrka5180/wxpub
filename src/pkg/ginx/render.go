package ginx

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Render struct {
	code int
	ctx  *gin.Context
}

func NewRender(c *gin.Context, code ...int) Render {
	r := Render{ctx: c}
	if len(code) > 0 {
		r.code = code[0]
	} else {
		r.code = 200
	}
	return r
}

func (r Render) Message(v interface{}, a ...interface{}) {
	if v == nil {
		if r.code == 200 {
			r.ctx.JSON(r.code, gin.H{"err": ""})
		} else {
			r.ctx.String(r.code, "")
		}
		return
	}

	switch t := v.(type) {
	case string:
		msg := fmt.Sprintf(t, a...)
		if r.code == 200 {
			r.ctx.JSON(r.code, gin.H{"err": msg})
		} else {
			r.ctx.String(r.code, msg)
		}
	case error:
		msg := fmt.Sprintf(t.Error(), a...)
		if r.code == 200 {
			r.ctx.JSON(r.code, gin.H{"err": msg})
		} else {
			r.ctx.String(r.code, msg)
		}
	}
}

func (r Render) Data(data interface{}, err interface{}, a ...interface{}) {
	if err == nil {
		r.ctx.JSON(r.code, gin.H{"dat": data, "err": ""})
		return
	}

	r.Message(err, a...)
}

func (r Render) DataString(data string, err interface{}, a ...interface{}) {
	if err == nil {
		r.ctx.String(r.code, data)
		return
	}

	r.Message(err, a...)
}

func (r Render) RawString(data string) {
	r.ctx.String(r.code, data)
}
