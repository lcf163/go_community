package controller

import "go-community/models"

// 专门用来放接口文档用到的 model
// 因为接口文档返回的数据格式是一致的，但是具体的 data 类型不一致

type _ResponsePostList struct {
	Code    MyCode                  `json:"code"`    // 业务响应状态码
	Message string                  `json:"message"` // 提示信息
	Data    []*models.ApiPostDetail `json:"data"`    // 数据
}
