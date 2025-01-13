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
	p.PostId = snowflake.GetID()

	// 创建帖子，保存到数据库
	if err := mysql.CreatePost(p); err != nil {
		zap.L().Error("mysql.CreatePost failed", zap.Error(err))
		return err
	}
	if err := redis.CreatePost(p.PostId, p.CommunityId); err != nil {
		zap.L().Error("redis.CreatePost failed", zap.Error(err))
		return err
	}
	return nil
}

// GetPostById 根据帖子ID查询帖子详情
func GetPostById(postId int64) (data *models.ApiPostDetail, err error) {
	// 查询帖子信息
	post, err := mysql.GetPostById(postId)
	if err != nil {
		zap.L().Error("mysql.GetPostById(postId) failed",
			zap.Int64("postId", postId),
			zap.Error(err))
		return nil, err
	}

	// 查询作者信息
	user, err := mysql.GetUserById(post.AuthorId)
	if err != nil {
		zap.L().Error("mysql.GetUserById(post.AuthorId) failed",
			zap.Int64("AuthorId", post.AuthorId),
			zap.Error(err))
		return nil, err
	}

	// 查询每篇帖子的投票数
	voteData, err := redis.GetPostVoteNum(strconv.FormatInt(postId, 10))
	if err != nil {
		return nil, err
	}

	// 查询社区详细信息
	community, err := mysql.GetCommunityDetailById(post.CommunityId)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID(post.CommunityId) failed",
			zap.Int64("community_id", post.CommunityId),
			zap.Error(err))
		return nil, err
	}

	// 获取帖子的评论总数
	commentCount, err := mysql.GetCommentCount(postId)
	if err != nil {
		zap.L().Error("mysql.GetCommentCount failed",
			zap.Int64("post_id", postId),
			zap.Error(err))
		commentCount = 0 // 如果查询失败，默认为0
	}

	// 组装数据
	data = &models.ApiPostDetail{
		AuthorName:      user.UserName,
		VoteNum:         voteData,
		CommentCount:    commentCount,
		Post:            post,
		CommunityDetail: community,
	}
	return data, nil
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
		zap.L().Warn("redis.GetPostIdInOrder(p), return data is empty")
		return data, nil
	}
	zap.L().Debug("GetPostList2", zap.Any("ids: ", ids))

	// 根据 Id 在数据库 mysql 中查询帖子详细信息
	// 返回的数据需要按照给定的 id 的顺序，order by FIND_IN_SET(post_id, ?)
	posts, err := mysql.GetPostListByIds(ids)
	if err != nil {
		return nil, err
	}
	zap.L().Debug("GetPostList2", zap.Any("posts: ", posts))

	// 提前查询好每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
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

		// 获取评论数量
		commentCount, err := mysql.GetCommentCount(post.PostId)
		if err != nil {
			zap.L().Error("mysql.GetCommentCount failed",
				zap.Int64("post_id", post.PostId),
				zap.Error(err))
			commentCount = 0
		}

		// 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
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
	total, err := mysql.GetCommunityPostTotalCount(p.CommunityId)
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
	community, err := mysql.GetCommunityDetailById(p.CommunityId)
	if err != nil {
		zap.L().Error("mysql.GetCommunityByID(post.CommunityId) failed",
			zap.Int64("community_id", p.CommunityId),
			zap.Error(err))
		return
	}

	// 组合数据
	// 将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {
		// 过滤掉不属于该社区的帖子
		if post.CommunityId != p.CommunityId {
			continue
		}
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById(AuthorId) failed",
				zap.Int64("AuthorId", post.AuthorId),
				zap.Error(err))
			continue
		}

		// 获取评论数量
		commentCount, err := mysql.GetCommentCount(post.PostId)
		if err != nil {
			zap.L().Error("mysql.GetCommentCount failed",
				zap.Int64("post_id", post.PostId),
				zap.Error(err))
			commentCount = 0
		}

		// 接口数据拼接
		postDetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
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
	if p.CommunityId == 0 {
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
		ids = append(ids, strconv.Itoa(int(post.PostId)))
	}
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return nil, err
	}

	// 组合数据
	for idx, post := range posts {
		// 根据作者id查询作者信息
		user, err := mysql.GetUserById(post.AuthorId)
		if err != nil {
			zap.L().Error("mysql.GetUserById() failed",
				zap.Int64("postId", post.AuthorId),
				zap.Error(err))
			continue
		}

		// 根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailById(post.CommunityId)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailById() failed",
				zap.Int64("communityId", post.CommunityId),
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
		data.List = append(data.List, postDetail)
	}
	return data, nil
}

// UpdatePost 编辑帖子
func UpdatePost(userId int64, p *models.ParamUpdatePost) (err error) {
	// 把帖子id转换为int64
	postId, err := strconv.ParseInt(p.PostId, 10, 64)
	if err != nil {
		return err
	}

	// 判断帖子是否存在
	post, err := mysql.GetPostById(postId)
	if err != nil {
		return err
	}

	// 判断是否是帖子作者
	if post.AuthorId != userId {
		return mysql.ErrorNoPermission // 使用自定义错误
	}

	// 更新帖子
	return mysql.UpdatePost(postId, p.Title, p.Content)
}
