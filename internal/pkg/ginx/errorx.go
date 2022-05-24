package ginx

import "git.nova.net.cn/nova/misc/wx-public/proxy/internal/pkg/errorx"

func BombErr(code int, format string, p ...interface{}) {
	errorx.BombErr(code, format, p...)
}

func CustomErr(v interface{}, code ...int) {
	errorx.CustomErr(v, code...)
}
