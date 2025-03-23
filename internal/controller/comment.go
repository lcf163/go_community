package controller

import (
	"go_community/internal/dao/mysql"
	"go_community/internal/models"
	"go_community/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreateCommentHandler
// @Summary 创建评论/回复
// @Description 创建评论或回复评论
// @Tags 评论相关接口
// @Accept application/json
// @Produce application/json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param comment body models.ParamComment true "评论信息"
// @Success 1000 {object} ResponseData
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /comment [post]
func CreateCommentHandler(c *gin.Context) {
	// 参数校验
	p := new(models.ParamComment)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("CreateCommentHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 如果是回复,检查必填字段
	if p.ParentID != 0 && p.ReplyToUID == 0 {
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
	if err := service.CreateComment(userID, p); err != nil {
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

// UpdateCommentHandler
// @Summary 更新评论
// @Description 更新评论内容
// @Tags 评论相关接口
// @Accept application/json
// @Produce application/json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param comment body models.ParamUpdateComment true "更新评论信息"
// @Success 1000 {object} ResponseData
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1011 {object} ResponseData "无操作权限"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /comment [put]
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
	if err := service.UpdateComment(userID, p); err != nil {
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

// GetCommentListHandler
// @Summary 获取评论列表
// @Description 获取帖子评论或评论回复列表
// @Tags 评论相关接口
// @Accept application/json
// @Produce application/json
// @Param post_id query int false "帖子ID(获取帖子评论时必填)"
// @Param comment_id query int false "评论ID(获取评论回复时必填)"
// @Param page query int false "页码" minimum(1) default(1)
// @Param size query int false "每页数量" minimum(1) maximum(100) default(10)
// @Success 1000 {object} _ResponseCommentList
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /comments [get]
func GetCommentListHandler(c *gin.Context) {
	// 获取参数
	p := &models.ParamCommentList{}
	if err := c.ShouldBindQuery(p); err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 参数校验
	if p.PostID == 0 && p.CommentID == 0 {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 根据参数判断获取类型
	var data interface{}
	var err error
	if p.PostID != 0 {
		// 获取帖子评论列表
		data, err = service.GetCommentList(p.PostID, p.Page, p.Size)
	} else {
		// 获取评论回复列表
		data, err = service.GetCommentReplyList(p.CommentID)
	}

	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, data)
}

// GetCommentDetailHandler
// @Summary 获取评论详情
// @Description 获取评论的详细信息
// @Tags 评论相关接口
// @Accept application/json
// @Produce application/json
// @Param id path int true "评论ID"
// @Success 1000 {object} _ResponseCommentDetail
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /comment/{id} [get]
func GetCommentDetailHandler(c *gin.Context) {
	// 获取评论ID参数
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取评论详情
	data, err := service.GetCommentById(commentID)
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

// DeleteCommentHandler
// @Summary 删除评论
// @Description 删除指定评论
// @Tags 评论相关接口
// @Accept application/json
// @Produce application/json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "评论ID"
// @Success 1000 {object} ResponseData
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1011 {object} ResponseData "无操作权限"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /comment/{id} [delete]
func DeleteCommentHandler(c *gin.Context) {
	// 1. 获取评论ID
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
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

	// 3. 删除评论
	if err := service.DeleteComment(userID, commentID); err != nil {
		zap.L().Error("logic.DeleteComment failed",
			zap.Int64("comment_id", commentID),
			zap.Int64("user_id", userID),
			zap.Error(err))

		// 根据错误类型返回对应的错误码
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

// DeleteCommentWithRepliesHandler
// @Summary 删除评论及其回复
// @Description 删除指定评论及其所有回复
// @Tags 评论相关接口
// @Accept application/json
// @Produce application/json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "评论ID"
// @Success 1000 {object} ResponseData
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1011 {object} ResponseData "无操作权限"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /comments/{id} [delete]
func DeleteCommentWithRepliesHandler(c *gin.Context) {
	// 1. 获取评论ID
	commentIDStr := c.Param("id")
	commentID, err := strconv.ParseInt(commentIDStr, 10, 64)
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

	// 3. 删除评论及其回复
	if err := service.DeleteCommentWithReplies(userID, commentID); err != nil {
		zap.L().Error("logic.DeleteCommentWithReplies failed",
			zap.Int64("comment_id", commentID),
			zap.Int64("user_id", userID),
			zap.Error(err))

		if err == mysql.ErrorInvalidID {
			ResponseErrorWithMsg(c, CodeInvalidParams, "评论不存在或已删除")
			return
		}
		if err == mysql.ErrorNoPermission {
			ResponseErrorWithMsg(c, CodeNoPermission, "无权删除该评论")
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}
