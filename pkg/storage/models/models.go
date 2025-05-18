package models

import (
	"time"

	"gorm.io/gorm"
)

// RpcCallRecord 定义一条 gRPC 调用的记录
type RpcCallRecord struct {
	ID     uint   `gorm:"primaryKey;autoIncrement"`
	Method string `gorm:"size:255;index"` // 方法名
	// Request   string         `gorm:"type:text"`      // 请求内容 (json序列化)
	// Response  string         `gorm:"type:text"`      // 响应内容 (json序列化)
	Duration  int64          // 耗时，单位：毫秒
	Error     string         `gorm:"type:text"` // 错误信息
	Timestamp time.Time      `gorm:"index"`     // 调用时间
	DeletedAt gorm.DeletedAt `gorm:"index"`     // 软删除支持
}
