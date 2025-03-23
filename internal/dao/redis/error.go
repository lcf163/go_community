package redis

import "errors"

var (
	ErrorVoteTimeExpire = errors.New("投票时间已过")
	ErrorVoteRepeted    = errors.New("不允许重复投票")
)
