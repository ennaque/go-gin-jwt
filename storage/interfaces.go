package storage

import (
	"database/sql"
	"gorm.io/gorm"
)

type gormInterface interface {
	Transaction(fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) (err error)
	Unscoped() (tx *gorm.DB)
	Where(query interface{}, args ...interface{}) (tx *gorm.DB)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
	Create(value interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	AutoMigrate(dst ...interface{}) error
}
