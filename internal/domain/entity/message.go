package entity

type SendTmplMsgReq struct {
	// 获取到的凭证
	AccessToken string `json:"access_token" form:"access_token"`
	// 接收者openid
	ToUser string `json:"touser"`
	// 模板ID
	TemplateID string `json:"template_id"`
	// 模板数据
	Data interface{} `json:"data"`
}

type SendTmplMsgResp struct {
	MsgID int64 `json:"msgid"`
	ErrorInfo
}

func (r *SendTmplMsgReq) Validate() (errorMsg string) {
	if len(r.ToUser) <= 0 {
		errorMsg = "toUser is empty"
		return
	}
	if len(r.TemplateID) <= 0 {
		errorMsg = "template id is empty"
	}
	return
}
