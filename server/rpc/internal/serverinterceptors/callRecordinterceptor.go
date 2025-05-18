package serverinterceptors

import (
	"context"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/taluos/Malt/pkg/storage/models"

	"github.com/taluos/Malt/pkg/storage"

	"github.com/taluos/Malt/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func UnarySQLiteInterceptor(db *storage.SQLiteStorage) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		startTime := time.Now()

		var resp any
		var err error
		done := make(chan struct{})
		panicChan := make(chan any, 1)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- fmt.Sprintf("%+v\n\n%s", p, strings.TrimSpace(string(debug.Stack())))
				}
			}()
			resp, err = handler(ctx, req)
			close(done)
		}()

		select {
		case p := <-panicChan:
			duration := time.Since(startTime)
			saveRpcRecord(info.FullMethod, duration, db, fmt.Sprintf("panic: %v", p))
			panic(p)
		case <-done:
			duration := time.Since(startTime)
			saveRpcRecord(info.FullMethod, duration, db, err)
			return resp, err
		case <-ctx.Done():
			duration := time.Since(startTime)
			saveRpcRecord(info.FullMethod, duration, db, ctx.Err())
			return nil, status.Error(status.Code(ctx.Err()), ctx.Err().Error())
		}
	}
}

func saveRpcRecord(method string, duration time.Duration, db *storage.SQLiteStorage, err any) {

	errorStr := fmt.Sprintf("%v", err)

	record := models.RpcCallRecord{
		Method:    method,
		Duration:  duration.Milliseconds(),
		Error:     errorStr,
		Timestamp: time.Now(),
	}

	// 异步插入
	go func(rec models.RpcCallRecord) {
		if dbErr := db.Insert(&rec); dbErr != nil {
			// 这里可以加日志记录错误
			errors.New(fmt.Sprintf("save rpc call record error: %v\n", dbErr.Error()))
		}
	}(record)
}
