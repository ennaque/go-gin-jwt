package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRedisDeleteTokensSuccess(t *testing.T) {
	mockSt := &redisAdapterMock{}
	mockSt.On("Del", mock.Anything).Return(nil)
	redisSt := &RedisStorage{adapter: mockSt}
	err := redisSt.DeleteTokens("1", "uuid")

	assert.Nil(t, err)
}
func TestRedisDeleteTokensError(t *testing.T) {
	mockSt := &redisAdapterMock{}
	mockSt.On("Del", mock.Anything).Return(errors.New("delete error"))
	redisSt := &RedisStorage{adapter: mockSt}
	err := redisSt.DeleteTokens("1", "uuid")

	assert.Error(t, err)
	assert.Equal(t, "delete error", err.Error())
}

func TestRedisSaveTokenSuccess(t *testing.T) {
	mockSt := &redisAdapterMock{}
	mockSt.On("SaveMultipleInPipe", mock.Anything).Return(nil)
	redisSt := &RedisStorage{adapter: mockSt}
	err := redisSt.SaveTokens("1", "auuid", "ruuid", 123,
		321, "atoken", "rtoken")

	assert.Nil(t, err)
}

func TestRedisSaveTokenError(t *testing.T) {
	mockSt := &redisAdapterMock{}
	mockSt.On("SaveMultipleInPipe", mock.Anything).Return(errors.New("save error"))
	redisSt := &RedisStorage{adapter: mockSt}
	err := redisSt.SaveTokens("1", "auuid", "ruuid", 123,
		321, "atoken", "rtoken")

	assert.Error(t, err)
	assert.Equal(t, "failed to save token from storage", err.Error())
}

func TestIsExpiredSuccess(t *testing.T) {
	mockSt := &redisAdapterMock{}
	mockSt.On("Get", mock.Anything).Return("token", nil)
	redisSt := &RedisStorage{adapter: mockSt}
	err := redisSt._isExpired("key", "token")

	assert.Nil(t, err)
}
func TestIsExpiredError(t *testing.T) {
	mockSt := &redisAdapterMock{}
	mockSt.On("Get", mock.Anything).Return("token", errors.New("err"))
	redisSt := &RedisStorage{adapter: mockSt}
	err := redisSt._isExpired("key", "token")

	assert.Error(t, err)
	assert.Equal(t, "token has expired", err.Error())
}
func TestIsExpiredInvalidTokenError(t *testing.T) {
	mockSt := &redisAdapterMock{}
	mockSt.On("Get", mock.Anything).Return("wrong", nil)
	redisSt := &RedisStorage{adapter: mockSt}
	err := redisSt._isExpired("key", "token")

	assert.Error(t, err)
	assert.Equal(t, "token is not valid", err.Error())
}
func TestRedisHasAccessTokenSuccess(t *testing.T) {
	mockSt := &redisAdapterMock{}
	mockSt.On("Get", mock.Anything).Return("token", nil)
	redisSt := &RedisStorage{adapter: mockSt}
	err := redisSt.HasAccessToken("uuid", "token", "1")

	assert.Nil(t, err)
}
func TestRedisHasAccessTokenError(t *testing.T) {
	mockSt := &redisAdapterMock{}
	mockSt.On("Get", mock.Anything).Return("token", errors.New("err"))
	redisSt := &RedisStorage{adapter: mockSt}
	err := redisSt.HasAccessToken("uuid", "token", "1")

	assert.Error(t, err)
}
func TestRedisHasRefreshTokenSuccess(t *testing.T) {
	mockSt := &redisAdapterMock{}
	mockSt.On("Get", mock.Anything).Return("token", nil)
	redisSt := &RedisStorage{adapter: mockSt}
	err := redisSt.HasRefreshToken("uuid", "token", "1")

	assert.Nil(t, err)
}
func TestRedisHasRefreshTokenError(t *testing.T) {
	mockSt := &redisAdapterMock{}
	mockSt.On("Get", mock.Anything).Return("token", errors.New("err"))
	redisSt := &RedisStorage{adapter: mockSt}
	err := redisSt.HasRefreshToken("uuid", "token", "1")

	assert.Error(t, err)
}

func TestGetStorageKeys(t *testing.T) {
	redisSt := &RedisStorage{adapter: &redisAdapterMock{}}
	res := redisSt._getStorageKeys("1", "uuid1", "uuid2")

	assert.Equal(t, "1_uuid1", res[0])
	assert.Equal(t, "1_uuid2", res[1])
}

func TestRedisDeleteAllTokensSuccess(t *testing.T) {
	iteratorMock := &redisIteratorMock{}
	iteratorMock.On("Val", mock.Anything).Return("test")
	iteratorMock.On("Next", mock.Anything).Return(nil)
	adapterMock := &redisAdapterMock{}
	adapterMock.On("GetScanIterator", mock.Anything).Return(iteratorMock)
	adapterMock.On("Del", mock.Anything).Return(nil)
	redisSt := &RedisStorage{adapter: adapterMock}
	err := redisSt.DeleteAllTokens("1")

	assert.Nil(t, err)

}
func TestRedisDeleteAllTokenNoAuthUserError(t *testing.T) {
	iteratorMock := &redisIteratorMock{}
	iteratorMock.On("Val", mock.Anything).Return("")
	iteratorMock.On("Next", mock.Anything).Return("err")
	adapterMock := &redisAdapterMock{}
	adapterMock.On("GetScanIterator", mock.Anything).Return(iteratorMock)
	adapterMock.On("Del", mock.Anything).Return(nil)
	redisSt := &RedisStorage{adapter: adapterMock}
	err := redisSt.DeleteAllTokens("1")

	assert.Error(t, err)
	assert.Equal(t, "user is not authenticated", err.Error())
}
