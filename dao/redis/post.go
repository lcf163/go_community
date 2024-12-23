package redis

import (
	"go-community/models"
)

// getIdsFormKey 按照分数从大到小的顺序，查询指定数量的元素
func getIdsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// 3.ZRevRange 按照分数从大到小的顺序，查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}

// GetPostIdInOrder 获取帖子列表：按创建时间排序或者分数排序（查询出 ids 根据 order 从大到小排序）
func GetPostIdInOrder(p *models.ParamPostList) ([]string, error) {
	// 从 redis 获取 id
	// 1.根据请求中携带的 order 参数，确定要查询的 redis key
	key := getRedisKey(KeyPostTimeZSet) // 默认是时间
	if p.Order == models.OrderScore {   // 按照分数请求
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 2.确定查询的索引起始点
	return getIdsFormKey(key, p.Page, p.Size)
}
