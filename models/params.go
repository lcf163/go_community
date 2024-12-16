package models

// 定义请求参数的结构体

type ParamSignUp struct {
	UserName        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}
