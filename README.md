# Go Community

一个基于 Go 语言开发的社区系统，提供用户注册、登录、发帖、评论等功能。

## 功能特性

- 用户系统
  - 用户注册/登录
  - JWT 认证
  - 参数校验
  - 多语言支持

### 社区管理
- 社区创建与管理
  - 创建社区
  - 社区信息维护
- 社区列表展示
  - 分页查询
  - 社区详情
- 权限控制
  - 管理员权限
  - 用户权限

### 内容系统
- 帖子管理
  - 发布帖子
  - 编辑帖子
  - 删除帖子
- 帖子列表
  - 分页展示
  - 时间排序
  - 分数排序
- 投票系统
  - 帖子投票
  - 按时间/分数排序

### 性能优化
- 缓存设计
  - Redis 热点数据缓存
  - 定时更新策略
- 查询优化
  - 索引优化
  - 分页查询
- 并发控制
  - 限流中间件
  - 分布式锁

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
- 容器化: Docker
  - 多阶段构建
  - 镜像优化
- 编排: Docker Compose
  - 服务编排
  - 环境隔离
- 反向代理: Nginx
  - 负载均衡

## 项目亮点

### 部署
- Web服务器: Nginx
- 容器编排: Docker Compose

## 项目结构
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
3. **测试规范**
- 单元测试
- 集成测试
- 性能测试

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

2. 错误处理
- 统一错误码
- 规范的错误返回
- 详细的日志记录

3. Git提交
- 使用语义化的提交信息
- 每次提交保持功能单一

## 部署

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
- 项目地址: https://github.com/yourusername/go_community