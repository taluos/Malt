# HTTP 服务器模块

## 简介

HTTP 服务器模块是 Malt 框架的核心组件之一，基于 Gin 框架实现，提供了 RESTful API 服务能力。该模块支持多种功能特性，包括但不限于：健康检查、性能分析、指标监控、链路追踪以及 JWT 认证等。

## 主要特性

- **选项模式配置**：采用函数式选项模式，灵活配置服务器参数
- **健康检查**：提供 `/health` 接口，用于服务健康状态监控
- **性能分析**：集成 pprof，便于性能问题排查
- **指标监控**：支持 Prometheus 指标收集
- **链路追踪**：集成 OpenTelemetry，支持分布式追踪
- **JWT 认证**：内置 JWT 认证机制，保障 API 安全
- **错误提示**：支持多语言错误提示
- **参数验证**：内置多种参数验证器，如邮箱、手机号、用户名和密码等

## 使用示例

```go
package main

import (
    "context"
    "Malt/server/rest/httpServer"
    "github.com/gin-gonic/gin"
)

func main() {
    // 创建 HTTP 服务器实例
    server := httpServer.NewServer(
        httpServer.WithName("api-server"),
        httpServer.WithPort(8080),
        httpServer.WithMode("release"),
        httpServer.WithHealthz(true),
        httpServer.WithEnableMetrics(true),
    )
    
    // 注册路由
    server.GET("/api/v1/users", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "获取用户列表成功"})
    })
    
    // 启动服务器
    ctx := context.Background()
    if err := server.Start(ctx); err != nil {
        panic(err)
    }
}
```
