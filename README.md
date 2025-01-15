# Go Community

一个基于 Go 语言开发的社区系统，提供用户注册、登录、发帖、评论、投票等功能。

## 系统架构

### 整体架构
- 采用前后端分离架构
- RESTful API 设计
- 分层架构：controller -> logic -> dao
- 缓存设计：Redis + MySQL 双写一致性

### 目录结构
```bash
go_community/
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

## 功能

### 用户功能
- 用户注册/登录
  - JWT 认证
  - Token 刷新
  - 参数校验
- 用户信息
  - 获取用户信息
  - 作者信息展示（待开发）
  - 修改用户信息（待开发）

### 社区功能
- 社区列表展示
  - 获取社区列表
  - 获取社区详情
  - 分页查询
- 权限控制
  - 用户权限校验
  - 操作权限管理

### 帖子功能
- 帖子管理
  - 发布帖子
  - 编辑帖子
  - 删除帖子（待开发）
  - 帖子详情查询
- 帖子列表
  - 分页展示
  - 按时间排序
  - 按分数排序
  - 社区内帖子查询
  - 关键词搜索
  - 点赞数统计
- 投票功能
  - 帖子投票
  - 分数权重算法

### 评论功能
- 评论管理
  - 发表评论
  - 编辑评论（待开发）
  - 刪除评论（待开发）
  - 评论列表查询
  - 评论回复功能
  - 评论投票
- 评论展示
  - 按时间排序
  - 点赞数统计
  - 评论数统计

### 缓存设计
- Redis 缓存
  - 帖子时间/分数排序
  - 投票记录存储
  - 社区帖子集合
  - 评论时间/投票数据
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
- 文档: [Swagger](https://github.com/swaggo/gin-swagger)
  - API文档自动生成
  - 在线调试
- 认证: [jwt-go](https://github.com/dgrijalva/jwt-go)
  - Token生成与验证
  - 过期处理

### 开发工具
- 热重载: [Air](https://github.com/cosmtrek/air)
  - 代码热更新
  - 自定义配置
- API测试:
  - Postman: 接口调试
  - Swagger: 文档测试
- 监控调试:
  - pprof: 性能分析
  - Prometheus: 指标收集

### 部署
- 反向代理: Nginx
  - 负载均衡
- 容器化: Docker
  - 多阶段构建
  - 镜像优化
- 编排: Docker Compose
  - 服务编排
  - 环境隔离

## 项目亮点

## 项目结构

## 快速开始

### 环境要求
- Go 1.23+
- MySQL 8.0+
- Redis 6.0+
- Docker & Docker Compose

### 安装

1. 克隆项目
```bash
git clone https://github.com/yourusername/go_community.git
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
# 使用 Docker
```

### 部署

1. 使用 Docker
```bash
# 构建镜像
docker build -t go_community .

# 运行容器
docker run -p 8081:8081 go_community
```

2. 使用 Docker Compose
```bash
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
