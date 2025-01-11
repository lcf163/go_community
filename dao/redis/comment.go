package redis

import (
	"github.com/go-redis/redis"
	"time"
)

// CreateComment 创建评论时记录到Redis
func CreateComment(commentId int64) error {
	now := float64(time.Now().Unix())
	pipeline := client.TxPipeline()

	// 记录评论时间
	pipeline.ZAdd(getRedisKey(KeyCommentTimeZSet), redis.Z{
		Score:  now,
		Member: commentId,
	})

	_, err := pipeline.Exec()
	return err
}