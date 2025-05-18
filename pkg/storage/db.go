package storage

import (
	"Malt/pkg/errors"

	"os"

	"github.com/glebarez/sqlite"

	// "gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"fmt"
	"path/filepath"
)

type SQLiteStorage struct {
	path string // 数据库文件路径
	db   *gorm.DB
}

func NewSQLiteStorage(dbPath string, db *gorm.DB) *SQLiteStorage {
	return &SQLiteStorage{path: dbPath, db: db}
}

// Init db and check the path and permission
// if user not provide the db, it will create a default db in the given path(dbPath)
func (s *SQLiteStorage) Init() error {
	// 这里可以加额外初始化逻辑，比如设置PRAGMA
	var err error

	// 获取绝对路径
	absPath, err := filepath.Abs(s.path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// 确保目录存在
	dir := filepath.Dir(absPath)
	if err = ensureDir(dir); err != nil {
		return errors.Wrapf(err, "failed to ensure directory")
	}

	// 检查读写权限
	if err = checkDirRW(dir); err != nil {
		return errors.Wrapf(err, "directory permission error")
	}

	// 检查db文件是否存在
	_, err = os.Stat(absPath)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrapf(err, "failed to stat db file")
	}

	// 打开SQLite数据库
	if s.db == nil {
		s.db, err = gorm.Open(sqlite.Open(absPath), &gorm.Config{})
		if err != nil {
			return errors.Wrapf(err, "failed to open sqlite db")
		}
	}

	return nil
}

func (s *SQLiteStorage) HasTable(table any) (bool, error) {
	if s.db == nil {
		return false, errors.New("database not initialized")
	}
	return s.db.Migrator().HasTable(table), nil
}

func (s *SQLiteStorage) CreateTable(table any) error {
	if s.db == nil {
		return errors.New("database not initialized")
	}
	return s.db.AutoMigrate(table)
}

func (s *SQLiteStorage) Insert(record any) error {
	if s.db == nil {
		return errors.New("database not initialized")
	}
	return s.db.Create(record).Error
}

func (s *SQLiteStorage) Delete(record any) error {
	if s.db == nil {
		return errors.New("database not initialized")
	}
	return s.db.Delete(record).Error
}

/*
	func (s *SQLiteStorage) Find(record any) error {
		if s.db == nil {
			return errors.New("database not initialized")
		}
		return s.db.First(record).Error
	}
*/

func (s *SQLiteStorage) Close() error {
	if s.db == nil {
		return errors.New("database not initialized")
	}
	sqlDB, err := s.db.DB()
	if err != nil {
		return errors.New("failed to close db")
	}
	return sqlDB.Close()
}

// ensureDir 确保目录存在
func ensureDir(dir string) error {
	info, err := os.Stat(dir)
	if err == nil {
		if !info.IsDir() {
			return errors.New("path exists but is not a directory")
		}
		return nil
	}
	if os.IsNotExist(err) {
		// 不存在就创建
		return os.MkdirAll(dir, 0755)
	}
	return err
}

// checkDirRW 检查目录的读写权限
func checkDirRW(dir string) error {
	// 尝试创建一个临时文件来验证权限
	tempFile := filepath.Join(dir, ".permission_check")
	f, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	f.Close()
	// 删除临时文件
	return os.Remove(tempFile)
}
