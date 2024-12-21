package controller

import (
	"go-community/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 社区相关

// CommunityHandler 社区列表
func CommunityHandler(c *gin.Context) {
	// 查询到所有的社区（community_id, community_name），以列表的形式返回
	communityList, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("logic.GetCommunityList() failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错暴露给外面
		return
	}
	ResponseSuccess(c, communityList)
}

// CommunityDetailHandler 社区详情
func CommunityDetailHandler(c *gin.Context) {
	// 1.获取社区ID
	communityIdStr := c.Param("id")                              // 获取URL参数
	communityId, err := strconv.ParseInt(communityIdStr, 10, 64) // id字符串格式转换
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 2.根据ID获取社区详情
	communityList, err := logic.GetCommunityDetailById(int64(communityId))
	if err != nil {
		zap.L().Error("logic.GetCommunityByID() failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeSuccess, err.Error())
		return
	}
	ResponseSuccess(c, communityList)
}
