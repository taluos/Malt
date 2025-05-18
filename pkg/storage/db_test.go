package storage

import (
	"Malt/pkg/storage/models"
	"fmt"

	"os"
	"testing"
	"time"
)

const testDBFile = "./test_rpc_calls.db"

func cleanupTestDB() {
	_ = os.Remove(testDBFile)
}

func setupTestStorage(t *testing.T) *SQLiteStorage {
	cleanupTestDB()
	storage := NewSQLiteStorage(testDBFile, nil)
	return storage
}

func TestSQLiteStorage_Init(t *testing.T) {
	storage := setupTestStorage(t)

	defer storage.Close()

	if err := storage.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
}

func TestSQLiteStorage_CreateAndHasTable(t *testing.T) {
	var err error
	storage := setupTestStorage(t)
	defer storage.Close()

	if err = storage.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 确认一开始没有表
	exists, err := storage.HasTable(&models.RpcCallRecord{})
	if err != nil {
		t.Fatalf("HasTable failed: %v", err)
	}
	if exists {
		t.Fatalf("expected table not exists initially")
	}

	// 创建表
	err = storage.CreateTable(&models.RpcCallRecord{})
	if err != nil {
		t.Fatalf("CreateTable failed: %v", err)
	}

	// 确认表已经存在
	exists, err = storage.HasTable(&models.RpcCallRecord{})
	if err != nil {
		t.Fatalf("HasTable after create failed: %v", err)
	}
	if !exists {
		t.Fatalf("expected table to exist after creation")
	}
}

func TestSQLiteStorage_Insert(t *testing.T) {
	var err error
	storage := setupTestStorage(t)
	defer storage.Close()

	if err = storage.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// 先创建表
	if err = storage.CreateTable(&models.RpcCallRecord{}); err != nil {
		t.Fatalf("CreateTable failed: %v", err)
	}

	// 插入一条记录，字段对齐你的结构体
	record := models.RpcCallRecord{
		Method: "test.service/Method",
		// Request:   `{"input":"hello"}`,
		// Response:  `{"output":"world"}`,
		Duration:  120,
		Error:     "there is no error",
		Timestamp: time.Now(),
	}

	if err = storage.Insert(&record); err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// 查询验证
	var result models.RpcCallRecord
	err = storage.db.First(&result).Error

	var results []models.RpcCallRecord
	queryOptions := QueryOptions{
		Method: "test.service/Method",
	}
	results, err = storage.QueryRpcCallRecords(queryOptions)

	fmt.Println(result)
	fmt.Println(results)
	if err != nil {
		t.Fatalf("query inserted record failed: %v", err)
	}

	// 简单验证下字段
	if result.Method != record.Method {
		t.Fatalf("inserted record not match: got %+v", result)
	}
}
