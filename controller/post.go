package controller

import (
	"errors"
	"go_community/dao/mysql"
	"go_community/logic"
	"go_community/models"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler 创建帖子
func CreatePostHandler(c *gin.Context) {
	// 1.获取参数及校验参数
	//var post models.Post
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreatePostHandler with invalid param", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}
	// 从 c 取到当前发请求的用户 ID
	userID, err := getCurrentUserId(c)
	if err != nil {
		zap.L().Error("getCurrentUserID failed", zap.Error(err))
		ResponseError(c, CodeNotLogin)
		return
	}
	p.AuthorId = userID
	// 2.创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}

// PostDetailHandler 获取帖子详情
func PostDetailHandler(c *gin.Context) {
	// 1.获取参数（URL中获取帖子ID）
	postIdStr := c.Param("id")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		zap.L().Error("PostDetailHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
	}
	// 2.根据ID取出帖子数据（查数据库）
	post, err := logic.GetPostById(postId)
	if err != nil {
		zap.L().Error("logic.GetPostById failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 3.返回响应
	ResponseSuccess(c, post)
}

// GetPostListHandler 分页获取帖子列表
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	// 获取数据
	posts, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 返回响应
	ResponseSuccess(c, posts)
}

// GetPostListHandler2 分页获取帖子列表（按帖子的创建时间或者分数排序）
// @Summary 升级版帖子列表接口
// @Description 按社区按时间或分数排序查询帖子列表接口
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string false "Bearer 用户令牌"
// @Param object query models.ParamPostList false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} _ResponsePostList
// @Router /posts2 [get]
func GetPostListHandler2(c *gin.Context) {
	// GET 请求参数（query string）: /api/v1/post2?page=1&size=10&order=time
	p := &models.ParamPostList{
		Page:  1,                // 默认第1页
		Size:  10,               // 默认每页10条
		Order: models.OrderTime, // 默认按时间排序
	}
	// 获取分页参数
	// c.ShouldBind() 根据请求的数据类型，选择相应的方法去获取数据
	// c.ShouldBindJSON() 请求中携带 json 格式的数据，才能用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取数据
	posts, err := logic.GetPostListNew(p) // 更新：合二为一
	if err != nil {
		zap.L().Error("logic.GetPostListNew failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return // 添加return，避免错误情况下继续执行
	}

	// 返回响应
	ResponseSuccess(c, posts)
}

// PostSearchHandler 搜索帖子
func PostSearchHandler(c *gin.Context) {
	p := &models.ParamPostList{}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("PostSearchHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取数据
	data, err := logic.PostSearch(p)
	if err != nil {
		zap.L().Error("logic.PostSearch failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// UpdatePostHandler 更新帖子
func UpdatePostHandler(c *gin.Context) {
	// 1. 参数校验
	p := new(models.ParamUpdatePost)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdatePostHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 2. 获取当前用户ID
	userID, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	// 3. 更新帖子
	if err := logic.UpdatePost(userID, p); err != nil {
		zap.L().Error("logic.UpdatePost failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorNoPermission) {
			ResponseError(c, CodeNoPermission)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

// GetUserPostListHandler 根据用户ID获取帖子列表
func GetUserPostListHandler(c *gin.Context) {
	// 获取用户ID参数
	userIdStr := c.Param("id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取分页参数
	page, size := getPageInfo(c)

	// 获取数据
	data, err := logic.GetUserPostList(userId, page, size)
	if err != nil {
		zap.L().Error("logic.GetUserPostList failed",
			zap.Int64("user_id", userId),
			zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}

// DeletePostHandler 删除帖子
func DeletePostHandler(c *gin.Context) {
	// 1. 获取帖子ID
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 2. 获取当前登录用户ID
	userID, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	// 3. 删除帖子
	if err := logic.DeletePost(userID, postID); err != nil {
		zap.L().Error("logic.DeletePost failed",
			zap.Int64("post_id", postID),
			zap.Int64("user_id", userID),
			zap.Error(err))

		if err == mysql.ErrorInvalidID {
			ResponseErrorWithMsg(c, CodeInvalidParams, "帖子不存在或已删除")
			return
		}
		if err == mysql.ErrorNoPermission {
			ResponseErrorWithMsg(c, CodeNoPermission, "无权删除该帖子")
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}
