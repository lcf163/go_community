package controller

import (
	"go-community/logic"
	"go-community/models"
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
		zap.L().Error("logic.GetPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 3.返回响应
	ResponseSuccess(c, post)
}

// GetPostListHandler 获取帖子列表
func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	page, size := getPageInfo(c)
	// 获取数据
	posts, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("logic.GetPost failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 返回响应
	ResponseSuccess(c, posts)
}

// GetPostListHandler2 获取帖子列表（按帖子的创建时间或者分数排序）
// 根据前端传来的参数动态地获取帖子列表
func GetPostListHandler2(c *gin.Context) {
	// GET 请求参数（query string）: /api/v1/post2?page=1&size=10&order=time
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
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
		zap.L().Error("logic.GetPostList2 failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 返回响应
	ResponseSuccess(c, posts)
}
