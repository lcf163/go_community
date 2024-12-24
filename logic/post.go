package logic

import (
	"fmt"
	"go-community/dao/mysql"
	"go-community/dao/redis"
	"go-community/models"
	"go-community/pkg/snowflake"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(post *models.Post) (err error) {
	// 生成帖子ID
	postId := snowflake.GetID()
	post.PostId = postId
	// 创建帖子，保存到数据库
	if err := mysql.CreatePost(post); err != nil {
		zap.L().Error("mysql.CreatePost failed", zap.Error(err))
		return err
	}
	if err := redis.CreatePost(fmt.Sprint(postId)); err != nil {
		zap.L().Error("redis.CreatePost failed", zap.Error(err))
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
		zap.L().Error("mysql.GetUserById(post.AuthorId) failed",
			zap.Int64("AuthorId", post.AuthorId),
			zap.Error(err))
		return
	}
	// 根据社区id查询社区详细信息
	community, err := mysql.GetCommunityDetailById(post.CommunityId)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID(post.CommunityId) failed",
			zap.Int64("community_id", post.CommunityId),
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

// GetPostList 获取帖子列表
func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	// 查询并组合接口需要的数据
	// 查询帖子信息
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		zap.L().Error("mysql.GetPostList failed", zap.Error(err))
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(AuthorId) failed",
				zap.Int64("AuthorId", post.AuthorId),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID(post.CommunityId) failed",
				zap.Int64("community_id", post.CommunityId),
				zap.Error(err))
			continue
		}
		// 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}

// GetPostList2 获取帖子列表（按帖子的创建时间或者分数排序）
func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// redis 中查询 Id 列表
	ids, err := redis.GetPostIdInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIdInOrder(p), return data is empty")
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("ids: ", ids))
	// 根据 Id 在数据库 mysql 中查询帖子详细信息
	// 返回的数据需要按照给定的 id 的顺序，order by FIND_IN_SET(post_id, ?)
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		return
	}
	zap.L().Debug("GetPostList2", zap.Any("posts: ", posts))

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	// 组合数据
	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(AuthorId) failed",
				zap.Int64("AuthorId", post.AuthorId),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID(post.CommunityId) failed",
				zap.Int64("community_id", post.CommunityId),
				zap.Error(err))
			continue
		}
		// 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postDetail)
	}
	return
}
