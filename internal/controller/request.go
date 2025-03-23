package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

const CtxUserIDKey = "userID"

var ErrorUserNotLogin = errors.New("用户未登录")

// getCurrentUserID 获取当前登录的用户ID
func getCurrentUserId(c *gin.Context) (userId int64, err error) {
	_userID, ok := c.Get(CtxUserIDKey)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	userId, ok = _userID.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

// getPageInfo 分页参数
func getPageInfo(c *gin.Context) (int64, int64) {
	// 获取分页参数
	pageNumStr := c.Query("page")
	pageSizeStr := c.Query("size")

	var (
		page int64 // 页数，第几页
		size int64 // 每页几条数据
		err  error
	)
	page, err = strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}
