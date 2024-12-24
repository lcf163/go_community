package redis

import (
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	OneWeekInSeconds          = 7 * 24 * 3600        // 一周的秒数
	OneMonthInSeconds         = 4 * OneWeekInSeconds // 一个月的秒数
	VoteScore         float64 = 432                  // 每一票的值432分
	PostPerAge                = 20                   // 每页显示20条帖子
)

/*
	推荐阅读
	基于用户投票的排名算法：http://www.ruanyifeng.com/blog/algorithm/

	本项目使用简化版的投票分数
	投一票加432分（86400/200=432，200张赞成票就可以给帖子在首页续天）->《redis实战》
*/

/*
	为帖子投票
	投票分为四种情况：1.投赞成票(1) 2.投反对票(-1) 3.取消投票(0) 4.反转投票

	记录文章参与投票的人
	更新文章分数：赞成票要加分；反对票减分

	direction=1 时，有两种情况
		1.之前没有投过票，现在投赞成票  -> 更新分数和投票记录，差值的绝对值：1, +432
		2.之前投反对票，现在改为赞成票  -> 更新分数和投票记录，差值的绝对值：2, +432*2
	direction=0 时，有两种情况
		2.之前投过反对票，现在取消投票  -> 更新分数和投票记录，差值的绝对值：1, +432
		1.之前投过赞成票，现在取消投票  -> 更新分数和投票记录，差值的绝对值：1, -432
	direction=-1 时，有两种情况
		1.之前没有投过票，现在投反对票  -> 更新分数和投票记录，差值的绝对值：1, -432
		2.之前投赞成票，现在改为反对票  -> 更新分数和投票记录，差值的绝对值：2, -432*2

	投票的限制：
		每个帖子子发表之日起一个星期之内允许用户投票，超过一个星期就不允许投票了
		1.到期之后将 redis 中保存的赞成票数及反对票数存储到 mysql 表中
		2.到期之后删除那个 KeyPostVotedZSetPrefix
*/

// VoteForPost	为帖子投票
func VoteForPost(userId, postId string, v float64) (err error) {
	// 1.判断投票限制
	// 去redis取帖子发布时间
	postTime := client.ZScore(getRedisKey(KeyPostTimeZSet), postId).Val()
	if float64(time.Now().Unix())-postTime > OneWeekInSeconds {
		// 超过一个星期，就不允许投票了
		return ErrorVoteTimeExpire
	}
	// 2和3，需要放到一个事务中操作
	// 2.更新帖子的分数
	// 判断是否已经投过票，查当前用户给当前帖子的投票记录
	key := getRedisKey(KeyPostVotedZSetPrefix) + postId
	ov := client.ZScore(key, userId).Val()
	// 若这一次投票的值和之前保存的值一致，则提示不允许重复投票
	if v == ov {
		return ErrorVoteRepeted
	}
	var op float64
	if v > ov {
		op = 1
	} else {
		op = -1
	}
	diffAbs := math.Abs(ov - v)                // 计算两次投票的差值
	incrementScore := VoteScore * diffAbs * op // 计算分数（新增）
	pipeline := client.TxPipeline()            // 事务操作
	// ZIncrBy 用于将有序集合中的成员分数增加指定数量
	pipeline.ZIncrBy(getRedisKey(KeyPostScoreZSet), incrementScore, postId) // 更新分数
	// 3.记录用户为该帖子投票的数据
	if v == 0 {
		pipeline.ZRem(key, userId)
	} else {
		pipeline.ZAdd(key, redis.Z{ // 记录已投票
			Score:  v, // 赞成票还是反对票
			Member: userId,
		})
	}
	// 4、更新帖子的投票数
	//pipeline.HIncrBy(getRedisKey(KeyPostVotedZSetPrefix+postId), "votes", int64(op))

	_, err = pipeline.Exec() // 事务操作的提交
	return err
}

// CreatePost redis 存储帖子信息
func CreatePost(postId, communityId int64) (err error) {
	now := float64(time.Now().Unix())
	pipeline := client.TxPipeline() // 事务操作
	// 文章 hash
	//pipeline.HMSet(getRedisKey(KeyPostVotedZSetPrefix+postId), postInfo)
	// 帖子时间 ZSet
	pipeline.ZAdd(getRedisKey(KeyPostTimeZSet), redis.Z{
		Score:  now,
		Member: postId,
	})
	// 帖子分数 ZSet
	pipeline.ZAdd(getRedisKey(KeyPostScoreZSet), redis.Z{
		Score:  now + VoteScore,
		Member: postId,
	})
	// 把帖子id添加到社区 set
	communityKey := getRedisKey(KeyCommunityPostSetPrefix) + strconv.Itoa(int(communityId))
	pipeline.SAdd(communityKey, postId)
	_, err = pipeline.Exec() // 事务操作的提交
	return
}
