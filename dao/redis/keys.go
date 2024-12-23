package redis

// redis key

// redis key 注意使用命名空间的方式，方便查询和拆分
const (
	KeyPrefix              = "go-community:"
	KeyPostTimeZSet        = "post:time"   // 帖子及发帖时间
	KeyPostScoreZSet       = "post:score"  // 帖子及投票分数
	KeyPostVotedZSetPrefix = "post:voted:" // 记录用户及投票类型，参数是 postId
)

// getRedisKey redis key 拼接前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
