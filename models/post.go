package models

import (
	"encoding/json"
	"errors"
	"time"
)

// Post 帖子model，注意内存对齐
type Post struct {
	PostId      int64     `json:"post_id,string" db:"post_id"`
	AuthorId    int64     `json:"author_id" db:"author_id"`
	CommunityId int64     `json:"community_id" db:"community_id" binding:"required"`
	Status      int32     `json:"status" db:"status"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Content     string    `json:"content" db:"content" binding:"required"`
	CreateTime  time.Time `json:"create_time" db:"create_time"`
	UpdateTime  time.Time `json:"-" db:"update_time"`
}

// UnmarshalJSON 为POST类型实现自定义的 UnmarshalJSON 方法
func (p *Post) UnmarshalJSON(data []byte) (err error) {
	required := struct {
		Title       string `json:"title" db:"title"`
		Content     string `json:"content" db:"content"`
		CommunityId int64  `json:"community_id" db:"community_id"`
	}{}
	err = json.Unmarshal(data, &required)
	if err != nil {
		return
	} else if len(required.Title) == 0 {
		err = errors.New("帖子标题不能为空")
	} else if len(required.Content) == 0 {
		err = errors.New("帖子内容不能为空")
	} else if required.CommunityId == 0 {
		err = errors.New("未指定版块")
	} else {
		p.Title = required.Title
		p.Content = required.Content
		p.CommunityId = required.CommunityId
	}
	return
}

// ApiPostDetail 帖子返回的详情model
type ApiPostDetail struct {
	AuthorName       string             `json:"author_name"`   // 作者名
	AuthorAvatar     string             `json:"author_avatar"` // 头像相对路径
	VoteNum          int64              `json:"vote_num"`      // 投票数量
	CommentCount     int64              `json:"comment_count"` // 帖子评论的数量
	*Post                               // 嵌入帖子结构体
	*CommunityDetail `json:"community"` // 嵌入社区结构体
}

// ApiPostDetailRes 搜索帖子返回的model
type ApiPostDetailRes struct {
	Page Page             `json:"page"`
	List []*ApiPostDetail `json:"list"`
}

// ParamPostList 获取帖子列表 query string 参数
type ParamPostList struct {
	CommunityId int64  `json:"community_id" form:"community_id"`   // 可以为空
	Page        int64  `json:"page" form:"page"`                   // 页码
	Size        int64  `json:"size" form:"size"`                   // 每页数量
	Order       string `json:"order" form:"order" example:"score"` // 排序依据
	Search      string `json:"search" form:"search"`               // 关键字搜索
}

// ParamUpdatePost 更新帖子请求参数
type ParamUpdatePost struct {
	PostId  int64  `json:"post_id" binding:"required"` // 帖子id
	Title   string `json:"title" binding:"required"`   // 标题
	Content string `json:"content" binding:"required"` // 内容
}
