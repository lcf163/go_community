package controller

import (
	"go-community/dao/redis"
	"go-community/logic"
	"go-community/models"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

func VoteHandler(c *gin.Context) {
	// 1.获取请求参数和参数校验
	vote := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(vote); err != nil {
		zap.L().Error("VoteHandler with invalid param", zap.Error(err))
		errs, ok := err.(validator.ValidationErrors) // 类型断言
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		// 翻译并去除错误提示中的结构体标识
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParams, errData)
		return
	}
	// 从 c 取到当前发请求的用户 ID
	userID, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}
	// 2.业务逻辑处理
	if err := logic.VoteForPost(userID, vote); err != nil {
		zap.L().Error("logic.VoteForPost failed", zap.Error(err))
		switch err {
		case redis.ErrorVoteRepeted: // 重复投票
			ResponseError(c, CodeVoteRepeated)
		case redis.ErrorVoteTimeExpire: // 投票超时
			ResponseError(c, CodeVoteTimeExpire)
		default:
			ResponseError(c, CodeServerBusy)
		}
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}
