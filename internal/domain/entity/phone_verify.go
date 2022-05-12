package entity

type SendSmsReq struct {
	// 微信用户的openID
	OpenID string `json:"open_id"`
	// 目标手机号
	Phone string `json:"phone"`
}

type SendSmsResp struct {
	// 微信用户的openID
	OpenID string `json:"open_id"`
	// 验证方式id，目前固定为短信"sms"
	VerifyCodeID string `json:"verify_code_id"`
	ErrorInfo
}

type VerifyCodeRedisValue struct {
	VerifyCodeAnswer     string `json:"verify_code_answer"`
	VerifyCodeCreateTime int64  `json:"verify_code_create_time"`
}
