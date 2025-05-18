package main

import (
	agent "Malt/core/trace"
	httpserver "Malt/server/rest/httpServer"
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	traceSDK "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
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

func NewTracerProvider() *traceSDK.TracerProvider {

	agentOpt := agent.NewAgent("test http server", "http://localhost:4318", "ratio", 1.0, "collector",
		agent.WithTracerProviderOptions(traceSDK.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("test http server"),
			attribute.String("env", "test"),
		))),
	)

	tp := agent.InitAgent(agentOpt)
	return tp
}

func main() {

	_ = NewTracerProvider()

	r := httpserver.NewServer(
		httpserver.WithPort(8080),
		httpserver.WithEnableTracing(true),
		httpserver.WithMiddleware(gin.Recovery()),
	)

	r.GET("/", func(c *gin.Context) {})
	r.GET("/server", Server)
	r.GET("/gorm", Gorm)

	r.Start(context.Background())
}

func Server(c *gin.Context) {

	tp := NewTracerProvider()
	defer tp.Shutdown(context.Background())

	time.Sleep(time.Second * 1)

	c.JSON(200, gin.H{})
}

func Gorm(c *gin.Context) {

	tp := NewTracerProvider()
	defer tp.Shutdown(context.Background())

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
		panic(err)
	}

	if err := db.Use(tracing.NewPlugin(tracing.WithTracerProvider(tp))); err != nil {
		panic(err)
	}

	if err := db.WithContext(c.Request.Context()).Model(&User{}).Where("id > 10").Find(&User{}).Error; err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 1)

	c.JSON(200, gin.H{})
}
