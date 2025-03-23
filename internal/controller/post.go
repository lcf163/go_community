package controller

import (
	"errors"
	"go_community/internal/dao/mysql"
	"go_community/internal/models"
	"go_community/internal/service"
	"strconv"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// CreatePostHandler
// @Summary 创建帖子
// @Description 创建新帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param post body models.ParamPost true "帖子信息"
// @Success 1000 {object} ResponseData
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /post [post]
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
	p.AuthorID = userID
	// 2.创建帖子
	if err := service.CreatePost(p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}

// PostDetailHandler
// @Summary 获取帖子详情
// @Description 获取帖子的详细信息
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param id path int true "帖子ID"
// @Success 1000 {object} _ResponsePostDetail
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /post/{id} [get]
func PostDetailHandler(c *gin.Context) {
	// 1.获取参数（URL中获取帖子ID）
	postIDStr := c.Param("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("PostDetailHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
	}
	// 2.根据ID取出帖子数据（查数据库）
	post, err := service.GetPostById(postID)
	if err != nil {
		zap.L().Error("logic.GetPostById failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 3.返回响应
	ResponseSuccess(c, post)
}

// GetPostListHandler
// @Summary 获取帖子列表
// @Description 分页获取帖子列表
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param page query int false "页码" minimum(1) default(1)
// @Param size query int false "每页数量" minimum(1) maximum(10) default(5)
// @Success 1000 {object} _ResponsePostList
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /posts [get]
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	// 获取数据
	posts, err := service.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPostList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 返回响应
	ResponseSuccess(c, posts)
}

// GetPostListHandler2
// @Summary 升级版帖子列表接口
// @Description 分页获取帖子列表（按帖子的创建时间或者分数排序）
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param object query models.ParamPostListQueryNoSearch false "查询参数"
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
	posts, err := service.GetPostListNew(p) // 更新：合二为一
	if err != nil {
		zap.L().Error("logic.GetPostListNew failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return // 添加return，避免错误情况下继续执行
	}

	// 返回响应
	ResponseSuccess(c, posts)
}

// PostSearchHandler
// @Summary 搜索帖子
// @Description 根据关键词搜索帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param object query models.ParamPostListQueryWithSearch false "查询参数"
// @Success 1000 {object} _ResponsePostList
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /search [get]
func PostSearchHandler(c *gin.Context) {
	p := &models.ParamPostList{}
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("PostSearchHandler with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取数据
	data, err := service.PostSearch(p)
	if err != nil {
		zap.L().Error("logic.PostSearch failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, data)
}

// UpdatePostHandler
// @Summary 更新帖子
// @Description 更新帖子内容
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param post body models.ParamUpdatePost true "帖子更新信息"
// @Success 1000 {object} ResponseData
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1011 {object} ResponseData "无操作权限"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /post [put]
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
	if err := service.UpdatePost(userID, p); err != nil {
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

// GetUserPostListHandler
// @Summary 获取用户帖子列表
// @Description 获取指定用户发布的帖子列表
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Param id path int true "用户ID"
// @Param page query int false "页码" minimum(1) default(1)
// @Param size query int false "每页数量" minimum(1) maximum(10) default(5)
// @Success 1000 {object} _ResponsePostList
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1003 {object} ResponseData "用户不存在"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /posts/user/{id} [get]
func GetUserPostListHandler(c *gin.Context) {
	// 获取用户ID参数
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取分页参数
	page, size := getPageInfo(c)

	// 获取数据
	data, err := service.GetUserPostList(userID, page, size)
	if err != nil {
		zap.L().Error("logic.GetUserPostList failed",
			zap.Int64("user_id", userID),
			zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}

// DeletePostHandler
// @Summary 删除帖子
// @Description 删除指定帖子
// @Tags 帖子相关接口
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "帖子ID"
// @Success 1000 {object} ResponseData
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1011 {object} ResponseData "无操作权限"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /post/{id} [delete]
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
	if err := service.DeletePost(userID, postID); err != nil {
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
