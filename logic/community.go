package logic

import (
	"errors"
	"go_community/dao/mysql"
	"go_community/models"
	"go_community/pkg/snowflake"

	"go.uber.org/zap"
)

// GetCommunityList 查询分类社区列表
func GetCommunityList() ([]*models.Community, error) {
	// 数据库中查找到所有的 community 并返回
	return mysql.GetCommunityList()
}

// GetCommunityList2 查询分类社区列表（带分页）
func GetCommunityList2(p *models.ParamPage) (*models.ApiCommunityDetailRes, error) {
	// 获取总数
	total, err := mysql.GetCommunityTotalCount()
	if err != nil {
		zap.L().Error("mysql.GetCommunityTotalCount() failed", zap.Error(err))
		return nil, err
	}

	// 获取分页数据
	communities, err := mysql.GetCommunityList2(p)
	if err != nil {
		zap.L().Error("mysql.GetCommunityList2(p) failed",
			zap.Error(err),
			zap.Any("params", p))
		return nil, err
	}

	// 组装返回数据
	res := &models.ApiCommunityDetailRes{
		Page: &models.Page{
			Page:  p.Page,
			Size:  p.Size,
			Total: total,
		},
		List: communities,
	}

	return res, nil
}

// GetCommunityDetailById 根据ID查询社区详情
func GetCommunityDetailById(communityID int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailById(communityID)
}

// CreateCommunity 创建社区
func CreateCommunity(userID int64, community *models.CommunityDetail) error {
	// 检查社区名称是否已存在
	exists, err := mysql.GetCommunityDetailByName(community.CommunityName)
	if err != nil && err != mysql.ErrorInvalidID {
		zap.L().Error("mysql.GetCommunityDetailByName(community.CommunityName) failed",
			zap.String("name", community.CommunityName),
			zap.Error(err))
		return err
	}
	if exists != nil {
		return errors.New("社区名称已存在")
	}

	// 生成社区ID
	community.CommunityID = snowflake.GetID()

	// 创建社区
	if err := mysql.CreateCommunity(community); err != nil {
		zap.L().Error("mysql.CreateCommunity(community) failed",
			zap.Any("community", community),
			zap.Error(err))
		return err
	}
	return nil
}

// UpdateCommunity 更新社区信息
func UpdateCommunity(userID, communityID int64, communityName, introduction string) error {
	// 检查参数
	if communityName == "" && introduction == "" {
		return errors.New("更新的内容不能为空")
	}

	// 检查社区是否存在
	existingCommunity, err := mysql.GetCommunityDetailById(communityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailById(communityID) failed",
			zap.Int64("community_id", communityID),
			zap.Error(err))
		return err
	}

	// 如果不修改名称，使用原来的名称
	if communityName == "" {
		communityName = existingCommunity.CommunityName
	}
	// 如果不修改简介，使用原来的简介
	if introduction == "" {
		introduction = existingCommunity.Introduction
	}

	// 如果要修改名称，检查新名称是否已存在
	if communityName != existingCommunity.CommunityName {
		existingCommunity, err := mysql.GetCommunityDetailByName(communityName)
		if err != nil && err != mysql.ErrorInvalidID {
			zap.L().Error("mysql.GetCommunityDetailByName failed",
				zap.String("community_name", communityName),
				zap.Error(err))
			return err
		}
		// 如果找到同名社区，且不是当前社区
		if existingCommunity != nil && existingCommunity.CommunityID != communityID {
			return errors.New("社区名称已存在")
		}
	}

	// 更新社区信息
	if err := mysql.UpdateCommunity(userID, communityID, communityName, introduction); err != nil {
		zap.L().Error("mysql.UpdateCommunity failed",
			zap.Int64("community_id", communityID),
			zap.String("community_name", communityName),
			zap.String("introduction", introduction),
			zap.Error(err))
		return err
	}
	return nil
}

// DeleteCommunity 删除社区
func DeleteCommunity(userID int64, communityID int64) error {
	// 检查社区是否存在
	_, err := mysql.GetCommunityDetailById(communityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailById(communityID) failed",
			zap.Int64("community_id", communityID),
			zap.Error(err))
		return err
	}

	// 检查社区下是否有帖子
	count, err := mysql.GetCommunityPostTotalCount(communityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityPostTotalCount(communityID) failed",
			zap.Int64("community_id", communityID),
			zap.Error(err))
		return err
	}
	if count > 0 {
		return errors.New("该社区下还有帖子，无法删除")
	}

	// 删除社区
	if err := mysql.DeleteCommunity(communityID); err != nil {
		zap.L().Error("mysql.DeleteCommunity(communityID) failed",
			zap.Int64("community_id", communityID),
			zap.Error(err))
		return err
	}
	return nil
}
