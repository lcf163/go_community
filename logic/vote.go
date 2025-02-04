package logic

import (
	"errors"
	"go_community/dao/redis"
	"go_community/models"
	"strconv"

	"go.uber.org/zap"
)

const (
	TypePost    = 1 // 帖子
	TypeComment = 2 // 评论
)

// VoteForTarget 为帖子或评论投票
func VoteForTarget(userId int64, p *models.ParamVoteData) (voteNum int64, err error) {
	zap.L().Debug("VoteForTarget",
		zap.Int64("userId", userId),
		zap.Int64("targetId", p.TargetID),
		zap.Int8("targetType", p.TargetType),
		zap.Int8("direction", p.Direction))

	// 根据目标类型调用不同的投票函数
	switch p.TargetType {
	case TypePost:
		return redis.VoteForPost(
			strconv.FormatInt(userId, 10),
			strconv.FormatInt(p.TargetID, 10),
			float64(p.Direction))
	case TypeComment:
		return redis.VoteForComment(
			strconv.FormatInt(userId, 10),
			strconv.FormatInt(p.TargetID, 10),
			float64(p.Direction))
	default:
		return 0, errors.New("无效的投票目标类型")
	}
}
