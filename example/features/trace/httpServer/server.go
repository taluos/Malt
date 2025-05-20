package main

import (
	"context"
	"log"
	"os"
	"time"

	maltAgent "github.com/taluos/Malt/core/trace"
	httpserver "github.com/taluos/Malt/server/rest/httpServer"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/opentelemetry/tracing"
)

type BaseModel struct {
	ID uint32 `gorm:"primaryKey"` // Standard field for the primary key

	CreatedAt time.Time `gorm:"column:add_time"`    // 创建时间（由GORM自动管理）
	UpdatedAt time.Time `gorm:"column:update_time"` // 最后一次更新时间（由GORM自动管理）
	DeletedAt gorm.DeletedAt
	IsDeleted bool `gorm:"column:is_deleted"`
}

type User struct {
	BaseModel BaseModel `gorm:"embedded"`

	NickName string     `gorm:"type:varchar(20);default:UserName"` // 一个常规字符串字段
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null"`
	Password string     `gorm:"type:varchar(100);not null"`
	Email    *string    `gorm:"index:idx_email;unique;not null"` // 一个指向字符串的指针, allowing for null values 使用指针可以赋予非零值
	Birthday *time.Time `gorm:"type:datetime"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6)"`
	Role     int        `gorm:"column:role;default:1;type:int"`
}

// 全局 TracerProvider
var globalAgent *maltAgent.Agent

func NewTracerProvider(name string) *maltAgent.Agent {

	agent := maltAgent.NewAgent(name, "http://localhost:4318", "ratio", 1.0, "collector",
		maltAgent.WithTracerProviderOptions(traceSDK.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
			attribute.String("env", "test"),
		))),
	)

	return agent
}

func main() {
	var err error
	ctx := context.Background()
	// 初始化全局 TracerProvider
	globalAgent = NewTracerProvider("HTTP Server")
	defer globalAgent.Shutdown(ctx)

	// 获取 tracer
	tr := maltAgent.NewTracer(trace.SpanKindServer,
		maltAgent.WithTracerProvider(globalAgent.TracerProvider()),
		maltAgent.WithTracerName("test server"))

	// 创建根 span
	spanCtx, rootSpan := tr.Start(ctx,
		"http server root span", globalAgent.Propagator(), nil)
	defer tr.End(ctx, rootSpan, err)

	r := httpserver.NewServer(
		httpserver.WithPort(8080),
		httpserver.WithEnableTracing(true),
		httpserver.WithMiddleware(gin.Recovery()),
	)

	r.GET("/", func(c *gin.Context) {})
	r.GET("/server", Server)
	r.GET("/gorm", Gorm)

	// 使用带有 span 上下文的 context 启动服务器
	r.Start(spanCtx)
}

func Server(c *gin.Context) {
	// 从请求中提取 span 上下文
	ctx := c.Request.Context()
	var err error

	// 使用全局 agent 而不是创建新的
	// 获取 tracer
	tr := maltAgent.NewTracer(trace.SpanKindServer,
		maltAgent.WithTracerProvider(globalAgent.TracerProvider()),
		maltAgent.WithTracerName("server-handler"))

	// 创建新的 span
	ctx, span := tr.Start(ctx, "server-handler", globalAgent.Propagator(), nil)
	defer tr.End(ctx, span, err)

	// 添加一些属性到 span
	span.SetAttributes(attribute.String("handler", "server"))

	time.Sleep(time.Second * 1)

	c.JSON(200, gin.H{})
}

func Gorm(c *gin.Context) {
	// 从请求中提取 span 上下文
	ctx := c.Request.Context()
	var err error

	// 使用全局 agent 而不是创建新的
	// 获取 tracer
	tr := maltAgent.NewTracer(trace.SpanKindServer,
		maltAgent.WithTracerProvider(globalAgent.TracerProvider()),
		maltAgent.WithTracerName("gorm-handler"))

	// 创建新的 span
	ctx, span := tr.Start(ctx, "gorm-handler", globalAgent.Propagator(), nil)
	defer tr.End(ctx, span, err)

	dns := "root:root@tcp(192.168.142.137:3306)/shop?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		span.RecordError(err)
		panic(err)
	}

	// 使用带有 span 上下文的 context
	if err := db.Use(tracing.NewPlugin(tracing.WithTracerProvider(globalAgent.TracerProvider()))); err != nil {
		span.RecordError(err)
		panic(err)
	}

	if err := db.WithContext(ctx).Model(&User{}).Where("id > 10").Find(&User{}).Error; err != nil {
		span.RecordError(err)
		panic(err)
	}

	span.AddEvent("查询完成")
	time.Sleep(time.Second * 1)

	c.JSON(200, gin.H{})
}
