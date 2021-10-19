package storage

import (
	"github.com/ennaque/go-gin-jwt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var gwtTokensTablePrefix = "_gwt_token_data"

type gormStorage struct {
	con interface{ gormInterface }
}

func (gs *gormStorage) DeleteTokens(userId string, uuid ...string) error {
	err := gs.con.Transaction(func(tx gormInterface) error {
		for _, id := range uuid {
			if err := tx.Unscoped().Where(&tokenData{UserId: userId, Uuid: id}).Delete(&tokenData{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (gs *gormStorage) SaveTokens(userId string, accessUuid string, refreshUuid string, accessExpire int64,
	refreshExpire int64, accessToken string, refreshToken string) error {
	err := gs.con.Transaction(func(tx gormInterface) error {
		if accessErr := tx.Create(&tokenData{Token: accessToken, Uuid: accessUuid,
			Expire: accessExpire, UserId: userId, TokenType: "access"}).Error; accessErr != nil {
			return accessErr
		}
		if refreshErr := tx.Create(&tokenData{Token: refreshToken, Uuid: refreshUuid,
			Expire: refreshExpire, UserId: userId, TokenType: "refresh"}).Error; refreshErr != nil {
			return refreshErr
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (gs *gormStorage) HasRefreshToken(uuid string, token string, userId string) error {
	var data tokenData
	if err := gs.con.Where(&tokenData{Token: token, Uuid: uuid, UserId: userId, TokenType: "refresh"}).First(&data).Error; err != nil {
		return gwt.ErrTokenExpired
	}
	return nil
}
func (gs *gormStorage) HasAccessToken(uuid string, token string, userId string) error {
	var data tokenData
	if err := gs.con.Where(&tokenData{Token: token, Uuid: uuid, UserId: userId, TokenType: "access"}).First(&data).Error; err != nil {
		return gwt.ErrTokenExpired
	}
	return nil
}
func (gs *gormStorage) DeleteAllTokens(userId string) error {
	if err := gs.con.Unscoped().Where(&tokenData{UserId: userId}).Delete(&tokenData{}).Error; err != nil {
		return err
	}
	return nil
}

func InitGormStorage(con gormInterface, tablePrefix string) (gwt.StorageInterface, error) {
	if err := initDb(con, tablePrefix); err != nil {
		return nil, err
	}
	return &gormStorage{con: con}, nil
}

func initDb(con gormInterface, tablePrefix string) error {
	viper.Set("token_table_name", tablePrefix+gwtTokensTablePrefix)
	return con.AutoMigrate(&tokenData{})
}

type tokenData struct {
	gorm.Model
	Token     string `gorm:"type:string;not null;unique;index" valid:"required"`
	Uuid      string `gorm:"type:string;not null;unique;index" valid:"required"`
	Expire    int64  `gorm:"type:uint;not null;" valid:"required"`
	UserId    string `gorm:"type:string;not null;" valid:"required"`
	TokenType string `gorm:"type:string;not null" valid:"required"`
}

func (td *tokenData) TableName() string {
	return viper.Get("token_table_name").(string)
}
