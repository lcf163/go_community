package mysql

import (
	"go_community/models"
)

// CreateComment 创建评论
func CreateComment(comment *models.Comment) (err error) {
	sqlStr := `insert into comment(
		comment_id, parent_id, post_id, author_id, content)
		values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr,
		comment.CommentId,
		comment.ParentId,
		comment.PostId,
		comment.AuthorId,
		comment.Content)
	return
}

// GetCommentListByPostId 获取帖子的评论列表
func GetCommentListByPostId(postId int64) ([]*models.Comment, error) {
	sqlStr := `select id, comment_id, parent_id, post_id, author_id, content, create_time
        from comment where post_id = ? and status = 1 order by create_time desc`
	comments := make([]*models.Comment, 0)
	err := db.Select(&comments, sqlStr, postId)
	return comments, err
}

// UpdateComment 修改评论
func UpdateComment(commentId int64, content string) error {
	sqlStr := `update comment set content = ? where comment_id = ? and status = 1`
	_, err := db.Exec(sqlStr, content, commentId)
	return err
}

// GetCommentReplyCount 获取评论的回复数量
func GetCommentReplyCount(commentId int64) (int64, error) {
	sqlStr := `select count(*) from comment where parent_id = ? and status = 1`
	var count int64
	err := db.Get(&count, sqlStr, commentId)
	return count, err
}
