package mysql

import (
	"database/sql"
	"go-community/models"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(post *models.Post) (err error) {
	sqlStr := `insert into post(
	post_id, title, content, author_id, community_id)
	values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, post.PostId, post.Title,
		post.Content, post.AuthorId, post.CommunityId)
	if err != nil {
		zap.L().Error("insert post failed", zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

// GetPostById 根据ID获取帖子
func GetPostById(postId int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id = ?`
	err = db.Get(post, sqlStr, postId)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		zap.L().Error("query post failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}

func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	limit ?,?`
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	if err != nil {
		zap.L().Error("query post list failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}
