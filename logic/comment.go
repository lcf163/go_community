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
	comments, err := mysql.GetCommentListByPostId(postId)
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
