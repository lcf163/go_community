package logic

import (
	"go_community/dao/mysql"
	"go_community/dao/redis"
	"go_community/models"
	"go_community/pkg/snowflake"
	"strconv"

	"go.uber.org/zap"
)

// CreatePost 创建帖子
func CreatePost(p *models.Post) (err error) {
	// 生成帖子ID
	p.PostID = snowflake.GetID()

	// 创建帖子，保存到数据库
	if err := mysql.CreatePost(p); err != nil {
		zap.L().Error("mysql.CreatePost failed", zap.Error(err))
		return err
	}
	if err := redis.CreatePost(p.PostID, p.CommunityID); err != nil {
		zap.L().Error("redis.CreatePost failed", zap.Error(err))
		return err
	}
	return nil
}

// GetPostById 根据帖子ID查询帖子详情
func GetPostById(postID int64) (data *models.ApiPostDetail, err error) {
	// 查询帖子信息
	post, err := mysql.GetPostById(postID)
	if err != nil {
		zap.L().Error("mysql.GetPostById(postID) failed",
			zap.Int64("post_communityID", postID),
			zap.Error(err))
		return nil, err
	}

	// 查询作者信息
	user, err := mysql.GetUserById(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return nil, err
	}

	// 查询社区信息
	community, err := mysql.GetCommunityDetailById(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailById(post.CommunityID) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return nil, err
	}

	// 获取帖子投票数
	voteNum, err := redis.GetPostVoteNum(strconv.FormatInt(postID, 10))
	if err != nil {
		zap.L().Error("redis.GetPostVoteNum failed",
			zap.Int64("post_id", postID),
			zap.Error(err))
		voteNum = 0
	}

	// 获取评论数量
	commentCount, err := mysql.GetCommentCount(postID)
	if err != nil {
		zap.L().Error("mysql.GetCommentCount(postID) failed",
			zap.Int64("post_id", postID),
			zap.Error(err))
		commentCount = 0
	}

	// 组装数据
	data = &models.ApiPostDetail{
		AuthorName:      user.UserName,
		AuthorAvatar:    user.GetAvatarURL(),
		VoteNum:         voteNum,
		CommentCount:    commentCount,
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
	// 初始化返回数据结构
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
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
func GetPostList2(p *models.ParamPostList) (data *models.ApiPostDetailRes, err error) {
	// 初始化返回数据结构
	data = &models.ApiPostDetailRes{
		Page: models.Page{},
		List: make([]*models.ApiPostDetail, 0),
	}

	// 从mysql获取帖子列表总数
	total, err := mysql.GetPostTotalCount()
	if err != nil {
		return nil, err
	}
	data.Page.Total = total
	data.Page.Page = p.Page
	data.Page.Size = p.Size

	// redis 中查询 Id 列表
	ids, err := redis.GetPostIdsInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIdsInOrder(p), return data is empty")
		return data, nil
	}

	// 根据 Id 在数据库 mysql 中查询帖子详细信息
	// 返回的数据需要按照给定的 id 的顺序，order by FIND_IN_SET(post_id, ?)
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		return nil, err
	}

	// 添加数据一致性检查
	if len(posts) != len(ids) {
		zap.L().Warn("data inconsistency between Redis and MySQL",
			zap.Int("redis_count", len(ids)),
			zap.Int("mysql_count", len(posts)),
			zap.Strings("redis_ids", ids))

		// 找出在 MySQL 中不存在的 ID
		existingIds := make(map[string]bool)
		for _, post := range posts {
			existingIds[strconv.FormatInt(post.PostID, 10)] = true
		}

		var invalidIds []string
		for _, id := range ids {
			if !existingIds[id] {
				invalidIds = append(invalidIds, id)
			}
		}

		// 异步清理 Redis 中的无效数据
		if len(invalidIds) > 0 {
			go func() {
				if err := redis.RemoveInvalidPostIds(invalidIds); err != nil {
					zap.L().Error("failed to remove invalid post ids from redis",
						zap.Error(err),
						zap.Strings("invalid_ids", invalidIds))
				} else {
					zap.L().Info("successfully removed invalid post ids from redis",
						zap.Strings("invalid_ids", invalidIds))
				}
			}()
		}
	}

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	// 组合数据
	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityByID(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		// 获取评论数量
		commentCount, err := mysql.GetCommentCount(post.PostID)
		if err != nil {
			zap.L().Error("mysql.GetCommentCount(post.PostID) failed",
				zap.Int64("post_id", post.PostID),
				zap.Error(err))
			commentCount = 0
		}

		// 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			AuthorAvatar:    user.GetAvatarURL(),
			VoteNum:         voteData[idx],
			CommentCount:    commentCount,
			Post:            post,
			CommunityDetail: community,
		}
		data.List = append(data.List, postDetail)
	}
	return data, nil
}

