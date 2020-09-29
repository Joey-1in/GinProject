package param

type LoginParam struct {
	Name string `json:"name"`
	Password string `json:"pwd"`
	Id       string `json:"id"`    // 验证码id
	Value    string `json:"value"` // 验证码输入值
}