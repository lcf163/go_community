package controller

import (
	"go_community/dao/mysql"
	"go_community/logic"
	"go_community/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateCommentHandler 创建评论/回复
func CreateCommentHandler(c *gin.Context) {
	// 参数校验
	p := new(models.ParamComment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreateCommentHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 如果是回复,检查必填字段
	if p.ParentId != 0 && p.ReplyToUid == 0 {
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
		zap.L().Error("logic.CreateComment failed",
			zap.Error(err),
			zap.Any("params", p))
		if err == mysql.ErrorInvalidID {
			ResponseError(c, CodeInvalidParams)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

// UpdateCommentHandler 更新评论
func UpdateCommentHandler(c *gin.Context) {
	// 参数校验
	p := new(models.ParamUpdateComment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateCommentHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取当前用户ID
	userID, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	// 更新评论
	if err := logic.UpdateComment(userID, p); err != nil {
		zap.L().Error("logic.UpdateComment failed",
			zap.Error(err),
			zap.Any("params", p))

		if err == mysql.ErrorInvalidID {
			ResponseError(c, CodeInvalidParams)
			return
		}
		if err == mysql.ErrorNoPermission {
			ResponseError(c, CodeNoPermission)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

// GetCommentListHandler 获取评论列表（支持获取帖子评论和评论回复）
func GetCommentListHandler(c *gin.Context) {
	// 获取参数
	p := &models.ParamCommentList{}
	if err := c.ShouldBindQuery(p); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 参数校验
	if p.PostId == 0 && p.CommentId == 0 {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 根据参数判断获取类型
	var data interface{}
	var err error
	if p.PostId != 0 {
		// 获取帖子评论列表
		data, err = logic.GetCommentList(p.PostId, p.Page, p.Size)
	} else {
		// 获取评论回复列表
		data, err = logic.GetCommentReplyList(p.CommentId)
	}

	if err != nil {
		zap.L().Error("get comments failed",
			zap.Any("params", p),
			zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}

// GetCommentDetailHandler 获取单个评论详情
func GetCommentDetailHandler(c *gin.Context) {
	// 获取评论ID参数
	commentIDStr := c.Param("commentId")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取评论详情
	data, err := logic.GetCommentById(commentID)
	if err != nil {
		zap.L().Error("logic.GetCommentById failed",
			zap.Int64("comment_id", commentID),
			zap.Error(err))
		if err == mysql.ErrorInvalidID {
			ResponseError(c, CodeInvalidParams)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}
