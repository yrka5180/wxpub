package entity

type WXCheckReq struct {
	Signature string `json:"signature" form:"signature"`
	TimeStamp string `json:"time_stamp" form:"timestamp"`
	Nonce     string `json:"nonce" form:"nonce"`
	EchoStr   string `json:"echo_str" form:"echostr"`
}

func (u *WXCheckReq) Validate() (errorMessage string) {
	return
}
