package controller

import (
	"errors"
	"fmt"
	"go-community/dao/mysql"
	"go-community/logic"
	"go-community/models"
	"go-community/pkg/jwt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignUpHandler 处理注册请求
func SignUpHandler(c *gin.Context) {
	// 1.获取请求参数和参数校验
	//var p models.ParamSignUp
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUpHandler with invalid param", zap.Error(err))
		// 判断 err 是否为 validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 请求参数错误，直接返回响应
			ResponseError(c, CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 2.业务逻辑处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler 处理注册请求的函数
func LoginHandler(c *gin.Context) {
	// 1.获取请求参数和参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("LoginHandler with invalid param", zap.Error(err))
		// 判断 err 是否为 validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// 请求参数错误，直接返回响应
			ResponseError(c, CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 2.业务逻辑处理
	user, err := logic.Login(p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.UserName), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	// 3.返回响应
	ResponseSuccess(c, gin.H{
		"userID":       fmt.Sprintf("%d", user.UserId), // JS 的数字小于 1<<53 -1，int64: 1<<63 -1
		"userName":     user.UserName,
		"accessToken":  user.AccessToken,
		"refreshToken": user.RefreshToken,
	})
}

// RefreshTokenHandler 刷新accessToken
func RefreshTokenHandler(c *gin.Context) {
	rt := c.Query("refresh_token")
	// 客户端携带 Token 有三种方式 1.放在请求头 2.放在请求体 3.放在 URI
	// 这里假设 Token 放在 Header 的 Authorization 中，并使用 Bearer 开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		ResponseErrorWithMsg(c, CodeInvalidToken, "请求头缺少Auth Token")
		c.Abort()
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ResponseErrorWithMsg(c, CodeInvalidToken, "Token格式不对")
		c.Abort()
		return
	}
	aToken, rToken, err := jwt.RefreshToken(parts[1], rt)
	zap.L().Error("jwt.RefreshToken failed", zap.Error(err))
	c.JSON(http.StatusOK, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}
