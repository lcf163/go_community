package redis

import (
	"time"

	"github.com/go-redis/redis"
)

// GetCommentVoteNum 获取评论的投票数
func GetCommentVoteNum(commentId string) (voteNum int64, err error) {
	key := getRedisKey(KeyCommentVotedZSetPrefix + commentId)
	return client.ZCount(key, "1", "1").Result()
}

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

// DeleteCommentsVoteData 批量删除评论的点赞数据
func DeleteCommentsVoteData(commentIDs []string) error {
	pipeline := client.TxPipeline()
	
	// 批量删除评论的点赞数据和时间记录
	for _, commentID := range commentIDs {
		// 删除评论点赞记录
		pipeline.Del(getRedisKey(KeyCommentVotedZSetPrefix + commentID))
		// 删除评论时间记录
		pipeline.ZRem(getRedisKey(KeyCommentTimeZSet), commentID)
	}
	
	_, err := pipeline.Exec()
	return err
}
