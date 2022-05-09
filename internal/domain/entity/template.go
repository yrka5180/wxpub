package entity

type ListTemplateReq struct {
	// 获取到的凭证
	AccessToken string `json:"access_token" form:"access_token"`
}

type ListTemplateResp struct {
	// 获取到的凭证
	TemplateList []TemplateItem `json:"template_list"`
	ErrorInfo
}

type TemplateItem struct {
	TemplateID      string `json:"template_id"`
	Title           string `json:"title"`
	PrimaryIndustry string `json:"primary_industry"`
	DeputyIndustry  string `json:"deputy_industry"`
	Content         string `json:"content"`
	Example         string `json:"example"`
}

func (u *ListTemplateReq) Validate() (errorMessage string) {
	errorMessage = ""
	if len(u.AccessToken) <= 0 {
		errorMessage = "access token is empty"
	}

	return
}
