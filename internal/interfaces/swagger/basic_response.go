package swagger

type BasicError struct {
	// 错误信息，用于前端展示或者辅助定位错误源
	// in: header
	XCode int    `json:"x-code"` // x-code
	XMsg  string `json:"x-msg"`  // x-msg
}

// SuccessResp 操作成功(不需要返回内容，比如删除)
// swagger:response success
type SuccessResp struct {
	// in: header
	BasicError
}

// NoContentResp 操作成功(不需要返回内容，比如删除)
// swagger:response noContent
type NoContentResp struct {
}

// ErrBadRequest 请求异常，请求格式错误或者参数异常
// swagger:response badRequest
type ErrBadRequest struct {
	// in: header
	BasicError
}

// ErrRedirect 重定向
// swagger:response redirect
type ErrRedirect struct {
	// in: header
	BasicError
}

// ErrExpectationFailed 请求异常，请求格式错误或者参数异常
// swagger:response expectationFailed
type ErrExpectationFailed struct {
	// in: header
	BasicError
}

// ErrForbidden 没有权限访问
// swagger:response forbidden
type ErrForbidden struct {
	// in: header
	BasicError
}

// ErrUnauthorized 无法识别用户身份，令牌过期或者令牌被篡改
// swagger:response unauthorized
type ErrUnauthorized struct {
	// in: header
	BasicError
}

// ErrNotFound 找不到该资源
// swagger:response notfound
type ErrNotFound struct {
	// in: header
	BasicError
}

// ErrConflict 资源状态冲突
// swagger:response conflict
type ErrConflict struct {
	// in: header
	BasicError
}

// ErrServerInternal 服务器内部错误
// swagger:response serverError
type ErrServerInternal struct {
	// in: header
	BasicError
}
