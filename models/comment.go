package models

import "time"

// Comment 评论模型
type Comment struct {
	Id         int64     `json:"id,string" db:"id"`
	CommentId  int64     `json:"comment_id,string" db:"comment_id"`
	ParentId   int64     `json:"parent_id,string" db:"parent_id"` // 父评论id，如果是0表示是一级评论
	PostId     int64     `json:"post_id,string" db:"post_id"`
	AuthorId   int64     `json:"author_id,string" db:"author_id"`
	ReplyToUid int64     `json:"reply_to_uid,string" db:"reply_to_uid"` // 新增:被回复人的用户id
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

// ApiCommentDetail 评论详情
type ApiCommentDetail struct {
	CommentId     int64  `json:"comment_id,string"`
	ParentId      int64  `json:"parent_id,string"`
	PostId        int64  `json:"post_id,string"`
	AuthorId      int64  `json:"author_id,string"`
	AuthorName    string `json:"author_name"`
	AuthorAvatar  string `json:"author_avatar"`       // 评论作者头像
	ReplyToUid    int64  `json:"reply_to_uid,string"` // 被回复人ID
	ReplyToName   string `json:"reply_to_name"`       // 被回复人用户名
	ReplyToAvatar string `json:"reply_to_avatar"`     // 被回复人头像
	ReplyCount    int64  `json:"reply_count"`
	VoteNum       int64  `json:"vote_num"`
	Content       string `json:"content"`
	CreateTime    string `json:"create_time"`
}

// ParamCommentReply 创建评论回复的请求参数
type ParamCommentReply struct {
	ParentId   string `json:"parent_id" binding:"required"`              // 父评论id
	PostId     string `json:"post_id" binding:"required"`                // 帖子id
	ReplyToUid string `json:"reply_to_uid" binding:"required"`           // 被回复人的用户id
	Content    string `json:"content" binding:"required,min=1,max=1000"` // 修改最大长度限制
}

// ApiCommentListRes 评论列表接口响应数据
type ApiCommentListRes struct {
	Page *Page               `json:"page"` // 分页信息
	List []*ApiCommentDetail `json:"list"` // 评论列表
}
