package swagger

import "git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"

// swagger:parameters SendSms
type SendSmsRequest struct {
	// in: body
	Data entity.SendSmsReq `json:"data"`
}

type VerifyCodeRequest struct {
	// in: body
	Data entity.VerifyCodeReq `json:"data"`
}

// swagger:response APICaptchaResp
type APICaptchaResp struct {
	// in: body
	Body struct {
		Dat entity.CaptchaResp `json:"dat"`
		Err string             `json:"err"`
	}
}
