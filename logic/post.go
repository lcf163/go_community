package logic

import (
	"go-community/dao/mysql"
	"go-community/models"
	"go-community/pkg/snowflake"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(post *models.Post) (err error) {
	// 生成帖子ID
	postID := snowflake.GetID()
	post.PostID = postID
	// 创建帖子，保存到数据库
	if err := mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost failed", zap.Error(err))
		return err
	}
	return
}

// GetPostById 根据帖子ID查询帖子详情
func GetPostById(postId int64) (post *models.Post, err error) {
	// 查询帖子信息
	return mysql.GetPostById(postId)
}
