// mysql/post 测试

package mysql

import (
	"go_community/global"
	"go_community/internal/models"
	"testing"
	"time"
)

func init() {
	dbCfg := global.MySQLConfig{
		Host:         "127.0.0.1",
		User:         "root",
		Password:     "root",
		DbName:       "go_community",
		Port:         3307,
		MaxOpenConns: 10,
		MaxIdleConns: 10,
	}
	err := Init(&dbCfg)
	if err != nil {
		panic(err)
	}
}

func TestCreatePost(t *testing.T) {
	post := models.Post{
		PostID:      20,
		AuthorID:    1,
		CommunityID: 1,
		Status:      0,
		Title:       "test",
		Content:     "just a test",
		CreateTime:  time.Time{},
	}
	err := CreatePost(&post)
	if err != nil {
		t.Fatalf("CreatePost insert record into mysql failed, err:%v\n", err)
	}
	t.Logf("CreatePost insert record into mysql success")
}
