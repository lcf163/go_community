package mysql

import (
	"database/sql"
	"go_community/models"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"go.uber.org/zap"
)

// GetPostTotalCount 查询帖子总数
func GetPostTotalCount() (count int64, err error) {
	sqlStr := `select count(post_id) from post where status = 1`
	err = db.Get(&count, sqlStr)
	if err != nil {
		zap.L().Error("db.Get(&count, sqlStr) failed", zap.Error(err))
		return 0, err
	}
	return
}

// GetCommunityPostTotalCount 根据社区id查询数据库帖子总数
func GetCommunityPostTotalCount(communityId int64) (count int64, err error) {
	sqlStr := `select count(post_id) from post where community_id = ? and status = 1`
	err = db.Get(&count, sqlStr, communityId)
	if err != nil {
		zap.L().Error("db.Get(&count, sqlStr) failed", zap.Error(err))
		return 0, err
	}
	return
}

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
	sqlStr := `select post_id, title, content, author_id, community_id, status, create_time
	from post
	where post_id = ? and status = 1`

	// 打印SQL语句
	zap.L().Debug("GetPostById SQL",
		zap.String("sql", sqlStr),
		zap.Int64("post_id", postId))

	err = db.Get(post, sqlStr, postId)
	if err == sql.ErrNoRows {
		return nil, ErrorInvalidID
	}
	if err != nil {
		zap.L().Error("query post failed",
			zap.String("sql", sqlStr),
			zap.Int64("post_id", postId),
			zap.Error(err))
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
		zap.L().Error("query post list failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}

// GetPostListByIds 根据给定的 id 列表查询帖子数据
func GetPostListByIds(ids []string) (postList []*models.Post, err error) {
	sqlStr := `select post_id, title, content, author_id, community_id, create_time
	from post
	where post_id in (?)
	and status = 1
	order by FIND_IN_SET(post_id, ?)`
	// 动态填充 id
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return
	}
	// sqlx.In 返回带 `?` bindvar 的查询语句, 使用 Rebind() 重新绑定它
	query = db.Rebind(query)
	err = db.Select(&postList, query, args...)
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

	// 添加日志记录实际执行的 SQL
	zap.L().Debug("GetPostListByKeywords SQL",
		zap.String("sql", sqlStr),
		zap.String("keyword", p.Search),
		zap.Int64("offset", (p.Page-1)*p.Size),
		zap.Int64("size", p.Size))

	return
}

// UpdatePost 更新帖子
func UpdatePost(postId int64, title string, content string) error {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}

	sqlStr := `update post 
	set title = ?, 
	content = ?, 
	update_time = ? 
	where post_id = ?`
	_, err = db.Exec(sqlStr, title, content, time.Now().In(loc), postId)
	if err != nil {
		zap.L().Error("update post failed", zap.Error(err))
		return err
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
