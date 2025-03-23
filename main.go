package main

import (
	"fmt"
	"go.uber.org/zap"
	"go_community/global"
	"go_community/internal/controller"
	"go_community/internal/dao/mysql"
	"go_community/internal/dao/redis"
	"go_community/internal/middlewares"
	"go_community/internal/routers"
	"go_community/pkg/snowflake"
)

// Go Web开发较通用的脚手架模板

// @title go_community backend
// @version 1.0
// @description go_community API documentation
// @termsOfService http://swagger.io/terms/

// @host localhost:8081
// @BasePath /api/v1

func main() {
	// 1. 加载配置
	if err := global.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	fmt.Println(global.Conf)
	fmt.Println(global.Conf.LogConfig == nil)
	// 2. 初始化日志
	if err := middlewares.Init(global.Conf.LogConfig, global.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync()
	zap.L().Debug("logger init success...")
	// 3. 初始化MySQL连接
	if err := mysql.Init(global.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer mysql.Close()
	// 4. 初始化Redis连接
	if err := redis.Init(global.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()
	// 雪花算法生成 ID
	if err := snowflake.Init(global.Conf.StartTime, global.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}
	// 初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed, err:%v\n", err)
		return
	}
	// 5. 注册路由
	r := routers.SetupRouter(global.Conf.Mode)
	err := r.Run(fmt.Sprintf(":%d", global.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}

}
