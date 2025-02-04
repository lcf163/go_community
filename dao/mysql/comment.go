package mysql

import (
	"database/sql"
	"go_community/models"
	"strconv"

	"go.uber.org/zap"
)

// CreateComment 创建评论
func CreateComment(comment *models.Comment) (err error) {
	// 确保新创建的评论状态为1
	comment.Status = 1

	sqlStr := `insert into comment(
	comment_id, parent_id, post_id, author_id, reply_to_uid, content, status
	) values(?,?,?,?,?,?,?)`

	_, err = db.Exec(sqlStr,
		comment.CommentID,
		comment.ParentID,
		comment.PostID,
		comment.AuthorID,
		comment.ReplyToUID,
		comment.Content,
		comment.Status)

	if err != nil {
		zap.L().Error("CreateComment failed",
			zap.String("sql", sqlStr),
			zap.Any("comment", comment),
			zap.Error(err))
		return err
	}
	return nil
}

// UpdateComment 修改评论
func UpdateComment(commentId int64, content string) error {
	sqlStr := `update comment set content = ? where comment_id = ? and status = 1`
	result, err := db.Exec(sqlStr, content, commentId)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrorInvalidID
	}
	return nil
}

// GetCommentCount 获取帖子的评论数量
func GetCommentCount(postId int64) (count int64, err error) {
	sqlStr := `select count(comment_id) from comment where post_id = ? and parent_id = 0 and status = 1`
	err = db.Get(&count, sqlStr, postId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetCommentList 获取帖子的评论列表(分页)
func GetCommentList(postId int64, page, size int64) ([]*models.Comment, error) {
	sqlStr := `select comment_id, parent_id, post_id, author_id, content, create_time
	from comment 
	where post_id = ? and parent_id = 0 and status = 1 
	order by create_time desc
	limit ?, ?`

	comments := make([]*models.Comment, 0)
	err := db.Select(&comments, sqlStr, postId, (page-1)*size, size)
	return comments, err
}

// GetCommentReplyCount 获取评论的回复数量
func GetCommentReplyCount(commentId int64) (int64, error) {
	sqlStr := `select count(*) from comment where parent_id = ? and status = 1`
	var count int64
	err := db.Get(&count, sqlStr, commentId)
	return count, err
}

// GetCommentReplyList 获取评论的回复列表
func GetCommentReplyList(commentId int64) ([]*models.Comment, error) {
	sqlStr := `select comment_id, parent_id, post_id, author_id, reply_to_uid, content, status, create_time
	from comment 
	where parent_id = ? and status = 1 
	order by create_time desc`

	comments := make([]*models.Comment, 0)
	err := db.Select(&comments, sqlStr, commentId)
	return comments, err
}

// GetCommentById 根据评论ID获取评论
func GetCommentById(commentId int64) (comment *models.Comment, err error) {
	comment = new(models.Comment)
	sqlStr := `select comment_id, content, post_id, author_id, parent_id, create_time 
	from comment 
	where comment_id = ? and status = 1`
	err = db.Get(comment, sqlStr, commentId)
	if err == sql.ErrNoRows {
		return nil, ErrorInvalidID
	}
	return
}

// DeleteCommentWithTx 删除评论(使用事务)
func DeleteCommentWithTx(tx *sql.Tx, commentID int64) error {
	sqlStr := `update comment set status = 0 where comment_id = ? and status = 1`
	result, err := tx.Exec(sqlStr, commentID)
	if err != nil {
		return err
	}
	// 检查是否更新了记录
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrorInvalidID
	}

	return nil
}

// GetCommentRepliesIDs 获取评论的所有回复ID
func GetCommentRepliesIDs(tx *sql.Tx, commentID int64) ([]string, error) {
	sqlStr := `select comment_id from comment where parent_id = ? and status = 1`
	rows, err := tx.Query(sqlStr, commentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replyIDs []string
	for rows.Next() {
		var replyID int64
		if err := rows.Scan(&replyID); err != nil {
			return nil, err
		}
		replyIDs = append(replyIDs, strconv.FormatInt(replyID, 10))
	}
	return replyIDs, nil
}

// DeleteCommentRepliesWithTx 删除评论的所有回复(使用事务)
func DeleteCommentRepliesWithTx(tx *sql.Tx, commentID int64) error {
	sqlStr := `update comment set status = 0 where parent_id = ? and status = 1`
	_, err := tx.Exec(sqlStr, commentID)
	return err
}
