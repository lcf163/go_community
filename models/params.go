package models

import (
	"encoding/json"
	"errors"
	"strconv"
)

// 定义请求参数的结构体

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// parseID 解析ID，支持字符串和数字类型
func parseID(v interface{}) (int64, error) {
	switch v := v.(type) {
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case int64:
		return v, nil
	case json.Number:
		return v.Int64()
	default:
		return 0, errors.New("invalid id type")
	}
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

// ParamUpdateUser 修改用户名的参数
type ParamUpdateUser struct {
	Username string `json:"username" binding:"required"` // 用户名,必填
}

// ParamUpdatePassword 修改密码的参数
type ParamUpdatePassword struct {
	OldPassword string `json:"old_password" binding:"required"` // 旧密码
	NewPassword string `json:"new_password" binding:"required"` // 新密码
}

// ParamPost 创建帖子请求参数
type ParamPost struct {
	CommunityID int64  `json:"community_id" binding:"required"` // 社区ID
	Title       string `json:"title" binding:"required"`        // 标题
	Content     string `json:"content" binding:"required"`      // 内容
}

// UnmarshalJSON 自定义反序列化方法
func (p *ParamPost) UnmarshalJSON(data []byte) error {
	tmp := struct {
		CommunityID interface{} `json:"community_id"`
		Title       string      `json:"title"`
		Content     string      `json:"content"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// 处理必填字段 CommunityID
	if tmp.CommunityID == nil {
		return errors.New("community_id is required")
	}
	communityID, err := parseID(tmp.CommunityID)
	if err != nil {
		return errors.New("invalid community_id")
	}
	p.CommunityID = communityID
	p.Title = tmp.Title
	p.Content = tmp.Content
	return nil
}

// ParamPostListQueryNoSearch 获取帖子列表的请求参数（不搜索关键词）
type ParamPostListQueryNoSearch struct {
	CommunityID int64  `json:"community_id" form:"community_id"`   // 可以为空
	Page        int64  `json:"page" form:"page"`                   // 页码
	Size        int64  `json:"size" form:"size"`                   // 每页数量
	Order       string `json:"order" form:"order" example:"score"` // 排序依据
}

// ParamPostListQueryWithSearch 获取帖子列表的请求参数（不搜索社区ID）
type ParamPostListQueryWithSearch struct {
	Page   int64  `json:"page" form:"page"`     // 页码
	Size   int64  `json:"size" form:"size"`     // 每页数量
	Search string `json:"search" form:"search"` // 关键字搜索
}

// ParamPostList 获取帖子列表的请求参数
type ParamPostList struct {
	CommunityId int64  `json:"community_id" form:"community_id"`   // 可以为空
	Page        int64  `json:"page" form:"page"`                   // 页码
	Size        int64  `json:"size" form:"size"`                   // 每页数量
	Order       string `json:"order" form:"order" example:"score"` // 排序依据
	Search      string `json:"search" form:"search"`               // 关键字搜索
}

// ParamUpdatePost 更新帖子请求参数
type ParamUpdatePost struct {
	PostID  int64  `json:"post_id" binding:"required"` // 帖子id
	Title   string `json:"title" binding:"required"`   // 标题
	Content string `json:"content" binding:"required"` // 内容
}

// UnmarshalJSON 自定义反序列化方法
func (p *ParamUpdatePost) UnmarshalJSON(data []byte) error {
	tmp := struct {
		PostID  interface{} `json:"post_id"`
		Title   string      `json:"title"`
		Content string      `json:"content"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// 处理必填字段 PostId
	if tmp.PostID == nil {
		return errors.New("post_id is required")
	}
	postID, err := parseID(tmp.PostID)
	if err != nil {
		return errors.New("invalid post_id")
	}
	p.PostID = postID
	p.Title = tmp.Title
	p.Content = tmp.Content
	return nil
}

// ParamComment 创建评论/回复的请求参数
type ParamComment struct {
	PostID     int64  `json:"post_id" binding:"required"`                // 帖子id
	ParentID   int64  `json:"parent_id"`                                 // 父评论id（0表示创建评论，非0表示创建回复）
	ReplyToUID int64  `json:"reply_to_uid"`                              // 被回复人的用户id（parent_id不为0时必填）
	Content    string `json:"content" binding:"required,min=1,max=1000"` // 内容
}

// UnmarshalJSON 自定义反序列化方法
func (p *ParamComment) UnmarshalJSON(data []byte) error {
	tmp := struct {
		PostID     interface{} `json:"post_id"`
		ParentID   interface{} `json:"parent_id"`
		ReplyToUID interface{} `json:"reply_to_uid"`
		Content    string      `json:"content"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// 处理必填字段 PostID
	if tmp.PostID == nil {
		return errors.New("post_id is required")
	}
	postID, err := parseID(tmp.PostID)
	if err != nil {
		return errors.New("invalid post_id")
	}
	p.PostID = postID

	// 处理可选字段 ParentID
	if tmp.ParentID != nil {
		parentID, err := parseID(tmp.ParentID)
		if err != nil {
			return errors.New("invalid parent_id")
		}
		p.ParentID = parentID
	}

	// 处理可选字段 ReplyToUID
	if tmp.ReplyToUID != nil {
		replyToUID, err := parseID(tmp.ReplyToUID)
		if err != nil {
			return errors.New("invalid reply_to_uid")
		}
		p.ReplyToUID = replyToUID
	}

	p.Content = tmp.Content
	return nil
}

//// ParamComment 创建评论请求参数
//type ParamComment struct {
//	PostId   int64  `json:"post_id" binding:"required"` // 帖子id
//	ParentId int64  `json:"parent_id"`                  // 父评论id，可选
//	Content  string `json:"content" binding:"required"` // 评论内容
//}
//
//// ParamCommentReply 创建评论回复的请求参数
//type ParamCommentReply struct {
//	ParentId   int64  `json:"parent_id" binding:"required"`              // 父评论id
//	PostId     int64  `json:"post_id" binding:"required"`                // 帖子id
//	ReplyToUid int64  `json:"reply_to_uid" binding:"required"`           // 被回复人的用户id
//	Content    string `json:"content" binding:"required,min=1,max=1000"` // 修改最大长度限制
//}

// ParamUpdateComment 更新评论请求参数
type ParamUpdateComment struct {
	CommentID int64  `json:"comment_id" binding:"required"` // 评论id
	Content   string `json:"content" binding:"required"`    // 评论内容
}

// UnmarshalJSON 自定义反序列化方法
func (p *ParamUpdateComment) UnmarshalJSON(data []byte) error {
	tmp := struct {
		CommentID interface{} `json:"comment_id"`
		Content   string      `json:"content"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// 处理必填字段 CommentId
	if tmp.CommentID == nil {
		return errors.New("comment_id is required")
	}
	commentID, err := parseID(tmp.CommentID)
	if err != nil {
		return errors.New("invalid comment_id")
	}
	p.CommentID = commentID
	p.Content = tmp.Content
	return nil
}

//// ParamUpdateCommentReply 更新评论回复请求参数
//type ParamUpdateCommentReply struct {
//	CommentId int64  `json:"comment_id" binding:"required"` // 评论id
//	Content   string `json:"content" binding:"required"`    // 回复内容
//}

// ParamCommentList 获取评论列表的请求参数
type ParamCommentList struct {
	PostID    int64 `form:"post_id"`         // 帖子id,获取帖子评论时必填
	CommentID int64 `form:"comment_id"`      // 评论id,获取评论回复时必填
	Page      int64 `form:"page,default=1"`  // 页码
	Size      int64 `form:"size,default=10"` // 每页数量
}

// ParamUpdateCommunity 更新社区请求参数
type ParamUpdateCommunity struct {
	Name         string `json:"community_name" binding:"required"` // 评论id
	Introduction string `json:"introduction" binding:"required"`   // 评论内容
}

// ParamVoteData 投票数据
type ParamVoteData struct {
	TargetID   int64 `json:"target_id" binding:"required"`              // 投票目标ID
	TargetType int8  `json:"target_type" binding:"required,oneof=1 2"`  // 投票目标类型(1:帖子 2:评论)
	Direction  int8  `json:"direction" binding:"required,oneof=1 0 -1"` // 赞成票(1)、取消投票(0)、反对票(-1)
}

// UnmarshalJSON 自定义反序列化方法
func (p *ParamVoteData) UnmarshalJSON(data []byte) error {
	tmp := struct {
		TargetID   interface{} `json:"target_id"`
		TargetType int8        `json:"target_type"`
		Direction  int8        `json:"direction"`
	}{}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	// 处理必填字段 TargetID
	if tmp.TargetID == nil {
		return errors.New("缺少必填字段target_id")
	}
	targetID, err := parseID(tmp.TargetID)
	if err != nil {
		return errors.New("invalid target_id")
	}
	p.TargetID = targetID

	// 处理必填字段 TargetType
	if tmp.TargetType == 0 {
		return errors.New("缺少必填字段target_type")
	}
	if tmp.TargetType != 1 && tmp.TargetType != 2 {
		return errors.New("target_type必须是1或2")
	}
	p.TargetType = tmp.TargetType

	// 处理必填字段 Direction
	if tmp.Direction == 0 {
		return errors.New("缺少必填字段direction")
	}
	if tmp.Direction != 1 && tmp.Direction != -1 {
		return errors.New("direction必须是1或-1")
	}
	p.Direction = tmp.Direction

	return nil
}
