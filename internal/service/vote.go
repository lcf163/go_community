package service

import (
	"errors"
	"go_community/internal/dao/redis"
	"go_community/internal/models"
	"strconv"

	"go.uber.org/zap"
)

const (
	TypePost    = 1 // 帖子
	TypeComment = 2 // 评论
)

// VoteForTarget 为帖子或评论投票
func VoteForTarget(userID int64, p *models.ParamVoteData) (voteNum int64, err error) {
	zap.L().Debug("VoteForTarget",
		zap.Int64("userID", userID),
		zap.Int64("targetID", p.TargetID),
		zap.Int8("targetType", p.TargetType),
		zap.Int8("direction", p.Direction))

	// 根据目标类型调用不同的投票函数
	switch p.TargetType {
	case TypePost:
		return redis.VoteForPost(
			strconv.FormatInt(userID, 10),
			strconv.FormatInt(p.TargetID, 10),
			float64(p.Direction))
	case TypeComment:
		return redis.VoteForComment(
			strconv.FormatInt(userID, 10),
			strconv.FormatInt(p.TargetID, 10),
			float64(p.Direction))
	default:
		return 0, errors.New("无效的投票目标类型")
	}
}
