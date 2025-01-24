package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"go_community/models"
	"go_community/pkg/file"
	"strings"
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
	sqlStr := `select count(user_id) from user where username = ? and status = 1`
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
	// 设置随机默认头像 - 只存储文件名
	user.Avatar = file.GetRandomDefaultAvatar()
	// 设置默认状态为1
	user.Status = 1
	// 执行 SQL 语句入库
	sqlStr := `insert into user(user_id, username, password, avatar, status) values (?,?,?,?,?)`
	_, err = db.Exec(sqlStr, user.UserId, user.UserName, user.Password, user.Avatar, user.Status)
	return
}

// Login 用户登录
func Login(user *models.User) (err error) {
	originPassword := user.Password // 用户登录的原始密码
	sqlStr := `select user_id, username, password from user where username = ? and status = 1`
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
	sqlStr := `select user_id, username, avatar from user where user_id = ? and status = 1`
	err = db.Get(user, sqlStr, id)
	if err == sql.ErrNoRows {
		return nil, ErrorUserNotExist
	}
	return
}

// UpdateUserAvatar 更新用户头像
func UpdateUserAvatar(userId int64, avatarPath string) error {
	sqlStr := `update user set avatar = ? where user_id = ? and status = 1`
	result, err := db.Exec(sqlStr, avatarPath, userId)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrorInvalidID
	}
	return nil
}

// UpdateUserName 更新用户名
func UpdateUserName(userId int64, p *models.ParamUpdateUser) error {
	// 构建更新语句
	var updates []string
	var args []interface{}

	if p.Username != "" {
		updates = append(updates, "username = ?")
		args = append(args, p.Username)
	}

	// 如果没有要更新的字段
	if len(updates) == 0 {
		return nil
	}

	// 构建SQL语句，添加 status = 1 检查
	sqlStr := fmt.Sprintf("update user set %s where user_id = ? and status = 1",
		strings.Join(updates, ", "))
	args = append(args, userId)

	// 执行更新
	result, err := db.Exec(sqlStr, args...)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrorInvalidID
	}
	return nil
}

// CheckPassword 检查密码是否正确
func CheckPassword(userId int64, password string) error {
	sqlStr := `select password from user where user_id = ? and status = 1`
	var hashedPassword string
	if err := db.Get(&hashedPassword, sqlStr, userId); err != nil {
		if err == sql.ErrNoRows {
			return ErrorUserNotExist
		}
		return err
	}

	// 验证密码
	if hashedPassword != encryptPassword([]byte(password)) {
		return ErrorPasswordWrong
	}
	return nil
}

// UpdatePassword 更新密码
func UpdatePassword(userId int64, newPassword string) error {
	sqlStr := `update user set password = ? where user_id = ? and status = 1`
	result, err := db.Exec(sqlStr, encryptPassword([]byte(newPassword)), userId)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrorInvalidID
	}
	return nil
}
