package models

import (
	"encoding/json"
	"errors"
)

// 定义请求参数的结构体

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	UserName        string `json:"username" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登录请求参数
type ParamLogin struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	PostId    string `json:"post_id" binding:"required"`              // 帖子Id
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票1、反对票-1、取消投票0
}

func (v *ParamVoteData) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		PostId    string `json:"post_id"`
		Direction int8   `json:"direction"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.PostId) == 0 {
		err = errors.New("缺少必填字段post_id")
	} else if required.Direction == 0 {
		err = errors.New("缺少必填字段direction")
	} else {
		v.PostId = required.PostId
		v.Direction = required.Direction
	}
	return
}

// ParamPostList 获取帖子列表 query string 参数
type ParamPostList struct {
	Page  int64  `json:"page" form:"page"`                   // 页码
	Size  int64  `json:"size" form:"size"`                   // 每页数量
	Order string `json:"order" form:"order" example:"score"` // 排序依据
}
