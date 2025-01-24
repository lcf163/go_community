package mysql

import (
	"database/sql"
	"go_community/models"

	"go.uber.org/zap"
)

// GetCommunityList 查询分类社区列表
func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community where status = 1"
	err = db.Select(&communityList, sqlStr)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return
}

// GetCommunityList2 获取社区列表（带分页）
func GetCommunityList2(p *models.ParamPage) (communities []*models.CommunityDetail, err error) {
	sqlStr := `select community_id, community_name, introduction
	from community
	where status = 1
	ORDER BY create_time
	DESC
	limit ?,?`

	err = db.Select(&communities, sqlStr, (p.Page-1)*p.Size, p.Size)
	if err != nil {
		return nil, err
	}
	return
}

// GetCommunityTotalCount 查询分类社区总数
func GetCommunityTotalCount() (count int64, err error) {
	sqlStr := `select count(community_id) from community where status = 1`
	err = db.Get(&count, sqlStr)
	if err != nil {
		return 0, err
	}
	return
}

// GetCommunityDetailById 根据ID查询社区详情
func GetCommunityDetailById(id int64) (community *models.CommunityDetail, err error) {
	community = new(models.CommunityDetail)
	sqlStr := `select community_id, community_name, introduction, create_time
	from community
	where community_id = ? and status = 1`
	err = db.Get(community, sqlStr, id)
	if err == sql.ErrNoRows {
		err = ErrorInvalidID
		return
	}
	if err != nil {
		zap.L().Error("query community failed",
			zap.String("sql", sqlStr),
			zap.Error(err))
		err = ErrorQueryFailed
		return
	}
	return
}

// CreateCommunity 创建社区
func CreateCommunity(community *models.CommunityDetail) (err error) {
	// 设置默认状态为1
	community.Status = 1
	sqlStr := `insert into community(
	community_id, community_name, introduction, status)
	values(?,?,?,?)`
	_, err = db.Exec(sqlStr, community.CommunityId,
		community.CommunityName, community.Introduction, community.Status)
	if err != nil {
		zap.L().Error("CreateCommunity failed",
			zap.String("sql", sqlStr),
			zap.Any("community", community),
			zap.Error(err))
		err = ErrorInsertFailed
		return
	}
	return
}

// UpdateCommunity 更新社区信息
func UpdateCommunity(id int64, name, introduction string) error {
	sqlStr := `update community 
	set community_name = ?, introduction = ? 
	where community_id = ? and status = 1`
	result, err := db.Exec(sqlStr, name, introduction, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrorInvalidID
	}
	return nil
}

// DeleteCommunity 删除社区（软删除）
func DeleteCommunity(id int64) error {
	sqlStr := `update community set status = 0 where community_id = ? and status = 1`
	result, err := db.Exec(sqlStr, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrorInvalidID
	}
	return nil
}
