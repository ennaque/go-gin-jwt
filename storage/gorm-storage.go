package storage

import (
	"github.com/ennaque/go-gin-jwt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var gwtTokensTablePrefix = "_gwt_token_data"

type gormStorage struct {
	con     *gorm.DB
	adapter gormAdapterInterface
}

func (gs *gormStorage) DeleteTokens(userId string, uuid ...string) error {
	err := gs.adapter.Transaction(gs.con, func(tx *gorm.DB) error {
		for _, id := range uuid {
			if err := gs.adapter.DeleteUnscoped(tx, &tokenData{UserId: userId, Uuid: id}, &tokenData{}).Error; err != nil {
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
	err := gs.adapter.Transaction(gs.con, func(tx *gorm.DB) error {
		if accessErr := gs.adapter.Create(tx, &tokenData{Token: accessToken, Uuid: accessUuid,
			Expire: accessExpire, UserId: userId, TokenType: "access"}).Error; accessErr != nil {
			return accessErr
		}
		if refreshErr := gs.adapter.Create(tx, &tokenData{Token: refreshToken, Uuid: refreshUuid,
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
	if err := gs.adapter.SelectFirst(gs.con,
		&tokenData{Token: token, Uuid: uuid, UserId: userId, TokenType: "refresh"}, &data).Error; err != nil {
		return err
	}
	return nil
}
func (gs *gormStorage) HasAccessToken(uuid string, token string, userId string) error {
	var data tokenData
	if err := gs.adapter.SelectFirst(gs.con,
		&tokenData{Token: token, Uuid: uuid, UserId: userId, TokenType: "access"}, &data).Error; err != nil {
		return err
	}
	return nil
}
func (gs *gormStorage) DeleteAllTokens(userId string) error {
	if err := gs.adapter.DeleteUnscoped(gs.con, &tokenData{UserId: userId}, &tokenData{}).Error; err != nil {
		return err
	}
	return nil
}

func InitGormStorage(con *gorm.DB, tablePrefix string) (gwt.StorageInterface, error) {
	adapter := &gormAdapter{}
	viper.Set("token_table_name", tablePrefix+gwtTokensTablePrefix)
	if err := adapter.AutoMigrate(con, &tokenData{}); err != nil {
		return nil, err
	}
	return &gormStorage{con: con, adapter: &gormAdapter{}}, nil
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
