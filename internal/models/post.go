package models

import (
	"encoding/json"
	"errors"
	"time"
)

// Post 帖子模型（注意内存对齐）
type Post struct {
	PostID      int64     `json:"post_id,string" db:"post_id"`
	AuthorID    int64     `json:"author_id" db:"author_id"`
	CommunityID int64     `json:"community_id" db:"community_id" binding:"required"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	Status      int8      `json:"status" db:"status"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
	UpdateTime  time.Time `json:"update_time" db:"update_time"`
}

// UnmarshalJSON 为POST类型实现自定义的 UnmarshalJSON 方法
func (p *Post) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Title       string `json:"title" db:"title"`
		Content     string `json:"content" db:"content"`
		CommunityID int64  `json:"community_id" db:"community_id"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.Title) == 0 {
		err = errors.New("帖子标题不能为空")
	} else if len(required.Content) == 0 {
		err = errors.New("帖子内容不能为空")
	} else if required.CommunityID == 0 {
		err = errors.New("未指定版块")
	} else {
		p.Title = required.Title
		p.Content = required.Content
		p.CommunityID = required.CommunityID
	}
	return
}

// ApiPostDetail 帖子详情模型
type ApiPostDetail struct {
	AuthorName       string             `json:"author_name"`   // 作者名
	AuthorAvatar     string             `json:"author_avatar"` // 头像相对路径
	VoteNum          int64              `json:"vote_num"`      // 投票数量
	CommentCount     int64              `json:"comment_count"` // 帖子评论的数量
	*Post                               // 嵌入帖子结构体
	*CommunityDetail `json:"community"` // 嵌入社区结构体
}

// ApiPostDetailRes 搜索帖子返回模型
type ApiPostDetailRes struct {
	Page Page             `json:"page"`
	List []*ApiPostDetail `json:"list"`
}
