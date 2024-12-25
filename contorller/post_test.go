// post 单元测试

package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	url := "api/v1/post"
	r.POST(url, CreatePostHandler)

	body := `{
		"community_id": 5,
		"title": "test",
		"content": "just a test"
	}`
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// 判断响应的内容，是否按预期返回了需要登录的错误
	// 1.方法一：判断响应的内容是否包含指定的字符串
	//assert.Equal(t, w.Body.String(), "未登录")

	// 2.方法二：将响应的内容反序列化到 ResponseData，判断字段与预期是否一致
	res := new(ResponseData)
	if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
		t.Fatalf("json.Unmarshal w.Body failed, err:%v\n", err)
	}
	assert.Equal(t, res.Code, CodeNotLogin)
}
