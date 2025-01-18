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
		zap.String("targetId", p.TargetId),
		zap.Int8("targetType", p.TargetType),
		zap.Int8("direction", p.Direction))

	userIdStr := strconv.Itoa(int(userId))
	switch p.TargetType {
	case TypePost:
		return redis.VoteForPost(userIdStr, p.TargetId, float64(p.Direction))
	case TypeComment:
		return redis.VoteForComment(userIdStr, p.TargetId, float64(p.Direction))
	default:
		return 0, errors.New("无效的投票目标类型")
	}
}
