package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"go_community/models"
	"go_community/pkg/file"
)

// 把每一步数据库操作封装成函数
// 等待 logic 层根据业务需求调用

const secret = "liwenzhou.com"

// encryptPassword 对密码进行加密
func encryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 生成加密密码
	user.Password = encryptPassword([]byte(user.Password))
	// 设置随机默认头像
	user.Avatar = file.GetRandomDefaultAvatar()
	// 执行 SQL 语句入库
	sqlStr := `insert into user(user_id, username, password, avatar) values (?,?,?,?)`
	_, err = db.Exec(sqlStr, user.UserId, user.UserName, user.Password, user.Avatar)
	return
}

// Login 用户登录
func Login(user *models.User) (err error) {
	originPassword := user.Password // 用户登录的原始密码
	sqlStr := `select user_id, username, password from user where username = ?`
	err = db.Get(user, sqlStr, user.UserName)
	// 用户不存在
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	// 查询数据库出错
	if err != nil {
		return err
	}
	// 判断密码是否正确
	password := encryptPassword([]byte(originPassword))
	if user.Password != password {
		return ErrorPasswordWrong
	}
	return
}

// GetUserById 根据ID查询作者信息
func GetUserById(id int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id, username, avatar from user where user_id = ?`
	err = db.Get(user, sqlStr, id)
	return
}

// UpdateUserAvatar 更新用户头像
func UpdateUserAvatar(userId int64, avatarPath string) error {
	sqlStr := `update user set avatar = ? where user_id = ?`
	_, err := db.Exec(sqlStr, avatarPath, userId)
	return err
}
