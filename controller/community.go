package controller

import (
	"go_community/dao/mysql"
	"go_community/logic"
	"go_community/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 社区相关

// CommunityHandler
// @Summary 获取社区列表
// @Description 获取所有社区的信息列表
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Success 1000 {object} ResponseData{data=[]models.Community}
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /community [get]
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

// CommunityHandler2
// @Summary 获取社区列表
// @Description 获取社区列表（带分页）
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param page query int false "页码" minimum(1) default(1)
// @Param size query int false "每页数量" minimum(1) maximum(100) default(10)
// @Success 1000 {object} ResponseData{data=[]models.CommunityDetail}
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /community2 [get]
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
// @Summary 获取社区详情
// @Description 根据社区ID获取社区详细信息
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Param id path int true "社区ID"
// @Success 1000 {object} ResponseData{data=models.CommunityDetail}
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1009 {object} ResponseData "社区不存在"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /community/{id} [get]
func CommunityDetailHandler(c *gin.Context) {
	// 1.获取社区ID
	communityIDStr := c.Param("id")
	communityID, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 2.根据ID获取社区详情
	communityList, err := logic.GetCommunityDetailById(communityID)
	if err != nil {
		zap.L().Error("logic.GetCommunityDetailById failed", zap.Error(err))
		if err == mysql.ErrorInvalidID {
			ResponseError(c, CodeCommunityNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, communityList)
}

// CreateCommunityHandler 创建社区
// @Summary 创建社区
// @Description 创建新的社区
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param community body models.ParamUpdateCommunity true "社区信息"
// @Success 1000 {object} ResponseData
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1010 {object} ResponseData "社区已存在"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /community [post]
func CreateCommunityHandler(c *gin.Context) {
	// 检查用户权限
	userID, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	// 获取参数
	community := new(models.CommunityDetail)
	if err := c.ShouldBindJSON(community); err != nil {
		zap.L().Error("CreateCommunityHandler with invalid param",
			zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, "请检查社区名称和简介是否为空")
		return
	}

	// 创建社区
	if err := logic.CreateCommunity(userID, community); err != nil {
		zap.L().Error("logic.CreateCommunity failed",
			zap.Any("community", community),
			zap.Error(err))
		if err.Error() == "社区名称已存在" {
			ResponseError(c, CodeCommunityExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// UpdateCommunityHandler 更新社区信息
// @Summary 更新社区
// @Description 更新社区信息
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "社区ID"
// @Param community body models.ParamUpdateCommunity true "社区更新信息"
// @Success 1000 {object} ResponseData
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1009 {object} ResponseData "社区不存在"
// @Failure 1010 {object} ResponseData "社区名称已存在"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /community/{id} [put]
func UpdateCommunityHandler(c *gin.Context) {
	// 检查用户权限
	userID, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	// 获取社区ID
	communityIDStr := c.Param("id")
	communityID, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParams, "无效的社区ID")
		return
	}

	// 获取更新参数
	p := new(models.ParamUpdateCommunity)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateCommunityHandler with invalid param",
			zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, "请检查社区名称和简介格式是否正确")
		return
	}

	// 更新社区信息
	if err := logic.UpdateCommunity(userID, communityID, p.Name, p.Introduction); err != nil {
		zap.L().Error("logic.UpdateCommunity failed",
			zap.Int64("communityID", communityID),
			zap.Any("params", p),
			zap.Error(err))
		if err.Error() == "社区名称已存在" {
			ResponseError(c, CodeCommunityExist)
			return
		}
		if err == mysql.ErrorInvalidID {
			ResponseError(c, CodeCommunityNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}

// DeleteCommunityHandler 删除社区
// @Summary 删除社区
// @Description 删除指定社区
// @Tags 社区相关接口
// @Accept application/json
// @Produce application/json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param id path int true "社区ID"
// @Success 1000 {object} ResponseData
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1009 {object} ResponseData "社区不存在"
// @Failure 1014 {object} ResponseData "社区下存在帖子"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /community/{id} [delete]
func DeleteCommunityHandler(c *gin.Context) {
	// 检查用户权限
	userID, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	// 获取社区ID
	communityIDStr := c.Param("id")
	communityID, err := strconv.ParseInt(communityIDStr, 10, 64)
	if err != nil {
		ResponseErrorWithMsg(c, CodeInvalidParams, "无效的社区ID")
		return
	}

	// 2. 删除社区
	if err := logic.DeleteCommunity(userID, communityID); err != nil {
		zap.L().Error("logic.DeleteCommunity failed",
			zap.Int64("communityID", communityID),
			zap.Error(err))
		if err == mysql.ErrorInvalidID {
			ResponseError(c, CodeCommunityNotExist)
			return
		}
		if err.Error() == "该社区下还有帖子，无法删除" {
			ResponseError(c, CodeCommunityHasPost)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, nil)
}
