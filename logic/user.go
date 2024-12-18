package logic

import (
	"go-community/dao/mysql"
	"go-community/models"
	"go-community/pkg/jwt"
	"go-community/pkg/snowflake"
)

// 存放业务逻辑的代码

func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	if err := mysql.CheckUserExist(p.UserName); err != nil {
		return err
	}
	// 生成 UID
	userID := snowflake.GetID()
	// 构造一个 User 实例
	user := &models.User{
		UserID:   userID,
		UserName: p.UserName,
		Password: p.Password,
	}
	// 保存进数据库
	return mysql.InsertUser(user)
	// redis.xxx
}

func Login(p *models.ParamLogin) (aToken, rToken string, err error) {
	user := &models.User{
		UserName: p.UserName,
		Password: p.Password,
	}
	// 用户登录，传递的是指针
	if err := mysql.Login(user); err != nil {
		return "", "", err
	}
	// 生成 JWT token
	return jwt.GenToken(user.UserID)
}
