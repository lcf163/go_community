package controller

import (
	"go_community/logic"
	"go_community/models"
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
		zap.L().Error("logic.GetCommunityList failed", zap.Error(err))
		ResponseError(c, CodeServerBusy) // 不轻易把服务端报错暴露给外面
		return
	}
	ResponseSuccess(c, communityList)
}

// CommunityHandler2 获取社区列表（带分页）
func CommunityHandler2(c *gin.Context) {
	// GET 请求参数（query string）: /api/v1/community?page=1&size=10
	p := &models.ParamPage{
		Page: 1,  // 默认第1页
		Size: 10, // 默认每页10条
	}
	// 获取分页参数
	// c.ShouldBind() 根据请求的数据类型，选择相应的方法去获取数据
	// c.ShouldBindJSON() 请求中携带 json 格式的数据，才能用这个方法获取到数据
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("CommunityHandler2 with invalid params", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 查询到所有的社区（community_id, community_name, introduction），以列表的形式返回
	communityList, err := logic.GetCommunityList2(p)
	if err != nil {
		zap.L().Error("logic.GetCommunityList2 failed", zap.Error(err))
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
	communityList, err := logic.GetCommunityDetailById(communityId)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetailById failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeSuccess, err.Error())
		return
	}
	ResponseSuccess(c, communityList)
}
