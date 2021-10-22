package storage

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

func TestDeleteTokensSuccess(t *testing.T) {
	adapterMock := gormAdapterMock{}
	adapterMock.On("DeleteUnscoped", mock.Anything).Return(&gorm.DB{})
	gormSt := &gormStorage{con: &gorm.DB{}, adapter: &adapterMock}
	err := gormSt.DeleteTokens("1", "uuid")

	assert.Nil(t, err)
}

func TestDeleteTokensError(t *testing.T) {
	adapterMock := gormAdapterMock{}
	ret := &gorm.DB{}
	ret.Error = errors.New("err")
	adapterMock.On("DeleteUnscoped", mock.Anything).Return(ret)
	gormSt := &gormStorage{con: &gorm.DB{}, adapter: &adapterMock}
	err := gormSt.DeleteTokens("1", "uuid")

	assert.Error(t, err)
	assert.Equal(t, "err", err.Error())
}

func TestSaveTokensSuccess(t *testing.T) {
	adapterMock := gormAdapterMock{}
	adapterMock.On("Transaction", mock.Anything).Return(nil)
	adapterMock.On("Create", mock.Anything).Return(nil)
	gormSt := &gormStorage{con: &gorm.DB{}, adapter: &adapterMock}
	err := gormSt.SaveTokens("1", "auuid", "ruuid",
		123, 321, "atoken", "rtoken")

	assert.Nil(t, err)
}

func TestSaveTokensAccessError(t *testing.T) {
	adapterMock := gormAdapterMock{}
	adapterMock.On("Transaction", mock.Anything).Return(nil)
	adapterMock.On("Create", mock.Anything).Return("accessErr")
	gormSt := &gormStorage{con: &gorm.DB{}, adapter: &adapterMock}
	err := gormSt.SaveTokens("1", "auuid", "ruuid",
		123, 321, "atoken", "rtoken")

	assert.Error(t, err)
	assert.Equal(t, "accessErr", err.Error())
}

func TestSaveTokensRefreshError(t *testing.T) {
	adapterMock := gormAdapterMock{}
	adapterMock.On("Transaction", mock.Anything).Return(nil)
	adapterMock.On("Create", mock.Anything).Return("refreshErr")
	gormSt := &gormStorage{con: &gorm.DB{}, adapter: &adapterMock}
	err := gormSt.SaveTokens("1", "auuid", "ruuid",
		123, 321, "atoken", "rtoken")

	assert.Error(t, err)
	assert.Equal(t, "refreshErr", err.Error())
}

func TestHasRefreshTokenSuccess(t *testing.T) {
	adapterMock := gormAdapterMock{}
	adapterMock.On("SelectFirst", mock.Anything).Return(&gorm.DB{})
	gormSt := &gormStorage{con: &gorm.DB{}, adapter: &adapterMock}
	err := gormSt.HasRefreshToken("uuid", "rtoken", "1")

	assert.Nil(t, err)
}

func TestHasRefreshTokenError(t *testing.T) {
	adapterMock := gormAdapterMock{}
	ret := &gorm.DB{}
	ret.Error = errors.New("err")
	adapterMock.On("SelectFirst", mock.Anything).Return(ret)
	gormSt := &gormStorage{con: &gorm.DB{}, adapter: &adapterMock}
	err := gormSt.HasRefreshToken("uuid", "rtoken", "1")

	assert.Error(t, err)
	assert.Equal(t, "err", err.Error())
}

func TestHasAccessTokenSuccess(t *testing.T) {
	adapterMock := gormAdapterMock{}
	adapterMock.On("SelectFirst", mock.Anything).Return(&gorm.DB{})
	gormSt := &gormStorage{con: &gorm.DB{}, adapter: &adapterMock}
	err := gormSt.HasAccessToken("uuid", "atoken", "1")

	assert.Nil(t, err)
}

func TestHasAccessTokenError(t *testing.T) {
	adapterMock := gormAdapterMock{}
	ret := &gorm.DB{}
	ret.Error = errors.New("err")
	adapterMock.On("SelectFirst", mock.Anything).Return(ret)
	gormSt := &gormStorage{con: &gorm.DB{}, adapter: &adapterMock}
	err := gormSt.HasAccessToken("uuid", "atoken", "1")

	assert.Error(t, err)
	assert.Equal(t, "err", err.Error())
}

func TestDeleteAllTokensSuccess(t *testing.T) {
	adapterMock := gormAdapterMock{}
	adapterMock.On("DeleteUnscoped", mock.Anything).Return(&gorm.DB{})
	gormSt := &gormStorage{con: &gorm.DB{}, adapter: &adapterMock}
	err := gormSt.DeleteAllTokens("1")

	assert.Nil(t, err)
}

func TestDeleteAllTokensError(t *testing.T) {
	adapterMock := gormAdapterMock{}
	ret := &gorm.DB{}
	ret.Error = errors.New("err")
	adapterMock.On("DeleteUnscoped", mock.Anything).Return(ret)
	gormSt := &gormStorage{con: &gorm.DB{}, adapter: &adapterMock}
	err := gormSt.DeleteAllTokens("1")

	assert.Error(t, err)
	assert.Equal(t, "err", err.Error())
}

func TestInitGormStorageSuccess(t *testing.T) {
	//dbMock := gormMock{}
	//dbMock.On("AutoMigrate", mock.Anything).Return(nil)
	//_, err := InitGormStorage(&dbMock, "prefix")
	//
	//assert.Nil(t, err)
}

func TestInitGormStorageError(t *testing.T) {
	//dbMock := gormMock{}
	//dbMock.On("AutoMigrate", mock.Anything).Return(errors.New("migrate error"))
	//_, err := InitGormStorage(&dbMock, "prefix")
	//
	//assert.Error(t, err)
}

func TestGormTableName(t *testing.T) {
	viper.Set("token_table_name", "name")
	td := &tokenData{}

	assert.Equal(t, "name", td.TableName())
}
