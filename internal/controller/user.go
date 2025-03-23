package controller

import (
	"errors"
	"fmt"
	"go_community/internal/dao/mysql"
	"go_community/internal/models"
	"go_community/internal/service"
	pkg_file "go_community/pkg/file"
	"go_community/pkg/jwt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignUpHandler 处理注册请求
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param user body models.ParamSignUp true "注册信息"
// @Success 1000 {object} ResponseData
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1002 {object} ResponseData "用户名已存在"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /signup [post]
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
	if err := service.SignUp(p); err != nil {
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

// LoginHandler 处理登录请求
// @Summary 用户登录
// @Description 用户登录并获取token
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param user body models.ParamLogin true "登录信息"
// @Success 1000 {object} ResponseData{data=map[string]string{user_id=string,user_name=string,access_token=string,refresh_token=string}}
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1004 {object} ResponseData "用户名或密码错误"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /login [post]
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
	user, err := service.Login(p)
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
		"user_id":       fmt.Sprintf("%d", user.UserID), // JS 的数字小于 1<<53-1，int64: 1<<63-1
		"user_name":     user.UserName,
		"access_token":  user.AccessToken,
		"refresh_token": user.RefreshToken,
	})
}

// RefreshTokenHandler 刷新token
// @Summary 刷新访问令牌
// @Description 使用refresh_token刷新access_token
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "Bearer 用户令牌"
// @Param refresh_token query string true "刷新令牌"
// @Success 200 {object} map[string]string{access_token=string,refresh_token=string}
// @Failure 1006 {object} ResponseData "无效的Token"
// @Router /refresh_token [post]
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

// GetUserInfoHandler 获取用户信息
// @Summary 获取用户信息
// @Description 获取指定用户的详细信息
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param id path int true "用户ID"
// @Success 1000 {object} ResponseData{data=map[string]string{user_id=string,username=string,avatar=string}}
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1003 {object} ResponseData "用户不存在"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /user/{id} [get]
func GetUserInfoHandler(c *gin.Context) {
	// 获取用户ID参数
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 获取用户信息
	user, err := service.GetUserInfo(userID)
	if err != nil {
		zap.L().Error("logic.GetUserInfo failed",
			zap.Int64("user_id", userID),
			zap.Error(err))
		if err == mysql.ErrorUserNotExist {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, gin.H{
		"user_id":  fmt.Sprintf("%d", user.UserID),
		"username": user.UserName,
		"avatar":   user.GetAvatarURL(),
	})
}

// UpdateUserNameHandler 更新用户名
// @Summary 更新用户名
// @Description 修改当前登录用户的用户名
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param user body models.ParamUpdateUser true "用户名信息"
// @Success 1000 {object} ResponseData
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1002 {object} ResponseData "用户名已存在"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /user/name [put]
func UpdateUserNameHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	// 获取参数
	p := new(models.ParamUpdateUser)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdateUserNameHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 更新用户名
	if err := service.UpdateUserName(userID, p); err != nil {
		zap.L().Error("logic.UpdateUserName failed",
			zap.Int64("user_id", userID),
			zap.Error(err))
		if err == mysql.ErrorUserExist {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

// UpdatePasswordHandler 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param password body models.ParamUpdatePassword true "密码信息"
// @Success 1000 {object} ResponseData
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1004 {object} ResponseData "原密码错误"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /user/password [put]
func UpdatePasswordHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	// 获取参数
	p := new(models.ParamUpdatePassword)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("UpdatePasswordHandler with invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParams)
		return
	}

	// 修改密码
	if err := service.UpdatePassword(userID, p); err != nil {
		zap.L().Error("logic.UpdatePassword failed",
			zap.Int64("user_id", userID),
			zap.Error(err))
		if err == mysql.ErrorPasswordWrong {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, nil)
}

// UpdateAvatarHandler 更新用户头像
// @Summary 更新头像
// @Description 更新当前登录用户的头像
// @Tags 用户相关接口
// @Accept multipart/form-data
// @Produce application/json
// @Security Bearer
// @Param Authorization header string true "Bearer 用户令牌"
// @Param avatar formData file true "头像文件"
// @Success 1000 {object} ResponseData{data=map[string]string{avatar=string,message=string}}
// @Failure 1001 {object} ResponseData "参数错误"
// @Failure 1008 {object} ResponseData "未登录"
// @Failure 1005 {object} ResponseData "服务繁忙"
// @Router /user/avatar [post]
func UpdateAvatarHandler(c *gin.Context) {
	// 获取当前用户ID
	userID, err := getCurrentUserId(c)
	if err != nil {
		ResponseError(c, CodeNotLogin)
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		zap.L().Error("get form file failed", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParams, "请选择要上传的头像文件")
		return
	}

	// 更新头像
	filename, err := service.UpdateAvatar(userID, file)
	if err != nil {
		zap.L().Error("logic.UpdateAvatar failed",
			zap.Int64("user_id", userID),
			zap.Error(err))

		// 根据具体错误类型返回相应的错误信息
		switch err {
		case pkg_file.ErrorFileLimit:
			ResponseErrorWithMsg(c, CodeInvalidParams, "文件大小超出限制")
		case pkg_file.ErrorFileType:
			ResponseErrorWithMsg(c, CodeInvalidParams, "不支持的文件类型，请上传jpg/jpeg/png/gif格式的图片")
		case pkg_file.ErrorFileDirectory:
			ResponseErrorWithMsg(c, CodeServerBusy, "服务器存储错误")
		default:
			ResponseError(c, CodeServerBusy)
		}
		return
	}

	// 返回成功响应，包含完整的头像URL
	ResponseSuccess(c, gin.H{
		"avatar":  pkg_file.GetAvatarPath(filename),
		"message": "头像更新成功",
	})
}
