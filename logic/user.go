package logic

import (
	"go-community/dao/mysql"
	"go-community/models"
	"go-community/pkg/snowflake"
)

// 存放业务逻辑的代码

func SignUp(p *models.ParamSignUp) {
	// 判断用户是否存在
	mysql.QueryUserByUsername()
	// 生成 UID
	snowflake.GetID()
	// 保存进数据库
	mysql.InsertUser()
	// redis.xxx
}
