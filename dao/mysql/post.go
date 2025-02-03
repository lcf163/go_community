package mysql

import (
	"database/sql"
	"go_community/models"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

// GetPostTotalCount 查询帖子总数
func GetPostTotalCount() (count int64, err error) {
	sqlStr := `select count(post_id) from post where status = 1`
	err = db.Get(&count, sqlStr)
	if err != nil {
		return 0, err
	}
	return
}

// GetCommunityPostTotalCount 根据社区id查询数据库帖子总数
func GetCommunityPostTotalCount(communityId int64) (count int64, err error) {
	sqlStr := `select count(post_id) from post where community_id = ? and status = 1`
	err = db.Get(&count, sqlStr, communityId)
	if err != nil {
		return 0, err
	}
	return
}

// CreatePost 创建帖子
func CreatePost(post *models.Post) (err error) {
	// 设置默认状态为1
	post.Status = 1
	sqlStr := `insert into post(
	post_id, title, content, author_id, community_id, status)
	values(?,?,?,?,?,?)`
	_, err = db.Exec(sqlStr, post.PostId, post.Title,
		post.Content, post.AuthorId, post.CommunityId, post.Status)
	if err != nil {
		zap.L().Error("CreatePost failed",
			zap.String("sql", sqlStr),
			zap.Any("post", post),
			zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

// GetPostById 根据ID获取帖子
func GetPostById(postId int64) (post *models.Post, err error) {
	post = new(models.Post)
	sqlStr := `select post_id, title, content, author_id, community_id, create_time, update_time, status
	from post
	where post_id = ? and status = 1`

	err = db.Get(post, sqlStr, postId)
	if err == sql.ErrNoRows {
		return nil, ErrorInvalidID
	}
	if err != nil {
		return nil, err
	}
	return post, nil
}

// GetPostList 查询帖子列表
func GetPostList(page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where status = 1
	ORDER BY create_time 
	DESC 
	limit ?,?`
	posts = make([]*models.Post, 0, 2)
	err = db.Select(&posts, sqlStr, (page-1)*size, size)
	if err != nil {
		zap.L().Error("GetPostList failed",
			zap.String("sql", sqlStr),
			zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}

// GetPostListByIds 根据给定的id列表查询帖子数据
func GetPostListByIds(ids []string) (posts []*models.Post, err error) {
	// 使用 FIND_IN_SET 来确保返回的数据按照给定的id顺序
	// 使用 ORDER BY create_time DESC 来确保时间降序
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?) and status = 1
	ORDER BY create_time DESC`

	// 使用 sqlx.In 来动态生成 IN 查询语句
	query, args, err := sqlx.In(sqlStr, ids)
	if err != nil {
		return
	}

	// sqlx.In 返回带 `?` 的查询语句, 我们使用 Rebind() 重新绑定它
	query = db.Rebind(query)
	// 执行查询
	err = db.Select(&posts, query, args...)
	return
}

// GetPostListTotalCount 根据关键词查询帖子列表总数
func GetPostListTotalCount(p *models.ParamPostList) (count int64, err error) {
	// 根据帖子标题或者帖子内容模糊查询帖子列表总数
	sqlStr := `select count(post_id)
	from post
	where status = 1
	and (
		title like ?
		or content like ?
	)`
	keyword := "%" + p.Search + "%"
	err = db.Get(&count, sqlStr, keyword, keyword)
	if err != nil {
		return
	}
	return
}

// GetPostListByKeywords 根据关键词查询帖子列表
func GetPostListByKeywords(p *models.ParamPostList) (posts []*models.Post, err error) {
	// 根据帖子标题或者帖子内容模糊查询帖子列表
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where status = 1
	and (
		title like ?
		or content like ?
	)
	ORDER BY create_time
	DESC
	limit ?,?
	`
	// %keyword%
	p.Search = "%" + p.Search + "%"
	// 初始化 posts 切片
	posts = make([]*models.Post, 0, 2)
	// 执行查询
	err = db.Select(&posts, sqlStr, p.Search, p.Search, (p.Page-1)*p.Size, p.Size)
	if err != nil {
		// 添加日志记录实际执行的 SQL
		zap.L().Error("GetPostListByKeywords failed",
			zap.String("sql", sqlStr),
			zap.Any("param", p),
			zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}

// UpdatePost 更新帖子
func UpdatePost(postId int64, title string, content string) error {
	sqlStr := `update post 
	set title = ?, content = ?, update_time = ? 
	where post_id = ? and status = 1`
	result, err := db.Exec(sqlStr, title, content, time.Now(), postId)
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

// GetUserPostTotalCount 获取用户发帖总数
func GetUserPostTotalCount(userId int64) (count int64, err error) {
	sqlStr := `select count(post_id) from post where author_id = ? and status = 1`
	err = db.Get(&count, sqlStr, userId)
	return
}

// GetUserPostList 获取用户的帖子列表
func GetUserPostList(userId, page, size int64) (posts []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time 
	from post 
	where author_id = ? and status = 1
	order by create_time desc 
	limit ?,?`
	posts = make([]*models.Post, 0, size)
	err = db.Select(&posts, sqlStr, userId, (page-1)*size, size)
	return
}

// DeletePostWithTx 删除帖子(使用事务)
func DeletePostWithTx(tx *sql.Tx, postID int64) error {
	sqlStr := `update post set status = 0 where post_id = ? and status = 1`
	result, err := tx.Exec(sqlStr, postID)
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

// DeletePostCommentsWithTx 删除帖子下的所有评论(使用事务)
func DeletePostCommentsWithTx(tx *sql.Tx, postID int64) error {
	sqlStr := `update comment set status = 0 where post_id = ? and status = 1`
	_, err := tx.Exec(sqlStr, postID)
	return err
}

// GetPostCommentIDs 获取帖子下所有评论的ID
func GetPostCommentIDs(tx *sql.Tx, postID int64) ([]string, error) {
	sqlStr := `select comment_id from comment where post_id = ? and status = 1`
	rows, err := tx.Query(sqlStr, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commentIDs []string
	for rows.Next() {
		var commentID int64
		if err := rows.Scan(&commentID); err != nil {
			return nil, err
		}
		commentIDs = append(commentIDs, strconv.FormatInt(commentID, 10))
	}
	return commentIDs, nil
}