// GetCommunityPostList 根据社区id去查询帖子列表
func GetCommunityPostList(p *models.ParamPostList) (data *models.ApiPostDetailRes, err error) {
	// 初始化返回数据结构
	data = &models.ApiPostDetailRes{
		Page: models.Page{},
		List: make([]*models.ApiPostDetail, 0),
	}

	// 从mysql获取该社区下帖子列表总数
	total, err := mysql.GetCommunityPostTotalCount(p.CommunityID)
	if err != nil {
		return
	}
	data.Page.Total = total
	data.Page.Page = p.Page
	data.Page.Size = p.Size

	// redis 中查询 Id 列表
	ids, err := redis.GetCommunityPostIdsInOrder(p)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		zap.L().Warn("redis.GetCommunityPostIdsInOrder(p), return data is empty")
		return
	}
	zap.L().Debug("GetCommunityPostList", zap.Any("ids: ", ids))

	// 根据 Id 在数据库 mysql 中查询帖子详细信息
	// 返回的数据需要按照给定的 id 的顺序，order by FIND_IN_SET(post_id, ?)
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		return nil, err
	}
	zap.L().Debug("GetCommunityPostList", zap.Any("posts: ", posts))

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 根据社区id查询社区详细信息
	// 为了减少数据库的查询次数，这里将社区信息提前查询出来
	community, err := mysql.GetCommunityDetailById(p.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailById(p.CommunityID) failed",
			zap.Int64("community_id", p.CommunityID),
			zap.Error(err))
		return
	}

	// 组合数据
	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 过滤掉不属于该社区的帖子
		if post.CommunityID != p.CommunityID {
			continue
		}
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		// 获取评论数量
		commentCount, err := mysql.GetCommentCount(post.PostID)
		if err != nil {
			zap.L().Error("mysql.GetCommentCount(post.PostID) failed",
				zap.Int64("post_id", post.PostID),
				zap.Error(err))
			commentCount = 0
		}

		// 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			AuthorAvatar:    user.GetAvatarURL(),
			VoteNum:         voteData[idx],
			CommentCount:    commentCount,
			Post:            post,
			CommunityDetail: community,
		}
		data.List = append(data.List, postDetail)
	}
	return
}

// GetPostListNew 将两个查询帖子列表的逻辑合二为一
func GetPostListNew(p *models.ParamPostList) (data *models.ApiPostDetailRes, err error) {
	// 根据请求参数的不同，执行不同的业务逻辑
	if p.CommunityID == 0 {
		// 查询所有帖子
		data, err = GetPostList2(p)
	} else {
		// 根据社区id查询
		data, err = GetCommunityPostList(p)
	}

	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}

	return data, nil
}

// PostSearch 搜索帖子
func PostSearch(p *models.ParamPostList) (data *models.ApiPostDetailRes, err error) {
	// 初始化返回数据结构
	data = &models.ApiPostDetailRes{
		Page: models.Page{},
		List: make([]*models.ApiPostDetail, 0),
	}

	// 根据搜索条件去mysql查询符合条件的帖子列表总数
	total, err := mysql.GetPostListTotalCount(p)
	if err != nil {
		return nil, err
	}
	data.Page.Total = total
	data.Page.Size = p.Size
	data.Page.Page = p.Page

	// 根据搜索条件去mysql分页查询符合条件的帖子列表
	posts, err := mysql.GetPostListByKeywords(p)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return data, nil
	}
	// 查询出来的帖子id列表传入到redis接口获取帖子的投票数
	ids := make([]string, 0, len(posts))
	for _, post := range posts {
		ids = append(ids, strconv.Itoa(int(post.PostID)))
	}
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	// 组合数据
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserById(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailById(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		// 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			AuthorAvatar:    user.GetAvatarURL(),
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data.List = append(data.List, postDetail)
	}
	return data, nil
}

