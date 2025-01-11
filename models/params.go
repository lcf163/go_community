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
	TargetId   string `json:"target_id" binding:"required"`              // 投票目标ID
	TargetType int8   `json:"target_type" binding:"required,oneof=1 2"`  // 投票目标类型(1:帖子 2:评论)
	Direction  int8   `json:"direction" binding:"required,oneof=1 0 -1"` // 赞成票(1)、取消投票(0)、反对票(-1)
}

func (v *ParamVoteData) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		TargetId   string `json:"target_id"`
		TargetType int8   `json:"target_type"`
		Direction  int8   `json:"direction"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.TargetId) == 0 {
		err = errors.New("缺少必填字段target_id")
	} else if required.TargetType == 0 {
		err = errors.New("缺少必填字段target_type")
	} else if required.Direction == 0 {
		err = errors.New("缺少必填字段direction")
	} else {
		v.TargetId = required.TargetId
		v.TargetType = required.TargetType
		v.Direction = required.Direction
	}
	return
}

// ParamPostList 获取帖子列表 query string 参数
type ParamPostList struct {
	Search      string `json:"search" form:"search"`               // 关键字搜索
	CommunityId int64  `json:"community_id" form:"community_id"`   // 可以为空
	Page        int64  `json:"page" form:"page"`                   // 页码
	Size        int64  `json:"size" form:"size"`                   // 每页数量
	Order       string `json:"order" form:"order" example:"score"` // 排序依据
}

// ParamPage 分页参数
type ParamPage struct {
	Page int64 `json:"page" form:"page"` // 页码
	Size int64 `json:"size" form:"size"` // 每页数量
}

// Page 分页结构体
type Page struct {
	Total int64 `json:"total"` // 总数
	Page  int64 `json:"page"`  // 页码
	Size  int64 `json:"size"`  // 每页数量
}
