package storage

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func testRedisDeleteTokensInit(error string) error {
	redisCmdMock := &redisCmdMock{}
	if error == "" {
		redisCmdMock.On("Err", mock.Anything).Return(nil)
	} else {
		redisCmdMock.On("Err", mock.Anything).Return(errors.New(error))
	}
	redisMock := &redisMock{redisCmdMock: redisCmdMock}
	redisStorage := RedisStorage{Con: redisMock}
	return redisStorage.DeleteTokens("1", "uuid")
}
func TestRedisDeleteTokensSuccess(t *testing.T) {
	err := testRedisDeleteTokensInit("")
	assert.Nil(t, err)
}
func TestRedisDeleteTokensError(t *testing.T) {
	err := testRedisDeleteTokensInit("delete error")
	assert.Error(t, err)
}

func testInitRedisSaveToken(error string) error {
	pipelineMock := &redisPipelineMock{}
	pipelineMock.On("Set", mock.Anything).Return(nil)
	if error == "" {
		pipelineMock.On("Exec", mock.Anything).Return(nil)
	} else {
		pipelineMock.On("Exec", mock.Anything).Return(errors.New(error))
	}
	redisMock := &redisMock{redisPipelineMock: pipelineMock}
	redisMock.On("TxPipeline", mock.Anything).Return(nil)
	redisStorage := RedisStorage{Con: redisMock}

	return redisStorage.SaveTokens("1", "auuid", "ruuid", 123,
		1234, "access", "refresh")
}
func TestRedisSaveTokenSuccess(t *testing.T) {
	err := testInitRedisSaveToken("")
	assert.Nil(t, err)
}

func TestRedisSaveTokenError(t *testing.T) {
	err := testInitRedisSaveToken("save error")
	assert.Error(t, err)
}

func testRedisIsExpiredInit(error string) RedisStorage {
	redisCmdMock := &redisCmdMock{}
	redisCmdMock.On("Err", mock.Anything).Return(nil)
	if error == "" {
		redisCmdMock.On("Result", mock.Anything).Return("token")
	} else {
		redisCmdMock.On("Result", mock.Anything).Return(errors.New(error))
	}
	redisMock := &redisMock{redisCmdMock: redisCmdMock}
	redisStorage := RedisStorage{Con: redisMock}

	return redisStorage
}
func TestIsExpiredSuccess(t *testing.T) {
	redisStorage := testRedisIsExpiredInit("")
	err := redisStorage._isExpired("key", "token")
	assert.Nil(t, err)
}
func TestIsExpiredError(t *testing.T) {
	redisStorage := testRedisIsExpiredInit("expired")
	err := redisStorage._isExpired("key", "token")
	assert.Error(t, err)
	assert.Equal(t, "token has expired", err.Error())
}
func TestIsExpiredInvalidTokenError(t *testing.T) {
	redisStorage := testRedisIsExpiredInit("")
	err := redisStorage._isExpired("key", "invalid_token")
	assert.Error(t, err)
	assert.Equal(t, "token is not valid", err.Error())
}
func TestRedisHasAccessTokenSuccess(t *testing.T) {
	redisStorage := testRedisIsExpiredInit("")
	err := redisStorage.HasAccessToken("auuid", "token", "1")
	assert.Nil(t, err)
}
func TestRedisHasAccessTokenError(t *testing.T) {
	redisStorage := testRedisIsExpiredInit("")
	err := redisStorage.HasAccessToken("auuid", "invalid_token", "1")
	assert.Error(t, err)
	assert.Equal(t, "token is not valid", err.Error())
}
func TestRedisHasRefreshTokenSuccess(t *testing.T) {
	redisStorage := testRedisIsExpiredInit("")
	err := redisStorage.HasRefreshToken("ruuid", "token", "1")
	assert.Nil(t, err)
}
func TestRedisHasRefreshTokenError(t *testing.T) {
	redisStorage := testRedisIsExpiredInit("")
	err := redisStorage.HasRefreshToken("ruuid", "invalid_token", "1")
	assert.Error(t, err)
	assert.Equal(t, "token is not valid", err.Error())
}

func TestGetStorageKeys(t *testing.T) {
	redisStorage := RedisStorage{Con: &redisMock{}}
	res := redisStorage._getStorageKeys("1", "uuid1", "uuid2")

	assert.Equal(t, "1_uuid1", res[0])
	assert.Equal(t, "1_uuid2", res[1])
}

func TestRedisDeleteAllTokensSuccess(t *testing.T) {
	redisIteratorMock := &redisIteratorMock{}
	redisIteratorMock.On("Val", mock.Anything).Return("token")
	redisIteratorMock.On("Next", mock.Anything).Return(true)
	redisCmdMock := &redisCmdMock{}
	redisCmdMock.On("Iterator", mock.Anything).Return(redisIteratorMock)
	redisCmdMock.On("Err", mock.Anything).Return(nil)
	redisCmdMock.On("Del", mock.Anything).Return(nil)
	redisMock := &redisMock{redisCmdMock: redisCmdMock}
	redisStorage := RedisStorage{Con: redisMock}
	err := redisStorage.DeleteAllTokens("1")

	assert.Nil(t, err)
}
func TestRedisDeleteAllTokenNoAuthUserError(t *testing.T) {
	redisIteratorMock := &redisIteratorMock{}
	redisIteratorMock.On("Next", mock.Anything).Return(false)
	redisCmdMock := &redisCmdMock{}
	redisCmdMock.On("Iterator", mock.Anything).Return(redisIteratorMock)
	redisCmdMock.On("Err", mock.Anything).Return(nil)
	redisCmdMock.On("Del", mock.Anything).Return(nil)
	redisMock := &redisMock{redisCmdMock: redisCmdMock}
	redisStorage := RedisStorage{Con: redisMock}
	err := redisStorage.DeleteAllTokens("1")

	assert.Error(t, err)
	assert.Equal(t, "user is not authenticated", err.Error())
}
