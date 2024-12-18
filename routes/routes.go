package routes

import (
	controller "go-community/contorller"
	"go-community/logger"
	"go-community/middlewares"
	"go-community/settings"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务路由
	r.POST("/signup", controller.SignUpHandler)
	r.POST("/login", controller.LoginHandler)

	r.GET("/version", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		// 如果是已登录的用户，判断请求头中是否包括有效的 JWT token
		c.String(http.StatusOK, settings.Conf.Version)
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