// UpdatePost 编辑帖子
func UpdatePost(userID int64, p *models.ParamUpdatePost) (err error) {
	// 判断帖子是否存在
	post, err := mysql.GetPostById(p.PostID)
	if err != nil {
		return err
	}

	// 判断是否是帖子作者
	if post.AuthorID != userID {
		return mysql.ErrorNoPermission
	}

	// 更新帖子
	return mysql.UpdatePost(p.PostID, p.Title, p.Content)
}

// GetUserPostList 获取用户的帖子列表
func GetUserPostList(userID, page, size int64) (data *models.ApiPostDetailRes, err error) {
	// 初始化返回数据结构
	data = &models.ApiPostDetailRes{
		Page: models.Page{},
		List: make([]*models.ApiPostDetail, 0),
	}

	// 获取该用户的帖子总数
	total, err := mysql.GetUserPostTotalCount(userID)
	if err != nil {
		return nil, err
	}
	data.Page.Total = total
	data.Page.Size = size
	data.Page.Page = page

	// 查询该用户的帖子列表
	posts, err := mysql.GetUserPostList(userID, page, size)
	if err != nil {
		return nil, err
	}
	if len(posts) == 0 {
		return data, nil
	}

	// 提前查询好每篇帖子的投票数
	ids := make([]string, 0, len(posts))
	for _, post := range posts {
		ids = append(ids, strconv.FormatInt(post.PostID, 10))
	}
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	// 获取用户信息（只需要查询一次）
	user, err := mysql.GetUserById(userID)
	if err != nil {
		zap.L().Error("mysql.GetUserById(userID) failed",
			zap.Int64("user_id", userID),
			zap.Error(err))
		return nil, err
	}

	// 组装数据
	for idx, post := range posts {
		// 获取社区信息
		community, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailById(post.CommunityID) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}

		// 获取评论数量
		commentCount, err := mysql.GetCommentCount(post.PostID)
		if err != nil {
			zap.L().Error("mysql.GetCommentCount(post.PostID) failed",
				zap.Int64("post_id", post.PostID),
				zap.Error(err))
			commentCount = 0
		}

		// 组装数据
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			AuthorAvatar:    user.GetAvatarURL(),
			VoteNum:         voteData[idx],
			CommentCount:    commentCount,
			Post:            post,
			CommunityDetail: community,
		}
		data.List = append(data.List, postDetail)
	}

	return data, nil
}

// DeletePost 删除帖子
func DeletePost(userID, postID int64) error {
	// 1. 检查帖子是否存在
	post, err := mysql.GetPostById(postID)
	if err != nil {
		return err
	}
	if post == nil {
		return mysql.ErrorInvalidID
	}

	// 2. 检查是否有权限删除（是否是帖子作者）
	if post.AuthorID != userID {
		zap.L().Error("no permission to delete post",
			zap.Int64("post_id", postID),
			zap.Int64("user_id", userID),
			zap.Int64("author_id", post.AuthorID))
		return mysql.ErrorNoPermission
	}

	// 3. 开启事务
	tx, err := mysql.GetDB().Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// 4. 获取需要删除的评论ID列表
	commentIDs, err := mysql.GetPostCommentIDs(tx, postID)
	if err != nil {
		return err
	}

	// 5. 删除帖子(软删除)
	if err = mysql.DeletePostWithTx(tx, postID); err != nil {
		return err
	}

	// 6. 删除帖子下的所有评论(软删除)
	if err = mysql.DeletePostCommentsWithTx(tx, postID); err != nil {
		return err
	}

	// 7. 删除Redis中的相关数据(包括帖子和评论的数据)
	if err = redis.DeletePostData(strconv.FormatInt(postID, 10), commentIDs); err != nil {
		return err
	}

	return nil
}
