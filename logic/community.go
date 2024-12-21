package logic

import (
	"go-community/dao/mysql"
	"go-community/models"
)

// GetCommunityList 查询分类社区列表
func GetCommunityList() ([]*models.Community, error) {
	// 数据库中查找到所有的 community 并返回
	return mysql.GetCommunityList()
}

// GetCommunityDetailByID 根据ID查询社区详情
func GetCommunityDetailByID(id uint64) (*models.CommunityDetail, error) {
	return mysql.GetCommunityDetailByID(id)
}
