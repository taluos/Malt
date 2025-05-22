# REST 服务器内部组件

## 简介

本目录包含 REST 服务器的内部组件，这些组件为 REST 服务器提供了各种功能支持，包括但不限于：参数验证、性能分析、中间件、RBAC 权限控制等。这些组件不建议直接在应用代码中使用，而是通过 REST 服务器的接口进行调用。

## 目录结构

- **middlewares**：中间件组件，包括 CORS、JWT、Recovery 等
- **pprof**：性能分析组件，提供 CPU、内存、协程等性能数据
- **RBAC**：基于角色的访问控制组件，提供权限管理功能
- **validations**：参数验证组件，提供各种数据格式的验证功能

## 主要组件

### 中间件 (middlewares)

中间件组件提供了一系列 Gin 中间件，用于处理请求前后的逻辑：

- **CORS**：跨域资源共享中间件，允许跨域请求
- **JWT**：JWT 认证中间件，验证请求的 JWT Token
- **Recovery**：恢复中间件，捕获并处理 panic
- **Logger**：日志中间件，记录请求日志
- **Timeout**：超时中间件，控制请求处理时间

### 性能分析 (pprof)

性能分析组件集成了 Go 的 pprof 工具，提供以下功能：

- **CPU 分析**：分析 CPU 使用情况
- **内存分析**：分析内存分配情况
- **协程分析**：分析协程运行情况
- **阻塞分析**：分析协程阻塞情况

### RBAC 权限控制

RBAC 组件基于 Casbin 实现，提供基于角色的访问控制功能：

- **角色管理**：创建、删除、更新角色
- **权限管理**：分配、撤销权限
- **权限验证**：验证用户是否有权限访问资源

### 参数验证 (validations)

参数验证组件提供了一系列验证器，用于验证请求参数的格式：

- **Email**：验证邮箱格式
- **Mobile**：验证手机号格式
- **Password**：验证密码强度
- **Username**：验证用户名格式

## 使用示例

这些组件主要通过 REST 服务器的接口进行调用，不建议直接在应用代码中使用。以下是一些通过 REST 服务器使用这些组件的示例：

### 使用中间件

```go
server := httpserver.NewServer(
    httpserver.WithMiddleware(middlewares.CORS()),
    httpserver.WithMiddleware(middlewares.JWT(jwtConfig)),
)
```

### 使用性能分析

```go
server := httpserver.NewServer(
    httpserver.WithEnableProfiling(true),
)
```
