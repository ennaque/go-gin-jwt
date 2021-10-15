package storage

import (
	"database/sql"
	"gorm.io/gorm"
)

type gormInterface interface {
	Transaction(fc func(tx gormInterface) error, opts ...*sql.TxOptions) (err error)
	Unscoped() (tx gormInterface)
	Where(query interface{}, args ...interface{}) (tx gormInterface)
	Delete(value interface{}, conds ...interface{}) (tx *gorm.DB)
	Create(value interface{}) (tx *gorm.DB)
	First(dest interface{}, conds ...interface{}) (tx *gorm.DB)
	AutoMigrate(dst ...interface{}) error
}
