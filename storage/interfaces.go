package storage

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

type gormInterface interface {
	Transaction(fc func(tx interface{ gormInterface }) error, opts ...*sql.TxOptions) (err error)
	Unscoped() (tx interface{ gormInterface })
	Where(query interface{}, args ...interface{}) (tx interface{ gormInterface })
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
	Create(value interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	AutoMigrate(dst ...interface{}) error
}

type redisInterface interface {
	Del(ctx context.Context, keys ...string) redisIntCmdInterface
	Get(ctx context.Context, key string) redisStringCmdInterface
	Scan(ctx context.Context, cursor uint64, match string, count int64) redisScanCmdInterface
	TxPipeline() redisPipelineInterface
}

type redisPipelineInterface interface {
	Exec(ctx context.Context) ([]redis.Cmder, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) redisStatusCmdInterface
}

type redisBaseCmdInterface interface {
	Err() error
}

type redisStatusCmdInterface interface {
	redisBaseCmdInterface
}

type redisIntCmdInterface interface {
	redisBaseCmdInterface
	Val() int64
}

type redisStringCmdInterface interface {
	redisBaseCmdInterface
	Result() (string, error)
}

type redisScanCmdInterface interface {
	redisBaseCmdInterface
	Iterator() redisIteratorInterface
}

type redisIteratorInterface interface {
	Err() error
	Next(ctx context.Context) bool
	Val() string
}
