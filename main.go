package main

import (
	"fmt"
	controller "go_community/controller"
	"go_community/dao/mysql"
	"go_community/dao/redis"
	"go_community/logger"
	"go_community/pkg/snowflake"
	"go_community/routers"
	"go_community/settings"

	"go.uber.org/zap"
)

// Go Web开发较通用的脚手架模板

// @title go_community backend
// @version 1.0
// @description go_community API documentation
// @termsOfService http://swagger.io/terms/

// @host localhost:8081
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description 在值的前面添加 Bearer，并带上空格。例如：Bearer your_token_here

func main() {
	// 1. 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	fmt.Println(settings.Conf)
	fmt.Println(settings.Conf.LogConfig == nil)
	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")
	// 3. 初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()
	// 4. 初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()
	// 雪花算法生成 ID
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed, err:%v\n", err)
		return
	}
	// 5. 注册路由
	r := routers.SetupRouter(settings.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", settings.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}

}
