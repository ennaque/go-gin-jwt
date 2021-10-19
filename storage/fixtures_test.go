package storage

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"time"
)

type redisCmdFullInterface interface {
	redisIntCmdInterface
	redisStringCmdInterface
	redisScanCmdInterface
	redisIntCmdInterface
}

// redis mocks
type redisMock struct {
	mock.Mock
	redisCmdMock      redisCmdFullInterface
	redisPipelineMock redisPipelineInterface
}

func (m *redisMock) Del(ctx context.Context, keys ...string) redisIntCmdInterface {
	return m.redisCmdMock
}
func (m *redisMock) Get(ctx context.Context, key string) redisStringCmdInterface {
	return m.redisCmdMock
}
func (m *redisMock) Scan(ctx context.Context, cursor uint64, match string, count int64) redisScanCmdInterface {
	return m.redisCmdMock
}
func (m *redisMock) TxPipeline() redisPipelineInterface {
	return m.redisPipelineMock
}

type redisPipelineMock struct {
	mock.Mock
}

func (m *redisPipelineMock) Exec(ctx context.Context) ([]redis.Cmder, error) {
	if m.Called().Get(0) == nil {
		return nil, nil
	}

	return nil, m.Called().Get(0).(error)
}
func (m *redisPipelineMock) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) redisStatusCmdInterface {
	return &redisCmdMock{}
}

type redisIteratorMock struct {
	mock.Mock
	nextCalled bool
}

func (m *redisIteratorMock) Err() error {
	args := m.Called()
	return args.Error(0)
}
func (m *redisIteratorMock) Next(ctx context.Context) bool {
	if m.nextCalled {
		return false
	}
	m.nextCalled = true
	args := m.Called()
	return args.Bool(0)
}
func (m *redisIteratorMock) Val() string {
	return m.Called().String(0)
}

type redisCmdMock struct {
	mock.Mock
}

func (m *redisCmdMock) Err() error {
	args := m.Called()
	return args.Error(0)
}
func (m *redisCmdMock) Val() int64 {
	args := m.Called()
	return int64(args.Int(0))
}
func (m *redisCmdMock) Result() (string, error) {
	if _, ok := m.Called().Get(0).(string); !ok {
		return "", m.Called().Get(0).(error)
	}
	return m.Called().String(0), nil
}
func (m *redisCmdMock) Iterator() redisIteratorInterface {
	return m.Called().Get(0).(redisIteratorInterface)
}

// gorm mocks
type gormMock struct {
	mock.Mock
}

func (m *gormMock) AutoMigrate(dst ...interface{}) error {
	args := m.Called()
	return args.Error(0)
}
func (m *gormMock) First(dest interface{}, conds ...interface{}) (tx *gorm.DB) {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}
func (m *gormMock) Create(value interface{}) (tx *gorm.DB) {
	if m.Called().Get(0) == nil {
		return &gorm.DB{}
	}
	if value.(*tokenData).TokenType == "access" {
		if m.Called().Get(0).(int) == 1 {
			g := &gorm.DB{}
			g.Error = errors.New("access error")
			return g
		} else {
			return &gorm.DB{}
		}
	} else {
		if m.Called().Get(0).(int) == 0 {
			g := &gorm.DB{}
			g.Error = errors.New("refresh error")
			return g
		} else {
			return &gorm.DB{}
		}
	}
}
func (m *gormMock) Delete(value interface{}, conds ...interface{}) (tx *gorm.DB) {
	args := m.Called()
	return args.Get(0).(*gorm.DB)
}
func (m *gormMock) Where(query interface{}, args ...interface{}) (tx interface{ gormInterface }) {
	arguments := m.Called()
	return arguments.Get(0).(interface{ gormInterface })
}
func (m *gormMock) Unscoped() (tx interface{ gormInterface }) {
	arguments := m.Called()
	return arguments.Get(0).(interface{ gormInterface })
}
func (m *gormMock) Transaction(fc func(tx interface{ gormInterface }) error, opts ...*sql.TxOptions) (err error) {
	return fc(m.Called().Get(0).(interface{ gormInterface }))
}
