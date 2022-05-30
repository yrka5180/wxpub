package swagger

type BasicError struct {
	// 错误信息，用于前端展示或者辅助定位错误源
	XCode int    `json:"x-code"` // x-code
	XMsg  string `json:"x-msg"`  // x-msg
}

type BasicRender struct {
	// 统一返回体
	// data struct
	Dat interface{} `json:"dat"` // data struct
	// err msg
	Err string `json:"err"` // err msg
}

// SuccessResp 操作成功(不需要返回内容，比如删除)
// swagger:response success
type SuccessResp struct {
	// in: header
	BasicError
	// in: body
	Body struct {
		BasicRender
	}
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
	// in: body
	Body struct {
		BasicRender
	}
}

// ErrRedirect 重定向
// swagger:response redirect
type ErrRedirect struct {
	// in: header
	BasicError
	// in: body
	Body struct {
		BasicRender
	}
}

// ErrExpectationFailed 请求异常，请求格式错误或者参数异常
// swagger:response expectationFailed
type ErrExpectationFailed struct {
	// in: header
	BasicError
	// in: body
	Body struct {
		BasicRender
	}
}

// ErrForbidden 没有权限访问
// swagger:response forbidden
type ErrForbidden struct {
	// in: header
	BasicError
	// in: body
	Body struct {
		BasicRender
	}
}

// ErrUnauthorized 无法识别用户身份，令牌过期或者令牌被篡改
// swagger:response unauthorized
type ErrUnauthorized struct {
	// in: header
	BasicError
	// in: body
	Body struct {
		BasicRender
	}
}

// ErrNotFound 找不到该资源
// swagger:response notfound
type ErrNotFound struct {
	// in: header
	BasicError
	// in: body
	Body struct {
		BasicRender
	}
}

// ErrConflict 资源状态冲突
// swagger:response conflict
type ErrConflict struct {
	// in: header
	BasicError
	// in: body
	Body struct {
		BasicRender
	}
}

// ErrServerInternal 服务器内部错误
// swagger:response serverError
type ErrServerInternal struct {
	// in: header
	BasicError
	// in: body
	Body struct {
		BasicRender
	}
}
