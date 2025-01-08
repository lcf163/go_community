# Go Community

一个基于 Go 语言开发的现代化社区系统,提供用户注册、登录、发帖、评论等功能。

## 功能特性

- 用户系统
  - 用户注册/登录
  - JWT 认证
  - 参数校验
  - 多语言支持

- 社区管理
  - 社区创建
  - 社区列表
  - 社区详情

- 内容系统
  - 发布帖子
  - 帖子列表
  - 帖子详情
  - 帖子投票
  - 按时间/分数排序

- 其他特性
  - 雪花算法生成 ID
  - Redis 缓存
  - 统一错误处理
  - 日志记录
  - 接口文档

## 技术栈

### 后端
- 框架: [Gin](https://github.com/gin-gonic/gin)
- 数据库: MySQL + Redis
- ORM: [sqlx](https://github.com/jmoiron/sqlx)
- 配置: [Viper](https://github.com/spf13/viper)
- 日志: [Zap](https://github.com/uber-go/zap)
- 文档: [Swagger](https://github.com/swaggo/gin-swagger)
- 校验: [validator](https://github.com/go-playground/validator)
- JWT认证: [jwt-go](https://github.com/dgrijalva/jwt-go)

### 开发工具
- 热重载: [Air](https://github.com/cosmtrek/air)
- API测试: Postman
- 容器化: Docker

### 部署
- Web服务器: Nginx
- 容器编排: Docker Compose

## 项目结构

```bash
go-community/
├── config/ # 配置文件
├── controller/ # 控制器层
├── dao/ # 数据访问层
│ ├── mysql/ # MySQL操作
│ └── redis/ # Redis操作
├── logger/ # 日志模块
├── logic/ # 业务逻辑层
├── models/ # 数据模型
├── pkg/ # 公共工具包
├── routers/ # 路由配置
├── static/ # 静态资源
└── templates/ # 模板文件
```

## 快速开始

### 环境要求

- Go 1.23+
- MySQL 8.0+
- Redis 6.0+
- Docker & Docker Compose

### 安装

1. 克隆项目
```bash
git clone https://github.com/yourusername/go-community.git
cd go-community
```

2. 安装依赖
```bash
go mod download
```

3. 配置数据库
```bash
# 创建数据库和表结构
mysql -u root -p < models/create_tables.sql
```

4. 修改配置
```bash
# 复制配置文件模板
cp config.yaml.example config.yaml
# 修改配置文件
vim config.yaml
```

### 运行

1. 开发模式
```bash
# 使用 Air 热重载
air -c .air.conf
```

2. 生产模式
```bash
# 使用 Docker Compose
docker-compose up -d
```

## API 文档

启动服务后访问: http://localhost:8081/swagger/index.html

主要接口:
- POST /api/v1/signup - 用户注册
- POST /api/v1/login - 用户登录
- GET /api/v1/community - 社区列表
- POST /api/v1/post - 发布帖子
- GET /api/v1/posts - 帖子列表
- POST /api/v1/vote - 帖子投票

## 开发规范

1. 代码风格
- 遵循 Go 官方规范
- 使用 gofmt 格式化代码
- 添加必要的注释

2. Git提交
- 使用语义化的提交信息
- 每次提交保持功能单一

3. 错误处理
- 统一错误码
- 规范的错误返回
- 详细的日志记录

## 部署

1. 使用 Docker
```bash
# 构建镜像
docker build -t go-community .

# 运行容器
docker run -p 8081:8081 go-community
```

2. 使用 Docker Compose
```bash
docker-compose up -d
```

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交代码
4. 发起 Pull Request

## 许可证

[MIT License](LICENSE)

## 联系方式

- 作者: Your Name
- 邮箱: your.email@example.com
- 项目地址: https://github.com/yourusername/go-community