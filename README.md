# Go Gin 企业级 Web 服务模板

一个基于 **Gin + slog + Viper** 的 Go Web 服务基础模板，提供企业级项目结构与最佳实践，适合用于构建稳定、可维护的 Go Web API 服务。

---

# 一、项目介绍

本项目提供一套 **Go Web 服务标准工程结构**，适用于企业级开发，包含以下特性：

* 清晰的分层架构
* 结构化日志（slog）
* 配置管理（Viper）
* Gin 中间件体系
* 统一 API 返回格式
* 优雅关闭（Graceful Shutdown）
* 可扩展的项目结构
* 基础的认证授权逻辑

核心技术栈：

* Go
* Gin
* slog
* Viper
* jwt

---

# 二、项目结构

```
go-gin-template
│
├── cmd
│   └── server
│       └── main.go          // 程序入口
│
├── configs
│   └── config.yaml          // 配置文件
│
├── internal
│   ├── handler              // HTTP接口
│   │   └── user.go
│   │
│   ├── service              // 业务逻辑
│   │   └── user.go
│   │
│   ├── repository           // 数据库访问
│   │   └── user.go
│   │
│   ├── model                // 数据模型
│   │   └── user.go
│   │
│   ├── middleware           // Gin中间件
│   │   ├── logger.go
│   │   └── recovery.go
│   │
│   ├── router               // 路由配置
│   │   └── router.go
│   │
│   └── config               // 配置加载
│       └── config.go
│
├── pkg
│   ├── logger               // slog日志封装
│   │   └── logger.go
│   │
│   └── response             // API统一返回
│       └── response.go
│
├── go.mod
└── README.md
```

---

# 三、架构设计

项目采用典型的 **三层架构设计**：

```
Handler -> Service -> Repository -> Database
```

职责划分：

| 层级         | 职责            |
| ---------- | ------------- |
| Handler    | 处理 HTTP 请求与响应 |
| Service    | 处理业务逻辑        |
| Repository | 数据库访问         |
| Model      | 数据结构定义        |

这样可以确保：

* 代码职责清晰
* 易于测试
* 易于扩展

---

# 四、配置管理

项目使用 **Viper** 进行配置管理。

配置文件：

```
configs/config.yaml
```

示例：

```yaml
server:
  port: 8080

database:
  dsn: root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local

log:
  level: info
```

配置加载代码：

```
internal/config/config.go
```

负责读取并解析配置文件。

---

# 五、日志系统

项目使用 Go 官方提供的 **slog** 作为日志系统。

特点：

* 结构化日志
* JSON输出
* 易于接入日志平台

示例日志：

```json
{
  "time": "2026-03-25T12:00:00Z",
  "level": "INFO",
  "msg": "server start",
  "port": 8080
}
```

日志初始化位置：

```
pkg/logger/logger.go
```

---

# 六、统一 API 返回格式

所有接口返回统一结构：

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

字段说明：

| 字段   | 说明   |
| ---- | ---- |
| code | 状态码  |
| msg  | 提示信息 |
| data | 返回数据 |

封装位置：

```
pkg/response/response.go
```

使用示例：

```go
response.Success(c, data)

response.Fail(c, "error message")
```

---

# 七、中间件

项目使用 Gin 中间件体系。

常见中间件：

| 中间件       | 作用      |
| --------- | ------- |
| Logger    | 请求日志    |
| Recovery  | Panic恢复 |
| JWTAuth   | 认证（可扩展） |
| RateLimit | 限流（可扩展） |

示例：

```
r.Use(middleware.Logger())
r.Use(gin.Recovery())
```

---

# 八、优雅关闭

项目实现了 **Graceful Shutdown**。

当收到以下信号时：

* SIGINT
* SIGTERM

服务器会：

1. 停止接收新请求
2. 等待已有请求处理完成
3. 正常退出

实现位置：

```
cmd/server/main.go
```

---

# 九、启动项目

### 1 初始化模块

```
go mod init go-gin-app
```

### 2 安装依赖

```
go get github.com/gin-gonic/gin
go get github.com/spf13/viper
```

### 3 运行项目

```
go run cmd/server/main.go
```

---

# 十、测试接口

示例接口：

```
GET /user/:id
```

请求示例：

```
http://localhost:8080/user/1
```

返回示例：

```json
{
  "code": 0,
  "msg": "success",
  "data": {
    "id": 1,
    "name": "han"
  }
}
```

---

# 十一、开发规范

### 1 Handler 层

只负责：

* 解析请求
* 调用 Service
* 返回响应

不要写业务逻辑。

---

### 2 Service 层

负责：

* 业务逻辑
* 事务控制
* 调用 Repository

---

### 3 Repository 层

只负责：

* 数据库操作
* SQL查询

---

# 十二、推荐扩展

在此模板基础上，可以继续扩展：

### 1 用户认证

JWT 登录认证

### 2 数据库

MySQL + GORM

### 3 缓存

Redis

### 4 限流

golang.org/x/time/rate

### 5 监控

Prometheus + Grafana

### 6 容器部署

Docker + Kubernetes

---

# 十三、最佳实践

开发 Go Web 服务时建议遵循以下原则：

* 分层架构
* 统一错误处理
* 统一响应格式
* 结构化日志
* 优雅关闭
* 配置管理
* 中间件统一管理

---

# 十四、适用场景

本模板适用于：

* REST API 服务
* 微服务
* 企业后台服务
* SaaS 平台后端

---

# 十五、License

MIT License