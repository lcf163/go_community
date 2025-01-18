package controller

import (
	"go_community/dao/redis"
	"go_community/logic"
	"go_community/models"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// VoteHandler 投票处理
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

	// 获取当前用户
	userID, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	// 投票并获取最新点赞数
	voteNum, err := logic.VoteForTarget(userID, vote)
	if err != nil {
		zap.L().Error("logic.VoteForTarget failed", zap.Error(err))
		switch err {
		case redis.ErrorVoteRepeted:
			ResponseError(c, CodeVoteRepeated)
		case redis.ErrorVoteTimeExpire:
			ResponseError(c, CodeVoteTimeExpire)
		default:
			ResponseError(c, CodeServerBusy)
		}
		return
	}

	ResponseSuccess(c, gin.H{
		"vote_num": voteNum,
	})
}
