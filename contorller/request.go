package controller

import (
	"errors"
	"go-community/middlewares"

	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")

// getCurrentUserID 获取当前登录的用户ID
func getCurrentUserID(c *gin.Context) (userID int64, err error) {
	_userID, ok := c.Get(middlewares.CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userID, ok = _userID.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}
