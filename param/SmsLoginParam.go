package param

// 登录方式：手机号+短信登录 参数结构体
type SmsLogin struct {
	Phone string `json:"phone"`
	Code string `json:"code"`
}