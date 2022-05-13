package swagger

import "git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"

// swagger:response APISendTmplMessage
type APISendTmplMessage struct {
	// in: body
	Data entity.SendTmplMsgResp `json:"data"`
}
