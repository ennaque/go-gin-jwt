package storage

import (
	"database/sql"
	"errors"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

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
func (m *gormMock) Where(query interface{}, args ...interface{}) (tx gormInterface) {
	arguments := m.Called()
	return arguments.Get(0).(gormInterface)
}
func (m *gormMock) Unscoped() (tx gormInterface) {
	arguments := m.Called()
	return arguments.Get(0).(gormInterface)
}
func (m *gormMock) Transaction(fc func(tx gormInterface) error, opts ...*sql.TxOptions) (err error) {
	return fc(m.Called().Get(0).(gormInterface))
}
