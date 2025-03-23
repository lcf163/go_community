package controller

import (
	"github.com/go-playground/validator/v10"
	"go_community/internal/dao/redis"
	"go_community/internal/models"
	"go_community/internal/service"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// VoteHandler 投票处理
// @Summary 投票
// @Description 对帖子或评论进行投票
// @Tags 投票相关接口
// @Accept application/json
// @Produce application/json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param vote body models.ParamVoteData true "投票信息"
// @Success 1000 {object} ResponseData{data=map[string]int64{vote_num=int64}}
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1012 {object} ResponseData "重复投票"
// @Failure 1013 {object} ResponseData "投票时间已过"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /vote [post]
func VoteHandler(c *gin.Context) {
	// 1.获取请求参数和参数校验
	vote := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(vote); err != nil {
		zap.L().Error("VoteHandler with invalid param", zap.Error(err))

		// 处理 validator.ValidationErrors 类型的错误
		if errs, ok := err.(validator.ValidationErrors); ok {
			errData := removeTopStruct(errs.Translate(trans))
			ResponseErrorWithMsg(c, CodeInvalidParams, errData)
			return
		}

		// 处理 UnmarshalJSON 返回的错误
		ResponseErrorWithMsg(c, CodeInvalidParams, err.Error())
		return
	}

	// 获取当前用户
	userID, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	// 投票并获取最新点赞数
	voteNum, err := service.VoteForTarget(userID, vote)
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
