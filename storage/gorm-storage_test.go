package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
)

func TestDeleteTokensSuccess(t *testing.T) {
	dbMock := gormMock{}
	dbMock.On("Transaction", mock.Anything).Return(&dbMock)
	dbMock.On("Unscoped", mock.Anything).Return(&dbMock)
	dbMock.On("Where", mock.Anything).Return(&dbMock)
	dbMock.On("Delete", mock.Anything).Return(&gorm.DB{})

	gormStor := gormStorage{con: &dbMock}
	err := gormStor.DeleteTokens("1", "uuid")

	assert.Nil(t, err)
}

func TestDeleteTokensError(t *testing.T) {
	dbMock := gormMock{}
	gormError := &gorm.DB{}
	gormError.Error = errors.New("delete error")
	dbMock.On("Transaction", mock.Anything).Return(&dbMock)
	dbMock.On("Unscoped", mock.Anything).Return(&dbMock)
	dbMock.On("Where", mock.Anything).Return(&dbMock)
	dbMock.On("Delete", mock.Anything).Return(gormError)

	gormStor := gormStorage{con: &dbMock}
	err := gormStor.DeleteTokens("1", "uuid")

	assert.Error(t, err)
}

func TestSaveTokensSuccess(t *testing.T) {
	dbMock := gormMock{}
	dbMock.On("Transaction", mock.Anything).Return(&dbMock)
	dbMock.On("Create", mock.Anything).Return(nil)

	gormStor := gormStorage{con: &dbMock}
	err := gormStor.SaveTokens("1", "auuid", "ruuid", 123,
		1234, "access", "refresh")

	assert.Nil(t, err)
}

func TestSaveTokensAccessError(t *testing.T) {
	dbMock := gormMock{}
	gormError := &gorm.DB{}
	gormError.Error = errors.New("access error")
	dbMock.On("Transaction", mock.Anything).Return(&dbMock)
	dbMock.On("Create", mock.Anything).Return(1)

	gormStor := gormStorage{con: &dbMock}
	err := gormStor.SaveTokens("1", "auuid", "ruuid", 123,
		1234, "access", "refresh")

	assert.Error(t, err)
	assert.Equal(t, "access error", err.Error())
}

func TestSaveTokensRefreshError(t *testing.T) {
	dbMock := gormMock{}
	gormError := &gorm.DB{}
	gormError.Error = errors.New("refresh error")
	dbMock.On("Transaction", mock.Anything).Return(&dbMock)
	dbMock.On("Create", mock.Anything).Return(0)

	gormStor := gormStorage{con: &dbMock}
	err := gormStor.SaveTokens("1", "auuid", "ruuid", 123,
		1234, "access", "refresh")

	assert.Error(t, err)
	assert.Equal(t, "refresh error", err.Error())
}

func TestHasRefreshTokenSuccess(t *testing.T) {
	dbMock := gormMock{}
	dbMock.On("Where", mock.Anything).Return(&dbMock)
	dbMock.On("First", mock.Anything).Return(&gorm.DB{})

	gormStor := gormStorage{con: &dbMock}
	err := gormStor.HasRefreshToken("uuid", "token", "1")

	assert.Nil(t, err)
}

func TestHasRefreshTokenError(t *testing.T) {
	dbMock := gormMock{}
	dbErr := &gorm.DB{}
	dbErr.Error = errors.New("not found")
	dbMock.On("Where", mock.Anything).Return(&dbMock)
	dbMock.On("First", mock.Anything).Return(dbErr)

	gormStor := gormStorage{con: &dbMock}
	err := gormStor.HasRefreshToken("uuid", "token", "1")

	assert.Error(t, err)
}

func TestHasAccessTokenSuccess(t *testing.T) {
	dbMock := gormMock{}
	dbMock.On("Where", mock.Anything).Return(&dbMock)
	dbMock.On("First", mock.Anything).Return(&gorm.DB{})

	gormStor := gormStorage{con: &dbMock}
	err := gormStor.HasAccessToken("uuid", "token", "1")

	assert.Nil(t, err)
}

func TestHasAccessTokenError(t *testing.T) {
	dbMock := gormMock{}
	dbErr := &gorm.DB{}
	dbErr.Error = errors.New("not found")
	dbMock.On("Where", mock.Anything).Return(&dbMock)
	dbMock.On("First", mock.Anything).Return(dbErr)

	gormStor := gormStorage{con: &dbMock}
	err := gormStor.HasAccessToken("uuid", "token", "1")

	assert.Error(t, err)
}

func TestDeleteAllTokensSuccess(t *testing.T) {
	dbMock := gormMock{}
	dbMock.On("Unscoped", mock.Anything).Return(&dbMock)
	dbMock.On("Where", mock.Anything).Return(&dbMock)
	dbMock.On("Delete", mock.Anything).Return(&gorm.DB{})

	gormStor := gormStorage{con: &dbMock}
	err := gormStor.DeleteAllTokens("1")

	assert.Nil(t, err)
}

func TestDeleteAllTokensError(t *testing.T) {
	dbMock := gormMock{}
	dbErr := &gorm.DB{}
	dbErr.Error = errors.New("delete error")
	dbMock.On("Unscoped", mock.Anything).Return(&dbMock)
	dbMock.On("Where", mock.Anything).Return(&dbMock)
	dbMock.On("Delete", mock.Anything).Return(dbErr)

	gormStor := gormStorage{con: &dbMock}
	err := gormStor.DeleteAllTokens("1")

	assert.Error(t, err)
}

func TestInitGormStorageSuccess(t *testing.T) {
	dbMock := gormMock{}
	dbMock.On("AutoMigrate", mock.Anything).Return(nil)
	_, err := InitGormStorage(&dbMock, "prefix")

	assert.Nil(t, err)
}

func TestInitGormStorageError(t *testing.T) {
	dbMock := gormMock{}
	dbMock.On("AutoMigrate", mock.Anything).Return(errors.New("migrate error"))
	_, err := InitGormStorage(&dbMock, "prefix")

	assert.Error(t, err)
}
