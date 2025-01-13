package logic

import (
	"go_community/dao/mysql"
	"go_community/dao/redis"
	"go_community/models"
	"go_community/pkg/snowflake"
	"strconv"

	"go.uber.org/zap"
)

// CreateComment 创建评论
func CreateComment(userId int64, p *models.ParamComment) error {
	// 生成评论ID
	commentId := snowflake.GetID()

	// 转换帖子ID
	postId, err := strconv.ParseInt(p.PostId, 10, 64)
	if err != nil {
		return err
	}

	// 转换父评论ID
	var parentId int64
	if p.ParentId != "" {
		parentId, err = strconv.ParseInt(p.ParentId, 10, 64)
		if err != nil {
			return err
		}
	}

	comment := &models.Comment{
		CommentId: commentId,
		ParentId:  parentId,
		PostId:    postId,
		AuthorId:  userId,
		Content:   p.Content,
	}

	// 保存到MySQL
	if err := mysql.CreateComment(comment); err != nil {
		return err
	}

	// 保存到Redis
	return redis.CreateComment(commentId)
}

// GetCommentList 获取评论列表
func GetCommentList(postId int64) ([]*models.ApiCommentDetail, error) {
	// 查询评论列表
	comments, err := mysql.GetCommentList(postId)
	if err != nil {
		return nil, err
	}

	// 组装评论详情
	data := make([]*models.ApiCommentDetail, 0, len(comments))
	for _, comment := range comments {
		// 查询评论作者信息
		user, err := mysql.GetUserById(comment.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed",
				zap.Int64("author_id", comment.AuthorId),
				zap.Error(err))
			continue
		}

		// 获取回复数量
		replyCount, err := mysql.GetCommentReplyCount(comment.CommentId)
		if err != nil {
			zap.L().Error("mysql.GetCommentReplyCount failed",
				zap.Int64("comment_id", comment.CommentId),
				zap.Error(err))
			replyCount = 0
		}

		// 获取点赞数量
		voteNum, err := redis.GetCommentVoteNum(strconv.FormatInt(comment.CommentId, 10))
		if err != nil {
			zap.L().Error("redis.GetCommentVoteNum failed",
				zap.Int64("comment_id", comment.CommentId),
				zap.Error(err))
			voteNum = 0
		}

		commentDetail := &models.ApiCommentDetail{
			CommentId:  comment.CommentId,
			ParentId:   comment.ParentId,
			PostId:     comment.PostId,
			AuthorId:   comment.AuthorId,
			Content:    comment.Content,
			AuthorName: user.UserName,
			ReplyCount: replyCount,
			VoteNum:    voteNum,
			CreateTime: comment.CreateTime.Format("2006-01-02 15:04:05"),
		}
		data = append(data, commentDetail)
	}
	return data, nil
}

// GetCommentReplyList 获取评论的回复列表
func GetCommentReplyList(commentId int64) ([]*models.ApiCommentDetail, error) {
	// 查询回复列表
	comments, err := mysql.GetCommentReplyList(commentId)
	if err != nil {
		return nil, err
	}

	// 组装评论详情
	data := make([]*models.ApiCommentDetail, 0, len(comments))
	for _, comment := range comments {
		// 查询评论作者信息
		user, err := mysql.GetUserById(comment.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed",
				zap.Int64("author_id", comment.AuthorId),
				zap.Error(err))
			continue
		}

		// 获取回复数量
		replyCount, err := mysql.GetCommentReplyCount(comment.CommentId)
		if err != nil {
			zap.L().Error("mysql.GetCommentReplyCount failed",
				zap.Int64("comment_id", comment.CommentId),
				zap.Error(err))
			replyCount = 0
		}

		// 获取点赞数量
		voteNum, err := redis.GetCommentVoteNum(strconv.FormatInt(comment.CommentId, 10))
		if err != nil {
			zap.L().Error("redis.GetCommentVoteNum failed",
				zap.Int64("comment_id", comment.CommentId),
				zap.Error(err))
			voteNum = 0
		}

		commentDetail := &models.ApiCommentDetail{
			CommentId:  comment.CommentId,
			ParentId:   comment.ParentId,
			PostId:     comment.PostId,
			AuthorId:   comment.AuthorId,
			Content:    comment.Content,
			AuthorName: user.UserName,
			ReplyCount: replyCount,
			VoteNum:    voteNum,
			CreateTime: comment.CreateTime.Format("2006-01-02 15:04:05"),
		}
		data = append(data, commentDetail)
	}
	return data, nil
}

