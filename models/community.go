package models

import "time"

// Community 社区列表model
type Community struct {
	CommunityID   int64  `json:"community_id" db:"community_id"`
	CommunityName string `json:"community_name" db:"community_name"`
}

// CommunityDetail 社区详情model
type CommunityDetail struct {
	CommunityID   int64     `json:"community_id" db:"community_id"`
	CommunityName string    `json:"community_name" db:"community_name"`
	Introduction  string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime    time.Time `json:"create_time" db:"create_time"`
}
