package redis

import (
	"go_community/models"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// getIdsFormKey 按照分数从大到小的顺序，查询指定数量的元素
func getIdsFormKey(key string, page, size int64) ([]string, error) {
	start := (page - 1) * size
	end := start + size - 1
	// 3.ZRevRange 按照分数从大到小的顺序，查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}

// GetPostIdsInOrder 获取帖子列表：按创建时间排序或者分数排序（查询出 ids 根据 order 从大到小排序）
func GetPostIdsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从 redis 获取 id
	// 1.根据请求中携带的 order 参数，确定要查询的 redis key
	key := getRedisKey(KeyPostTimeZSet) // 默认是时间
	if p.Order == models.OrderScore {   // 按照分数请求
		key = getRedisKey(KeyPostScoreZSet)
	}
	// 2.确定查询的索引起始点
	return getIdsFormKey(key, p.Page, p.Size)
}

// GetPostVoteNum 获取单个帖子的点赞数
func GetPostVoteNum(postId string) (voteNum int64, err error) {
	// 构造存储投票数据的key
	key := getRedisKey(KeyPostVotedZSetPrefix + postId)

	// 查找key中分数是1的元素数量，即点赞数量
	return client.ZCount(key, "1", "1").Result()
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

// GetCommunityPostIdsInOrder 按社区查询ids(查询出的ids根据order从大到小排序)
func GetCommunityPostIdsInOrder(p *models.ParamPostList) ([]string, error) {
	// 从 redis 获取 id
	// 1.根据请求中携带的 order 参数，确定要查询的 redis key
	orderKey := getRedisKey(KeyPostTimeZSet) // 默认是时间
	if p.Order == models.OrderScore {        // 按照分数请求
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	// 使用 zinterstore: 把分区的帖子 set 与帖子分数 zset 生成一个新的 zset
	// 针对新的 zset，按之前的逻辑取数据
	// 利用缓存 key 减少 zinterstore 执行的次数
	// 社区的 key
	communityKey := getRedisKey(KeyCommunityPostSetPrefix + strconv.Itoa(int(p.CommunityId)))
	// 缓存的 key
	key := orderKey + strconv.Itoa(int(p.CommunityId))
	if client.Exists(key).Val() < 1 {
		// 不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX", // 将两个 zset 函数聚合时求最大值
		}, communityKey, orderKey) // zinterstore 计算
		pipeline.Expire(key, 60*time.Second) // 设置超时时间
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	// 2.确定查询的索引起始点
	return getIdsFormKey(key, p.Page, p.Size)
}
