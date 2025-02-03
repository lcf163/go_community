package controller

import "go_community/models"

// 接口文档用到的 model
// 因为接口文档返回的数据格式是一致的，但是具体的 data 类型不一致

// _ResponsePostList 帖子列表响应
type _ResponsePostList struct {
	Code    MyCode                  `json:"code" example:"1000"`       // 业务响应状态码
	Message string                  `json:"message" example:"success"` // 提示信息
	Data    []*models.ApiPostDetail `json:"data"`                      // 帖子列表数据
}

// _ResponsePostDetail 帖子详情响应
type _ResponsePostDetail struct {
	Code    MyCode                `json:"code" example:"1000"`       // 业务响应状态码
	Message string                `json:"message" example:"success"` // 提示信息
	Data    *models.ApiPostDetail `json:"data"`                      // 帖子详情数据
}

// _ResponseCommentList 评论列表响应
type _ResponseCommentList struct {
	Code    MyCode                    `json:"code" example:"1000"`       // 业务响应状态码
	Message string                    `json:"message" example:"success"` // 提示信息
	Data    *models.ApiCommentListRes `json:"data"`                      // 评论列表数据
}

// _ResponseCommentDetail 评论详情响应
type _ResponseCommentDetail struct {
	Code    MyCode                   `json:"code" example:"1000"`       // 业务响应状态码
	Message string                   `json:"message" example:"success"` // 提示信息
	Data    *models.ApiCommentDetail `json:"data"`                      // 评论详情数据
}
