package models

import (
	"encoding/json"
	"errors"
	"go_community/pkg/file"
)

// User 定义请求参数结构体
type User struct {
	UserId       int64  `json:"user_id,string" db:"user_id"`
	UserName     string `json:"username" db:"username"`
	Password     string `json:"password" db:"password"`
	Avatar       string `json:"avatar" db:"avatar"` // 头像相对路径
	AccessToken  string
	RefreshToken string
}

// UnmarshalJSON 为User类型实现自定义的UnmarshalJSON方法
func (u *User) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		UserName string `json:"username" db:"username"`
		Password string `json:"password" db:"password"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.UserName) == 0 {
		err = errors.New("缺少必填字段username")
	} else if len(required.Password) == 0 {
		err = errors.New("缺少必填字段password")
	} else {
		u.UserName = required.UserName
		u.Password = required.Password
	}
	return
}

// GetAvatarURL 获取头像完整路径
func (u *User) GetAvatarURL() string {
	return file.GetAvatarPath(u.Avatar)
}
