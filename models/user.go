package models

import (
	"encoding/json"
	"errors"
	"go_community/pkg/file"
	"strings"
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
	// 如果头像为空，返回随机默认头像
	if u.Avatar == "" {
		return file.GetAvatarPath("")
	}
	// 如果头像已经是完整URL（以http开头），直接返回
	if strings.HasPrefix(u.Avatar, "http") {
		return u.Avatar
	}
	// 返回完整的URL路径，此时 u.Avatar 只包含文件名
	return file.GetAvatarPath(u.Avatar)
}

// ParamUpdateUser 修改用户名的参数
type ParamUpdateUser struct {
	Username string `json:"username" binding:"required"` // 用户名,必填
}

// ParamUpdatePassword 修改密码的参数
type ParamUpdatePassword struct {
	OldPassword string `json:"old_password" binding:"required"` // 旧密码
	NewPassword string `json:"new_password" binding:"required"` // 新密码
}
