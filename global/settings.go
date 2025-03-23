package global

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局变量，用来保存程序的所有配置信息
var Conf = new(AppConfig)

// 定义环境常量
const (
	ModeDev  = "dev"
	ModeProd = "prod"
)

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"` // dev/prod
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
	Avatar       AvatarConfig  `mapstructure:"avatar"`
	Swagger      SwaggerConfig `mapstructure:"swagger"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// AvatarConfig 头像配置
type AvatarConfig struct {
	BaseURL string `mapstructure:"base_url"`
	MaxSize int64  `mapstructure:"max_size"`
	Domain  struct {
		Dev  string `mapstructure:"dev"`
		Prod string `mapstructure:"prod"`
	} `mapstructure:"domain"`
}

type SwaggerConfig struct {
	Domain struct {
		Dev  string `mapstructure:"dev"`
		Prod string `mapstructure:"prod"`
	} `mapstructure:"domain"`
}

// IsDevMode 判断是否为开发环境
func (c *AppConfig) IsDevMode() bool {
	return c.Mode == ModeDev
}

// IsProdMode 判断是否为生产环境
func (c *AppConfig) IsProdMode() bool {
	return c.Mode == ModeProd
}

// GetDomain 根据运行模式获取对应的域名
func (a *AvatarConfig) GetDomain() string {
	if Conf.IsProdMode() {
		return a.Domain.Prod
	}
	return a.Domain.Dev
}

// GetSwaggerHost 根据运行模式获取 Swagger Host
func (c *AppConfig) GetSwaggerHost() string {
	if c.IsProdMode() {
		return c.Swagger.Domain.Prod
	}
	return c.Swagger.Domain.Dev
}

// Init 初始化读取配置文件
func Init() (err error) {
	// 方式1：直接指定配置文件路径（相对路径或者绝对路径）
	// 相对路径：执行的可执行文件的相对路径
	// 绝对路径：系统中实际的文件路径
	viper.SetConfigFile("./configs/config.yaml")

	// 方式2：指定配置文件名和位置，viper自行查找可用的配置文件
	// 配置文件名不需要带后缀
	// 配置文件位置可以配置多个
	//viper.SetConfigName("config") // 指定配置文件名称（不需要带后缀）
	//viper.SetConfigType("yaml")   // 指定配置文件类型(专用于从远程获取配置信息时指定配置文件类型的)
	//viper.AddConfigPath(".") // 指定查找配置文件的路径（这里使用相对路径）
	//viper.AddConfigPath(".conf") // 指定查找配置文件的路径（这里使用相对路径）

	// 读取配置信息
	err = viper.ReadInConfig()
	if err != nil {
		// 读取配置信息失败
		panic(fmt.Errorf("ReadInConfig failed, err:%v\n", err))
	}
	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("Unmarshal failed, err:%v\n", err))
	}

	viper.WatchConfig()
	// 监听配置文件变化
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(Conf); err != nil {
			panic(fmt.Errorf("Unmarshal failed, err:%v\n", err))
		}
	})
	return err
}
