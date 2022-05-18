package swagger

import "git.nova.net.cn/nova/misc/wx-public/proxy/internal/domain/entity"

// swagger:response APIListUser
type APIListUser struct {
	// in: body
	Data []entity.User `json:"data"`
}
