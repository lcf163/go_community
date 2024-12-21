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
func GetPostById(postId int64) (data *models.ApiPostDetail, err error) {
	// 查询并组合接口需要的数据
	// 查询帖子信息
	post, err := mysql.GetPostById(postId)
	if err != nil {
		zap.L().Error("mysql.GetPostById(postId) failed",
			zap.Int64("postId", postId),
			zap.Error(err))
		return
	}
	// 根据作者id查询作者信息
	user, err := mysql.GetUserById(post.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserById(AuthorId) failed",
			zap.Int64("AuthorId", post.AuthorId),
			zap.Error(err))
		return
	}
	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailById(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID(post.CommunityID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	// 接口数据拼接
	data = &models.ApiPostDetail{
		AuthorName:      user.UserName,
		Post:            post,
		CommunityDetail: community,
	}
	return
}
