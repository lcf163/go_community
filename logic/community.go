package logic

import (
	"go_community/dao/mysql"
	"go_community/models"

	"go.uber.org/zap"
)

// GetCommunityList 查询分类社区列表
func GetCommunityList() ([]*models.Community, error) {
	// 数据库中查找到所有的 community 并返回
	return mysql.GetCommunityList()
}

// GetCommunityList2 查询分类社区列表（带分页）
func GetCommunityList2(p *models.ParamPage) (*models.ApiCommunityDetailRes, error) {
	// 参数校验和默认值设置
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size < 1 || p.Size > 100 {
		p.Size = 10
	}

	// 获取总数
	total, err := mysql.GetCommunityTotalCount()
	if err != nil {
		zap.L().Error("mysql.GetCommunityTotalCount failed", zap.Error(err))
		return nil, err
	}

	// 获取分页数据
	communities, err := mysql.GetCommunityList2(p)
	if err != nil {
		zap.L().Error("mysql.GetCommunityList failed",
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
func GetCommunityDetailById(id int64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailById(id)
}
