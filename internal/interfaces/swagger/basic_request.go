package swagger

// swagger:parameters SendSms VerifyAndUpdatePhone GenCaptcha SendTmplMessage TmplMsgStatus
type BasicParam struct {
	// 所有请求必须携带一个随机产生的 id，采用标准 UUID 格式，如：123e4567-e89b-12d3-a456-426655440000
	// in: header
	TraceID string `json:"x-nova-trace-id"`

	// 超时时间，以秒为单位，如：10
	// in: header
	Timeout int `json:"x-nova-timeout"`
}

// swagger:parameters TmplMsgStatus
type APIPathIDParam struct {
	// in: path
	ID int `json:"id" validate:"required"`
}
