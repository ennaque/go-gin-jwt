package storage

import (
	"context"
	"github.com/ennaque/go-gin-jwt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisStorage struct {
	adapter redisAdapterInterface
}

func (rs *RedisStorage) DeleteTokens(userId string, uuid ...string) error {
	if err := rs.adapter.Del(context.Background(), rs._getStorageKeys(userId, uuid...)...); err != nil {
		return err
	}
	return nil
}

func (rs *RedisStorage) SaveTokens(
	userId string,
	accessUuid string,
	refreshUuid string,
	accessExpire int64,
	refreshExpire int64,
	accessToken string,
	refreshToken string) error {
	_, err := rs.adapter.SaveMultipleInPipe(
		context.Background(),
		redisValue{key: rs._getStorageKey("a"+userId, accessUuid), value: accessToken,
			expiration: time.Unix(accessExpire, 0).Sub(time.Now())},
		redisValue{key: rs._getStorageKey("r"+userId, refreshUuid), value: refreshToken,
			expiration: time.Unix(refreshExpire, 0).Sub(time.Now())},
	)
	if err != nil {
		return gwt.ErrCannotSaveToken
	}
	return nil
}

func (rs *RedisStorage) HasRefreshToken(uuid string, token string, userId string) error {
	return rs._isExpired("a"+rs._getStorageKey(userId, uuid), token)
}

func (rs *RedisStorage) HasAccessToken(uuid string, token string, userId string) error {
	return rs._isExpired("r"+rs._getStorageKey(userId, uuid), token)
}

func (rs *RedisStorage) DeleteAllTokens(userId string) error {
	userIdUuidKeys := rs._getUserIdUuidStorageKeys(userId)
	if len(userIdUuidKeys) == 0 {
		return gwt.ErrNotAuthUser
	}
	return rs.adapter.Del(context.Background(), userIdUuidKeys...)
}

func (rs *RedisStorage) _isExpired(key string, token string) error {
	tkn, err := rs.adapter.Get(context.Background(), key)
	if err != nil {
		return gwt.ErrTokenExpired
	}
	if tkn != token {
		return gwt.ErrTokenInvalid
	}
	return nil
}

func (rs *RedisStorage) _getUserIdUuidStorageKeys(userId string) []string {
	var keysToDelete []string
	iter := rs.adapter.GetScanIterator(context.Background(), 0, "[ar]"+userId+"_*", 0)
	for iter.Next(context.Background()) {
		keysToDelete = append(keysToDelete, iter.Val())
	}
	return keysToDelete
}

func (rs *RedisStorage) _getStorageKeys(userId string, uuids ...string) []string {
	var keys []string
	for _, key := range uuids {
		keys = append(keys, rs._getStorageKey(userId, key))
	}
	return keys
}

func (rs *RedisStorage) _getStorageKey(userId string, uuid string) string {
	return userId + "_" + uuid
}

func InitRedisStorage(client *redis.Client) gwt.StorageInterface {
	return &RedisStorage{adapter: &redisAdapter{con: client}}
}
