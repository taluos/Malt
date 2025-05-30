# Malt

Malt 是一个轻量的 Go 微服务框架，旨在简化分布式系统的开发和部署。框架提供了丰富的功能组件，包括 HTTP 服务器、RPC 客户端/服务器、指标监控、链路追踪等，帮助开发者快速构建高性能、可观测的微服务应用。

## 特性

- HTTP 服务 ：基于 Gin/Fiber 的 RESTful API 服务，支持中间件、参数验证、JWT 认证等
- RPC 通信 ：基于 gRPC 的高性能 RPC 服务，支持服务发现、负载均衡等
- 可观测性 ：集成 Prometheus 指标监控、OpenTelemetry 链路追踪
- 性能分析 ：内置 pprof 性能分析工具，便于问题排查
- 权限控制 ：支持 RBAC 基于角色的访问控制
- 错误显示 ：支持多语言错误提示和消息翻译
- 选项模式 ：采用函数式选项模式，灵活配置各组件参数

## quick start

### 安装

```bash
go get github.com/taluos/Malt
```
