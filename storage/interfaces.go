package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type gormAdapterInterface interface {
	Transaction(db *gorm.DB, fc func(tx *gorm.DB) error) error
	DeleteUnscoped(db *gorm.DB, query interface{}, model interface{}) *gorm.DB
	Create(db *gorm.DB, value interface{}) *gorm.DB
	SelectFirst(db *gorm.DB, query interface{}, destination interface{}) *gorm.DB
	AutoMigrate(db *gorm.DB, dst ...interface{}) error
}

type redisAdapterInterface interface {
	Del(ctx context.Context, keys ...string) error
	SaveMultipleInPipe(ctx context.Context, values ...redisValue) ([]redis.Cmder, error)
	Get(ctx context.Context, key string) (string, error)
	GetScanIterator(ctx context.Context, cursor uint64, match string, count int64) redisIteratorInterface
}
