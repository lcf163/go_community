package mysql

import (
	"database/sql"
	"go_community/models"

	"go.uber.org/zap"
)

// GetCommunityList 查询分类社区列表
func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	err = db.Select(&communityList, sqlStr)
	if err == sql.ErrNoRows { // 查询为空
		zap.L().Warn("there is no community in db")
		err = nil
	}
	return
}

// GetCommunityList2 获取社区列表
func GetCommunityList2(p *models.ParamPage) ([]*models.CommunityDetail, error) {
	sqlStr := `select community_id, community_name, introduction
	from community
	ORDER BY create_time
	DESC
	limit ?,?`

	communities := make([]*models.CommunityDetail, 0, 2)
	err := db.Select(&communities, sqlStr, (p.Page-1)*p.Size, p.Size)

	// 添加日志
	zap.L().Debug("GetCommunityList SQL",
		zap.String("sql", sqlStr),
		zap.Int64("page", p.Page),
		zap.Int64("size", p.Size))

	return communities, err
}

// GetCommunityTotalCount 查询分类社区总数
func GetCommunityTotalCount() (count int64, err error) {
	sqlStr := `select count(community_id) from community`
	err = db.Get(&count, sqlStr)
	if err != nil {
		zap.L().Error("db.Get(&count, sqlStr) failed", zap.Error(err))
		return 0, err
	}
	return
}

// GetCommunityDetailById 根据ID查询社区详情
func GetCommunityDetailById(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time
	from community
	where community_id = ?`
	err = db.Get(community, sqlStr, id)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		zap.L().Error("query community failed", zap.String("sql", sqlStr), zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}
