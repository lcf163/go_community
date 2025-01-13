package models

import "time"

// Comment 评论模型
type Comment struct {
	Id         int64     `json:"id,string" db:"id"`
	CommentId  int64     `json:"comment_id,string" db:"comment_id"`
	ParentId   int64     `json:"parent_id,string" db:"parent_id"` // 父评论id，如果是0表示是一级评论
	PostId     int64     `json:"post_id,string" db:"post_id"`
	AuthorId   int64     `json:"author_id,string" db:"author_id"`
	Content    string    `json:"content" db:"content"`
	Status     int8      `json:"status" db:"status"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
	UpdateTime time.Time `json:"update_time" db:"update_time"`
}

// ParamComment 创建评论请求参数
type ParamComment struct {
	PostId   string `json:"post_id" binding:"required"` // 帖子id
	ParentId string `json:"parent_id"`                  // 父评论id，可选
	Content  string `json:"content" binding:"required"` // 评论内容
}

// ApiCommentDetail 评论详情接口结构体
type ApiCommentDetail struct {
	CommentId  int64  `json:"comment_id,string" db:"comment_id"` // 评论id
	ParentId   int64  `json:"parent_id,string"`                  // 父评论id
	PostId     int64  `json:"post_id,string"`                    // 帖子id
	AuthorId   int64  `json:"author_id,string"`                  // 作者id
	ReplyCount int64  `json:"reply_count"`                       // 回复数量
	VoteNum    int64  `json:"vote_num"`                          // 投票数量
	AuthorName string `json:"author_name"`                       // 作者名字
	Content    string `json:"content"`                           // 评论内容
	CreateTime string `json:"create_time"`                       // 创建时间
}

// ParamCommentReply 创建评论回复的请求参数
type ParamCommentReply struct {
	ParentId string `json:"parent_id" binding:"required"`             // 父评论id
	PostId   string `json:"post_id" binding:"required"`               // 帖子id
	Content  string `json:"content" binding:"required,min=1,max=500"` // 回复内容
}
