package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"go-community/models"
)

// 把每一步数据库操作封装成函数
// 等待 logic 层根据业务需求调用

const secret = "liwenzhou.com"

var (
	ErrorUserExist       = errors.New("用户已存在")
	ErrorUserNotExist    = errors.New("用户不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
)

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
	password := encryptPassword([]byte(user.Password))
	// 执行 SQL 语句入库
	sqlStr := `insert into user(user_id, username, password) values (?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.UserName, password)
	return
}

// encryptPassword 密码加密
func encryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func Login(user *models.User) (err error) {
	originPassword := user.Password // 用户登录的原始密码
	sqlStr := `select user_id, username, password from user where username = ?`
	err = db.Get(user, sqlStr, user.UserName)
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	if err != nil {
		// 查询数据库出错
		return err
	}
	// 判断密码是否正确
	password := encryptPassword([]byte(originPassword))
	if user.Password != password {
		return ErrorInvalidPassword
	}
	return
}