// CreateCommentReply 创建评论回复
func CreateCommentReply(userId int64, p *models.ParamCommentReply) error {
	// 生成评论ID
	commentId := snowflake.GetID()

	// 转换父评论ID和帖子ID
	parentId, err := strconv.ParseInt(p.ParentId, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt(p.ParentId) failed",
			zap.String("parent_id", p.ParentId),
			zap.Error(err))
		return mysql.ErrorInvalidID
	}
	postId, err := strconv.ParseInt(p.PostId, 10, 64)
	if err != nil {
		zap.L().Error("strconv.ParseInt(p.PostId) failed",
			zap.String("post_id", p.PostId),
			zap.Error(err))
		return mysql.ErrorInvalidID
	}

	// 检查帖子是否存在且状态正常
	post, err := mysql.GetPostById(postId)
	if err != nil {
		zap.L().Error("mysql.GetPostById failed",
			zap.Int64("post_id", postId),
			zap.Error(err))
		return err
	}
	if post == nil || post.Status != 1 {
		return mysql.ErrorInvalidID
	}

	// 检查父评论是否存在且属于该帖子
	parentComment, err := mysql.GetCommentById(parentId)
	if err != nil {
		zap.L().Error("mysql.GetCommentById failed",
			zap.Int64("parent_id", parentId),
			zap.Error(err))
		return err
	}
	if parentComment == nil || parentComment.Status != 1 {
		return mysql.ErrorInvalidID
	}
	// 检查父评论是否属于该帖子
	if parentComment.PostId != postId {
		zap.L().Error("parent comment not belong to the post",
			zap.Int64("parent_id", parentId),
			zap.Int64("post_id", postId),
			zap.Int64("parent_post_id", parentComment.PostId))
		return mysql.ErrorInvalidID
	}

	// 创建评论对象
	comment := &models.Comment{
		CommentId: commentId,
		ParentId:  parentId,
		PostId:    postId,
		AuthorId:  userId,
		Content:   p.Content,
		Status:    1, // 设置状态为有效
	}

	// 保存到MySQL
	if err := mysql.CreateComment(comment); err != nil {
		zap.L().Error("mysql.CreateComment failed",
			zap.Any("comment", comment),
			zap.Error(err))
		return err
	}

	// 保存到Redis
	if err := redis.CreateComment(commentId); err != nil {
		zap.L().Error("redis.CreateComment failed",
			zap.Int64("comment_id", commentId),
			zap.Error(err))
	}

	return nil
}

// GetCommentById 根据ID获取评论详情
func GetCommentById(commentId int64) (*models.ApiCommentDetail, error) {
	// 查询评论
	comment, err := mysql.GetCommentById(commentId)
	if err != nil {
		return nil, err
	}
	if comment == nil || comment.Status != 1 {
		return nil, mysql.ErrorInvalidID
	}

	// 查询评论作者信息
	user, err := mysql.GetUserById(comment.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserById failed",
			zap.Int64("author_id", comment.AuthorId),
			zap.Error(err))
		return nil, err
	}

	// 获取回复数量
	replyCount, err := mysql.GetCommentReplyCount(comment.CommentId)
	if err != nil {
		zap.L().Error("mysql.GetCommentReplyCount failed",
			zap.Int64("comment_id", comment.CommentId),
			zap.Error(err))
		replyCount = 0
	}

	// 获取点赞数量
	voteNum, err := redis.GetCommentVoteNum(strconv.FormatInt(comment.CommentId, 10))
	if err != nil {
		zap.L().Error("redis.GetCommentVoteNum failed",
			zap.Int64("comment_id", comment.CommentId),
			zap.Error(err))
		voteNum = 0
	}

	// 组装评论详情
	commentDetail := &models.ApiCommentDetail{
		CommentId:  comment.CommentId,
		ParentId:   comment.ParentId,
		PostId:     comment.PostId,
		AuthorId:   comment.AuthorId,
		Content:    comment.Content,
		AuthorName: user.UserName,
		ReplyCount: replyCount,
		VoteNum:    voteNum,
		CreateTime: comment.CreateTime.Format("2006-01-02 15:04:05"),
	}

	return commentDetail, nil
}
