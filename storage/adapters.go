package storage

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

type redisValue struct {
	key        string
	value      interface{}
	expiration time.Duration
}

type redisIteratorInterface interface {
	Err() error
	Next(ctx context.Context) bool
	Val() string
}

type redisAdapter struct {
	con *redis.Client
}

func (a *redisAdapter) Del(ctx context.Context, keys ...string) error {
	return a.con.Del(ctx, keys...).Err()
}
func (a *redisAdapter) SaveMultipleInPipe(ctx context.Context, values ...redisValue) ([]redis.Cmder, error) {
	pipe := a.con.TxPipeline()
	for _, data := range values {
		pipe.Set(ctx, data.key, data.value, data.expiration)
	}
	return pipe.Exec(ctx)
}
func (a *redisAdapter) Get(ctx context.Context, key string) (string, error) {
	return a.con.Get(ctx, key).Result()
}
func (a *redisAdapter) GetScanIterator(ctx context.Context, cursor uint64, match string, count int64) redisIteratorInterface {
	return a.con.Scan(ctx, cursor, match, count).Iterator()
}

type gormAdapter struct{}

func (a *gormAdapter) Transaction(db *gorm.DB, fc func(tx *gorm.DB) error) error {
	return db.Transaction(fc)
}
func (a *gormAdapter) DeleteUnscoped(db *gorm.DB, query interface{}, model interface{}) *gorm.DB {
	return db.Unscoped().Where(query).Delete(model)
}
func (a *gormAdapter) Create(db *gorm.DB, value interface{}) *gorm.DB {
	return db.Create(value)
}
func (a *gormAdapter) SelectFirst(db *gorm.DB, query interface{}, destination interface{}) *gorm.DB {
	return db.Where(query).First(destination)
}
func (a *gormAdapter) AutoMigrate(db *gorm.DB, dst ...interface{}) error {
	return db.AutoMigrate(dst...)
}
