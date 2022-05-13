package swagger

import "git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"

// swagger:response APIGetAccessTokenResp
type APIGetAccessTokenResp struct {
	// in: body
	Data entity.GetAccessTokenResp `json:"data"`
}
