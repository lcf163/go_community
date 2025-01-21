package logic

import (
	"go_community/dao/mysql"
	"go_community/dao/redis"
	"go_community/models"
	"go_community/pkg/snowflake"
	"strconv"

	"go.uber.org/zap"
)

// CreateComment 创建评论/回复
func CreateComment(userId int64, p *models.ParamComment) error {
	// 检查帖子是否存在
	post, err := mysql.GetPostById(p.PostId)
	if err != nil || post == nil || post.Status != 1 {
		zap.L().Error("mysql.GetPostById failed",
			zap.Int64("post_id", p.PostId),
			zap.Error(err))
		return mysql.ErrorInvalidID
	}

	// 如果是回复,需要额外检查
	if p.ParentId != 0 {
		// 检查父评论是否存在
		parentComment, err := mysql.GetCommentById(p.ParentId)
		if err != nil || parentComment == nil || parentComment.Status != 1 {
			zap.L().Error("mysql.GetCommentById failed",
				zap.Int64("parent_id", p.ParentId),
				zap.Error(err))
			return mysql.ErrorInvalidID
		}
		// 检查父评论是否属于指定的帖子
		if parentComment.PostId != p.PostId {
			zap.L().Error("parent comment does not belong to the specified post",
				zap.Int64("comment_id", p.ParentId),
				zap.Int64("post_id", p.PostId))
			return mysql.ErrorInvalidID
		}
	}

	// 生成评论ID
	commentId := snowflake.GetID()

	// 创建评论
	comment := &models.Comment{
		CommentId:  commentId,
		ParentId:   p.ParentId,
		PostId:     p.PostId,
		AuthorId:   userId,
		ReplyToUid: p.ReplyToUid,
		Content:    p.Content,
		Status:     1,
	}

	// 保存到数据库
	if err := mysql.CreateComment(comment); err != nil {
		return err
	}

	// 保存到Redis
	return redis.CreateComment(commentId)
}

// GetCommentList 获取评论列表
func GetCommentList(postId int64, page, size int64) (*models.ApiCommentListRes, error) {
	// 获取评论总数
	total, err := mysql.GetCommentCount(postId)
	if err != nil {
		return nil, err
	}

	// 获取分页数据
	comments, err := mysql.GetCommentList(postId, page, size)
	if err != nil {
		return nil, err
	}

	// 组装评论详情
	data := make([]*models.ApiCommentDetail, 0, len(comments))
	for _, comment := range comments {
		// 获取评论作者信息
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

		// 获取评论点赞数
		voteNum, err := redis.GetCommentVoteNum(strconv.FormatInt(comment.CommentId, 10))
		if err != nil {
			zap.L().Error("redis.GetCommentVoteNum failed",
				zap.Int64("comment_id", comment.CommentId),
				zap.Error(err))
			voteNum = 0
		}

		commentDetail := &models.ApiCommentDetail{
			CommentId:    comment.CommentId,
			ParentId:     comment.ParentId,
			PostId:       comment.PostId,
			AuthorId:     comment.AuthorId,
			Content:      comment.Content,
			AuthorName:   user.UserName,
			AuthorAvatar: user.Avatar,
			ReplyCount:   replyCount,
			VoteNum:      voteNum,
			CreateTime:   comment.CreateTime.Format("2006-01-02 15:04:05"),
		}
		data = append(data, commentDetail)
	}

	// 组装返回数据
	return &models.ApiCommentListRes{
		Page: &models.Page{
			Page:  page,
			Size:  size,
			Total: total,
		},
		List: data,
	}, nil
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
		author, err := mysql.GetUserById(comment.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById failed",
				zap.Int64("author_id", comment.AuthorId),
				zap.Error(err))
			continue
		}

		// 查询被回复人信息（仅当 ReplyToUid 不为 0 时）
		var replyToUser *models.User
		if comment.ReplyToUid != 0 {
			replyToUser, err = mysql.GetUserById(comment.ReplyToUid)
			if err != nil {
				zap.L().Error("mysql.GetUserById failed",
					zap.Int64("reply_to_uid", comment.ReplyToUid),
					zap.Error(err))
				continue
			}
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
			CommentId:    comment.CommentId,
			ParentId:     comment.ParentId,
			PostId:       comment.PostId,
			AuthorId:     comment.AuthorId,
			Content:      comment.Content,
			AuthorName:   author.UserName,
			AuthorAvatar: author.GetAvatarURL(),
			ReplyCount:   replyCount,
			VoteNum:      voteNum,
			CreateTime:   comment.CreateTime.Format("2006-01-02 15:04:05"),
		}

		// 只有在有被回复用户时才设置被回复人信息
		if replyToUser != nil {
			commentDetail.ReplyToUid = comment.ReplyToUid
			commentDetail.ReplyToUserName = replyToUser.UserName
			commentDetail.ReplyToUserAvatar = replyToUser.GetAvatarURL()
		}

		data = append(data, commentDetail)
	}

	return data, nil
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

// UpdateComment 更新评论
func UpdateComment(userId int64, p *models.ParamUpdateComment) error {
	// 检查评论是否存在
	comment, err := mysql.GetCommentById(p.CommentId)
	if err != nil || comment == nil || comment.Status != 1 {
		zap.L().Error("mysql.GetCommentById failed",
			zap.Int64("comment_id", p.CommentId),
			zap.Error(err))
		return mysql.ErrorInvalidID
	}

	// 检查是否是评论作者
	if comment.AuthorId != userId {
		zap.L().Error("no permission to update comment",
			zap.Int64("comment_id", p.CommentId),
			zap.Int64("user_id", userId),
			zap.Int64("author_id", comment.AuthorId))
		return mysql.ErrorNoPermission
	}

	// 更新评论内容
	return mysql.UpdateComment(p.CommentId, p.Content)
}

// DeleteComment 删除评论
func DeleteComment(userID, commentID int64) error {
	// 1. 检查评论是否存在
	comment, err := mysql.GetCommentById(commentID)
	if err != nil {
		return err
	}
	if comment == nil || comment.Status == 0 {
		return mysql.ErrorInvalidID
	}

	// 2. 检查是否是评论作者
	if comment.AuthorId != userID {
		zap.L().Error("no permission to delete comment",
			zap.Int64("comment_id", commentID),
			zap.Int64("user_id", userID),
			zap.Int64("author_id", comment.AuthorId))
		return mysql.ErrorNoPermission
	}

	// 3. 开启事务
	tx, err := mysql.GetDB().Begin()
	if err != nil {
		return err
	}

	// 使用 defer 处理事务的提交和回滚
	defer func() {
		if p := recover(); p != nil {
			// 发生 panic 时回滚事务
			tx.Rollback()
			panic(p) // 重新抛出 panic
		} else if err != nil {
			// 有错误时回滚事务
			tx.Rollback()
		} else {
			// 无错误时提交事务
			err = tx.Commit()
		}
	}()

	// 4. 软删除评论（将状态设为0）
	if err = mysql.DeleteCommentWithTx(tx, commentID); err != nil {
		return err // defer 中会处理回滚
	}

	// 5. 删除Redis中的评论数据
	if err = redis.DeleteCommentVote(strconv.FormatInt(commentID, 10)); err != nil {
		return err // defer 中会处理回滚
	}

	return nil // defer 中会处理提交
}
