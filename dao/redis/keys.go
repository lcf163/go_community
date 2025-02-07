package redis

// redis key

// redis key 注意使用命名空间的方式，方便查询和拆分
const (
	KeyPrefix                 = "go_community:"
	KeyPostTimeZSet           = "post:time"      // 帖子及发帖时间
	KeyPostScoreZSet          = "post:score"     // 帖子及投票分数
	KeyPostVotedZSetPrefix    = "post:voted:"    // 记录用户及投票类型
	KeyCommunityPostSetPrefix = "community:"     // 保存每个分区下帖子的id
	KeyCommentTimeZSet        = "comment:time"   // 评论及发布时间
	KeyCommentVotedZSetPrefix = "comment:voted:" // 记录用户为评论投票的数据
)

// getRedisKey redis key 拼接前缀
func getRedisKey(key string) string {
	return KeyPrefix + key
}
