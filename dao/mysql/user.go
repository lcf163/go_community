package mysql

import (
	"crypto/md5"
	"encoding/hex"
)

// 把每一步数据库操作封装成函数
// 等待 logic 层根据业务需求调用

const secret = "liwenzhou.com"

func encryptPassword(data []byte) (result string) {
	h := md5.New()
	h.Write([]byte(secret))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func QueryUserByUsername() {
	// ...
}

func InsertUser() {
	// 执行 SQL 语句入库
}
