package routers

import (
	controller "go_community/controller"
	_ "go_community/docs" // 千万不要忘了导入把你上一步生成的docs
	"go_community/logger"
	"go_community/middlewares"
	"net/http"

	"github.com/gin-contrib/cors"
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
	r.Use(cors.Default())                               // 默认允许所有跨域请求
	// 自定义跨域请求 CORS 相关配置项
	//r.Use(cors.New(cors.Config{
	//	AllowOrigins:     []string{"https://foo.com"},
	//	AllowMethods:     []string{"PUT", "PATCH"},
	//	AllowHeaders:     []string{"Origin"},
	//	ExposeHeaders:    []string{"Content-Length"},
	//	AllowCredentials: true,
	//	AllowOriginFunc: func(origin string) bool {
	//		return origin == "https://github.com"
	//	},
	//	MaxAge: 12 * time.Hour,
	//}))

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
	// 帖子业务
	v1.GET("/posts", controller.GetPostListHandler)   // 分页展示帖子列表
	v1.GET("/posts2", controller.GetPostListHandler2) // 分页展示帖子列表：帖子的发布时间或分数排序
	v1.GET("/post/:id", controller.PostDetailHandler) // 查询帖子详情
	v1.GET("/search", controller.PostSearchHandler)   // 搜索帖子
	// 社区业务
	v1.GET("/community", controller.CommunityHandler)           // 获取分类社区列表
	v1.GET("/community2", controller.CommunityHandler2)         // 获取分类社区列表
	v1.GET("/community/:id", controller.CommunityDetailHandler) // 根据ID查找社区详情

	// 使用 JWT 认证中间件
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.POST("/post", controller.CreatePostHandler) // 创建帖子
		v1.PUT("/post", controller.UpdatePostHandler)  // 更新帖子
		v1.POST("/vote", controller.VoteHandler)       // 投票

		// 评论业务
		v1.POST("/comment", controller.CreateCommentHandler)         // 创建评论
		v1.GET("/comment/:postId", controller.GetCommentListHandler) // 获取评论列表
	}

	pprof.Register(r) // 注册 pprof 相关路由

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})

	return r
}
