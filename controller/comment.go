package controller

import (
	"go_community/logic"
	"go_community/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateCommentHandler 创建评论
func CreateCommentHandler(c *gin.Context) {
	// 参数校验
	p := new(models.ParamComment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreateCommentHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取当前用户ID
	userID, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	// 创建评论
	if err := logic.CreateComment(userID, p); err != nil {
		zap.L().Error("logic.CreateComment failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

// GetCommentListHandler 获取评论列表
func GetCommentListHandler(c *gin.Context) {
	// 获取帖子ID
	postIDStr := c.Param("postId")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取评论列表
	data, err := logic.GetCommentList(postID)
	if err != nil {
		zap.L().Error("logic.GetCommentList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}
