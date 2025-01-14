package logic

import (
	"go_community/dao/mysql"
	"go_community/models"
	"go_community/pkg/jwt"
	"go_community/pkg/snowflake"
)

// 存放业务逻辑的代码

// SignUp 注册业务逻辑
func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	if err := mysql.CheckUserExist(p.UserName); err != nil {
		return err
	}
	// 生成 UID
	userId := snowflake.GetID()
	// 构造一个 User 实例
	user := &models.User{
		UserId:   userId,
		UserName: p.UserName,
		Password: p.Password,
	}
	// 保存进数据库
	return mysql.InsertUser(user)
	// redis.xxx
}

// Login 登录业务逻辑
func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		UserName: p.UserName,
		Password: p.Password,
	}
	// 用户登录，传递的是指针
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成 JWT token
	//return jwt.GenToken(user.UserID)
	accessToken, refreshToken, err := jwt.GenToken(user.UserId)
	if err != nil {
		return
	}
	user.AccessToken = accessToken
	user.RefreshToken = refreshToken
	return
}

// GetUserInfo 获取用户信息
func GetUserInfo(userId int64) (*models.User, error) {
	// 从数据库中查询用户信息
	user, err := mysql.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, mysql.ErrorUserNotExist
	}
	return user, nil
}
