package models

// User 定义请求参数结构体
type User struct {
	UserID       int64  `json:"user_id" db:"user_id"`
	UserName     string `json:"username" db:"username"`
	Password     string `json:"password" db:"password"`
	AccessToken  string
	RefreshToken string
}
