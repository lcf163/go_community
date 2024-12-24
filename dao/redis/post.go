package redis

import (
	"go-community/models"

	"github.com/go-redis/redis"
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

// GetPostVoteData 根据ids查询每篇帖子的投赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {
	//data = make([]int64, 0, len(ids))
	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPrefix + id)
	//	// 查找key中分数是1的元素数量 -> 统计每篇帖子的赞成票的数量
	//	v := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}

	// 使用 pipeline 一次发送多条命令，减少 RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1") // ZCount 会返回分数在 min 和 max 范围内的成员数量
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	data = make([]int64, 0, len(cmders))
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)
	}
	return
}
