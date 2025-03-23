package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"go_community/global"
)

var db *sqlx.DB

// Init 初始化MySQL连接
func Init(cfg *global.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Asia%%2FShanghai",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName,
	)
	// 也可以使用MustConnect，连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}

func GetDB() *sqlx.DB {
	return db
}

// Close 关闭MySQL连接
func Close() {
	_ = db.Close()
}
