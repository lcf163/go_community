package models

import (
	"time"
)

// Community 社区列表model
type Community struct {
	CommunityId   int64  `json:"community_id" db:"community_id"`
	CommunityName string `json:"community_name" db:"community_name"`
	Status        int8   `json:"status" db:"status"`
}

// CommunityDetail 社区详情model
type CommunityDetail struct {
	CommunityId   int64     `json:"community_id" db:"community_id"`
	CommunityName string    `json:"community_name" db:"community_name"`
	Introduction  string    `json:"introduction,omitempty" db:"introduction"`
	Status        int8      `json:"status" db:"status"`
	CreateTime    time.Time `json:"create_time" db:"create_time"`
}

// ApiCommunityDetailRes 社区列表接口响应数据
type ApiCommunityDetailRes struct {
	Page *Page              `json:"page"` // 分页信息
	List []*CommunityDetail `json:"list"` // 社区列表
}
