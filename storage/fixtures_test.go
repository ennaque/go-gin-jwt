package storage

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type gormAdapterMock struct {
	mock.Mock
}

func (m *gormAdapterMock) Transaction(db *gorm.DB, fc func(tx *gorm.DB) error) error {
	return fc(db)
}
func (m *gormAdapterMock) DeleteUnscoped(db *gorm.DB, query interface{}, model interface{}) *gorm.DB {
	return m.Called().Get(0).(*gorm.DB)
}
func (m *gormAdapterMock) Create(db *gorm.DB, value interface{}) *gorm.DB {
	if m.Called().Get(0) == nil {
		return &gorm.DB{}
	}
	if value.(*tokenData).TokenType == "access" {
		if m.Called().String(0) == "accessErr" {
			g := &gorm.DB{}
			g.Error = errors.New("accessErr")
			return g
		}
		return &gorm.DB{}
	} else {
		if m.Called().String(0) == "refreshErr" {
			g := &gorm.DB{}
			g.Error = errors.New("refreshErr")
			return g
		}
		return &gorm.DB{}
	}
}
func (m *gormAdapterMock) SelectFirst(db *gorm.DB, query interface{}, destination interface{}) *gorm.DB {
	return m.Called().Get(0).(*gorm.DB)
}
func (m *gormAdapterMock) AutoMigrate(db *gorm.DB, dst ...interface{}) error {
	return m.Called().Error(0)
}

type redisAdapterMock struct {
	mock.Mock
}

func (m *redisAdapterMock) Del(ctx context.Context, keys ...string) error {
	return m.Called().Error(0)
}
func (m *redisAdapterMock) SaveMultipleInPipe(ctx context.Context, values ...redisValue) ([]redis.Cmder, error) {
	return nil, m.Called().Error(0)
}
func (m *redisAdapterMock) Get(ctx context.Context, key string) (string, error) {
	return m.Called().String(0), m.Called().Error(1)
}
func (m *redisAdapterMock) GetScanIterator(ctx context.Context, cursor uint64, match string, count int64) redisIteratorInterface {
	return m.Called().Get(0).(redisIteratorInterface)
}

type redisIteratorMock struct {
	mock.Mock
	init bool
}

func (m *redisIteratorMock) Err() error {
	return m.Called().Error(0)
}
func (m *redisIteratorMock) Next(ctx context.Context) bool {
	if m.Called().Get(0) != nil {
		return false
	}
	if m.init == true {
		return false
	}
	m.setInit()
	return true
}
func (m *redisIteratorMock) Val() string {
	return m.Called().String(0)
}
func (m *redisIteratorMock) setInit() {
	m.init = true
}
