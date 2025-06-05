# Go Community

一个基于 Go 语言开发的社区系统，提供用户注册、登录、发帖、评论、投票等功能。

## 系统架构

### 整体架构
- 采用前后端分离架构
- RESTful API 设计
- 分层架构：controller -> logic -> dao
- 缓存设计：Redis + MySQL

### 目录结构
```bash
go_community/
├── config/      # 配置文件
├── docs/        # 文档
├── global/      # 全局变量
├── internal/    # 内部模块
│   ├── controller/  # 控制器层
│   ├── dao/         # 数据访问层
│   │  ├── mysql/      # MySQL操作
│   │  └── redis/      # Redis操作
│   ├── middlewares/ # 中间件
│   ├── models/      # 数据模型
│   ├── routers/     # 路由配置
│   └── service/     # 业务逻辑层
├── pkg/         # 模块包
├── scripts/     # 脚本文件
└── storage/     # 生成的临时文件
```

## 功能

### 用户功能
- 用户认证
  - 注册账号 POST `/api/v1/signup`
  - 用户登录 POST `/api/v1/login`
  - 刷新Token GET `/api/v1/refresh_token`
- 用户信息
  - 获取用户信息 GET `/api/v1/user/:id`
  - 修改用户名 PUT `/api/v1/user/name`
  - 修改密码 PUT `/api/v1/user/password`
  - 更新头像 POST `/api/v1/user/avatar`

### 社区功能
- 社区管理
  - 创建社区 POST `/api/v1/community`
  - 更新社区 PUT `/api/v1/community/:id`
  - 删除社区 DELETE `/api/v1/community/:id`
- 社区查询
  - 获取社区列表 GET `/api/v1/community`
  - 分页获取社区 GET `/api/v1/community2`
  - 获取社区详情 GET `/api/v1/community/:id`

### 帖子功能
- 帖子管理
  - 发布帖子 POST `/api/v1/post`
  - 编辑帖子 PUT `/api/v1/post`
  - 删除帖子 DELETE `/api/v1/post/:id`
- 帖子查询
  - 获取帖子列表(分页) GET `/api/v1/posts`
  - 获取帖子列表(分页+排序) GET `/api/v1/posts2`
  - 获取用户帖子列表 GET `/api/v1/posts/user/:id`
  - 获取帖子详情 GET `/api/v1/post/:id`
  - 搜索帖子 GET `/api/v1/search`

### 评论功能
- 评论管理
  - 发表评论/回复 POST `/api/v1/comment`
  - 编辑评论 PUT `/api/v1/comment`
  - 删除评论 DELETE `/api/v1/comment/:id`
  - 删除评论及回复 DELETE `/api/v1/comments/:id`
- 评论查询
  - 获取评论列表 GET `/api/v1/comments`
  - 获取评论详情 GET `/api/v1/comment/:id`

### 投票功能
- 投票操作
  - 投票(帖子/评论) POST `/api/v1/vote`
  - 支持点赞和踩

### 其他特性
- 跨域支持 (CORS)
- API 文档 (Swagger)
- 性能分析 (pprof)
- 404 处理

### 缓存设计
- Redis 缓存
  - 帖子
  	- 时间/分数排序
  	- 投票数据
  - 评论
  	- 时间排序
  	- 投票数据
  - 社区
  	- id查询帖子集合
- 双写一致性
  - MySQL 持久化存储
  - Redis 实时计数

### 性能优化
- 接口优化
  - 分页查询
  - 批量获取
  - 并发控制
- 缓存优化
  - 热点数据缓存
  - 定时更新策略

## 技术栈

### 后端
- 框架: [Gin](https://github.com/gin-gonic/gin)
  - 路由管理
  - 中间件支持
  - 参数绑定
- 数据库:
  - MySQL: 持久化存储
  - Redis: 缓存、计数器
- ORM: [sqlx](https://github.com/jmoiron/sqlx)
  - 原生SQL支持
  - 性能优化
- 配置: [Viper](https://github.com/spf13/viper)
  - 多环境配置
  - 热重载
- 日志: [Zap](https://github.com/uber-go/zap)
  - 结构化日志
  - 性能优化
- 认证: [jwt-go](https://github.com/dgrijalva/jwt-go)
  - Token生成与验证
  - 过期处理
- 文档: [Swagger](https://github.com/swaggo/gin-swagger)
  - API文档自动生成
  - 在线调试

### 开发工具
- 热重载: [Air](https://github.com/cosmtrek/air)
  - 代码热更新
  - 自定义配置
- API测试:
  - Postman: 接口调试
  - Swagger: 文档测试
- 监控调试:
  - pprof: 性能分析

### 部署
- 反向代理: Nginx
  - 负载均衡
- 容器化: Docker
  - 多阶段构建
  - 镜像优化
- 编排: Docker Compose
  - 服务编排
  - 环境隔离

## 快速开始

### 环境要求
- Go 1.23+
- MySQL 8.0+
- Redis 5.0+
- Docker & Docker Compose

### 安装

1. 克隆项目
```bash
git clone https://github.com/xxx/go_community.git
cd go_community
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
# 修改配置文件
vim conf/config.yaml
```

### 运行

1. 开发模式
```bash
# 修改配置文件 conf/config.yaml
mode: "dev"

# 使用 Air 热重载
air -c .air.conf
```

2. 生产模式
```bash
# 修改配置文件 conf/config.yaml
mode: "prod"

# 使用 Docker
```

### 部署

1. 使用 Docker
```bash
# 构建镜像
docker build . -t go_community_app

# 拉取 MySQL8.0 镜像
docker pull mysql:8.0      
# 运行 MySQL 容器
docker run --name mysql8 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=your_password -d mysql:8.0

# 拉取 Redis 镜像
docker pull redis:5.0.7
# 运行 Redis 容器
docker run --name redis507 -p 6379:6379 -d redis:5.0.7 redis-server --appendonly yes

# 运行容器
docker run --link=mysql8:mysql8 --link=redis507:redis507 -p 8088:8081 go_community_app
```

2. 使用 Docker Compose
```bash
docker-compose up -d
```

## API 文档

启动服务后访问: http://localhost:8081/swagger/index.html

## 开发规范

1. 代码风格
- 遵循 Go 官方规范
- 使用 gofmt 格式化代码
- 添加必要的注释

2. 错误处理
- 统一错误码
- 规范的错误返回
- 详细的日志记录

3. Git提交
- 使用语义化的提交信息
- 每次提交保持功能单一

## 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交代码
4. 发起 Pull Request
