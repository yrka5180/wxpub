package swagger

import "git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"

// swagger:parameters SendTmplMessage
type SendTmplMessageRequest struct {
	// in: body
	Data entity.SendTmplMsgReq `json:"data"`
}

// swagger:response APISendTmplMessage
type APISendTmplMessage struct {
	// in: body
	Body struct {
		Dat entity.SendTmplMsgResp `json:"dat"`
		Err string                 `json:"err"`
	}
}

// swagger:response APITmplMsgStatusResp
type APITmplMsgStatusResp struct {
	// in: body
	Body struct {
		Dat entity.TmplMsgStatusResp `json:"dat"`
		Err string                   `json:"err"`
	}
}
