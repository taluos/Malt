package storage

import (
	"time"

	"github.com/taluos/Malt/pkg/storage/models"
)

type QueryOptions struct {
	Method    string
	HasError  *bool
	StartTime *time.Time
	EndTime   *time.Time
	Limit     int
	Offset    int
	OrderBy   string
}

func (s *SQLiteStorage) QueryRpcCallRecords(opts QueryOptions) ([]models.RpcCallRecord, error) {
	db := s.db.Model(&models.RpcCallRecord{})

	if opts.Method != "" {
		// db = db.Where("method = ?", opts.Method)
		db = db.Model(models.RpcCallRecord{Method: opts.Method})
	}
	if opts.HasError != nil {
		if *opts.HasError {
			db = db.Where("error != ''")
		} else {
			db = db.Where("error = ''")
		}
	}
	if opts.StartTime != nil {
		db = db.Where("timestamp >= ?", *opts.StartTime)
	}
	if opts.EndTime != nil {
		db = db.Where("timestamp <= ?", *opts.EndTime)
	}
	if opts.OrderBy != "" {
		db = db.Order(opts.OrderBy)
	} else {
		db = db.Order("timestamp desc")
	}
	if opts.Limit > 0 {
		db = db.Limit(opts.Limit)
	}
	if opts.Offset > 0 {
		db = db.Offset(opts.Offset)
	}

	var records []models.RpcCallRecord
	err := db.Find(&records).Error
	return records, err
}
