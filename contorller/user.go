package controller

import (
	"fmt"
	"go-community/logic"
	"go-community/models"
	"net/http"

	"github.com/go-playground/validator/v10"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// SignUpHandler 处理注册请求的函数
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
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"msg": removeTopStruct(errs.Translate(trans)), // 类型错误则进行翻译
			})
			return
		}
	}
	// 手动对请求参数进行详细的业务规则校验
	//if len(p.UserName) == 0 || len(p.Password) == 0 || len(p.ConfirmPassword) == 0 || p.ConfirmPassword != p.Password {
	//	// 请求参数错误，直接返回响应
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数错误",
	//	})
	//	return
	//}
	fmt.Println(p)
	// 2.业务处理
	logic.SignUp(p)
	// 3.返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}
