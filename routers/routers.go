package routers

import (
	controller "go-community/contorller"
	"go-community/logger"
	"go-community/middlewares"
	"net/http"

	_ "go-community/docs" // 千万不要忘了导入把你上一步生成的docs

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	gs "github.com/swaggo/gin-swagger"
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
	// r.Use(logger.GinLogger(), logger.GinRecovery(true), // Recovery 中间件：recover 项目可能出现的 panic，并使用 zap 记录相关日志
	// 	middlewares.RateLimitMiddleware(2*time.Second, 1), // 限流中间件（全局限流）：每两秒钟添加1个令牌
	// )
	r.Use(logger.GinLogger(), logger.GinRecovery(true)) // Recovery 中间件：recover 项目可能出现的 panic，并使用 zap 记录相关日志

	// 前端相关
	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// 注册 swagger api 相关路由
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	// 注册路由
	v1 := r.Group("/api/v1")
	// 登录业务
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.GET("/refresh_token", controller.RefreshTokenHandler)

	v1.GET("/posts", controller.GetPostListHandler)             // 分页展示帖子列表
	v1.GET("/posts2", controller.GetPostListHandler2)           // 分页展示帖子列表：帖子的发布时间或分数排序
	v1.GET("/community", controller.CommunityHandler)           // 获取分类社区列表
	v1.GET("/community/:id", controller.CommunityDetailHandler) // 根据ID查找社区详情
	v1.GET("/post/:id", controller.PostDetailHandler)           // 查询帖子详情

	// 使用 JWT 认证中间件
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.POST("/post", controller.CreatePostHandler) // 创建帖子
		v1.POST("/vote", controller.VoteHandler)       // 投票
	}

	pprof.Register(r) // 注册 pprof 相关路由

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
