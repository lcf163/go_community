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
	// 设置中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	v1 := r.Group("/api/v1")
	// 注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.GET("/refresh_token", controller.RefreshTokenHandler)

	r.GET("/version", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		// 如果是已登录的用户，判断请求头中是否包括有效的 JWT token
		c.String(http.StatusOK, settings.Conf.Version)
		//c.String(http.StatusOK, "xxx")
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
