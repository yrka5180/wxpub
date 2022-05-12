package entity

// AuthN 用于passport获取用户信息
type AuthN struct {
	Active    bool                   `json:"active"`
	Sub       string                 `json:"sub"`
	Username  string                 `json:"username"`
	Aud       []string               `json:"aud"`
	ClientID  string                 `json:"client_id"`
	Exp       int64                  `json:"exp"`
	Ext       map[string]interface{} `json:"ext"`
	Iat       int64                  `json:"iat"`
	Iss       string                 `json:"iss"`
	Nbf       int64                  `json:"nbf"`
	ObSub     string                 `json:"obfuscated_subject"`
	Scope     string                 `json:"scope"`
	TokenType string                 `json:"token_type"`
	TokenUse  string                 `json:"token_use"`
}
