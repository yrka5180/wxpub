package ginx

import (
	"net/http"
	"strconv"

	"git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/errorx"
	"github.com/gin-gonic/gin"
)

func BindJSON(c *gin.Context, param interface{}) {
	err := c.ShouldBindJSON(param)
	if err != nil {
		errorx.BombErr(http.StatusBadRequest, "bin json body invalid,err:%v", err)
	}
}

func URLParamStr(c *gin.Context, field string) string {
	val := c.Param(field)
	if val == "" {
		errorx.BombErr(http.StatusBadRequest, "url param[%s] is empty", field)
	}
	return val
}

func URLParamInt64(c *gin.Context, field string) int64 {
	strVal := URLParamStr(c, field)
	intVal, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		errorx.BombErr(http.StatusBadRequest, "cannot convert %s to int64", intVal)
	}
	return intVal
}

func URLParamInt(c *gin.Context, field string) int {
	return int(URLParamInt64(c, field))
}

func QueryStr(c *gin.Context, key string, defaultVal ...string) string {
	val := c.Query(key)
	if val != "" {
		return val
	}
	if len(defaultVal) == 0 {
		errorx.BombErr(http.StatusBadRequest, "query param[%s] is necessary", key)
	}
	return defaultVal[0]
}

func QueryInt(c *gin.Context, key string, defaultVal ...int) int {
	strVal := c.Query(key)
	if strVal != "" {
		intVal, err := strconv.Atoi(strVal)
		if err != nil {
			errorx.BombErr(http.StatusBadRequest, "cannot convert [%s] to int", strVal)
		}
		return intVal
	}
	if len(defaultVal) == 0 {
		errorx.BombErr(http.StatusBadRequest, "query param[%s] is necessary", key)
	}
	return defaultVal[0]
}

func QueryInt64(c *gin.Context, key string, defaultVal ...int64) int64 {
	strVal := c.Query(key)
	if strVal != "" {
		intVal, err := strconv.ParseInt(strVal, 10, 64)
		if err != nil {
			errorx.BombErr(http.StatusBadRequest, "cannot convert [%s] to int64", strVal)
		}
		return intVal
	}
	if len(defaultVal) == 0 {
		errorx.BombErr(http.StatusBadRequest, "query param[%s] is necessary", key)
	}
	return defaultVal[0]
}

func QueryBool(c *gin.Context, key string, defaultVal ...bool) bool {
	strVal := c.Query(key)
	if strVal != "" {
		if strVal == "true" || strVal == "1" || strVal == "on" || strVal == "checked" || strVal == "yes" || strVal == "Y" {
			return true
		} else if strVal == "false" || strVal == "0" || strVal == "off" || strVal == "no" || strVal == "N" {
			return false
		} else {
			errorx.BombErr(http.StatusBadRequest, "unknown arg[%s] value: %s", key, strVal)
		}
	}
	if len(defaultVal) == 0 {
		errorx.BombErr(http.StatusBadRequest, "arg[%s] is necessary", key)
	}
	return defaultVal[0]
}
