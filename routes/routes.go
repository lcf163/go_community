package routes

import (
	controller "go-community/contorller"
	"go-community/logger"
	"go-community/middlewares"
	"go-community/settings"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(mode string) *gin.Engine {
	// 设置成发布模式
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	// 初始化路由（没有默认中间件）
	r := gin.New()
	// 设置中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册路由
	v1 := r.Group("/api/v1")
	// 登录业务
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.GET("/refresh_token", controller.RefreshTokenHandler)

	// 使用 JWT 认证中间件
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
