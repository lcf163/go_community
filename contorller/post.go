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
	userID, err := getCurrentUserID(c)
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
	}
	// 3.返回响应
	ResponseSuccess(c, post)
}
