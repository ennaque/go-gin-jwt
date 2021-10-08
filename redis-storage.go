package gwt

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisStorage struct {
	con *redis.Client
}

func (rs *redisStorage) deleteTokens(userId string, uuid ...string) error {
	if err := rs.con.Del(context.Background(), rs._getStorageKeys(userId, uuid...)...).Err(); err != nil {
		return errCannotDeleteToken
	}
	return nil
}

func (rs *redisStorage) saveTokens(access *accessTokenData, refresh *refreshTokenData) error {
	pipe := rs.con.TxPipeline()
	pipe.Set(context.Background(), rs._getStorageKey("a"+access.userId, access.uuid), access.token, time.Unix(access.expire, 0).Sub(time.Now()))
	pipe.Set(context.Background(), rs._getStorageKey("r"+refresh.userId, refresh.uuid), refresh.token, time.Unix(refresh.expire, 0).Sub(time.Now()))
	_, err := pipe.Exec(context.Background())
	if err != nil {
		return errCannotSaveToken
	}
	return nil
}

func (rs *redisStorage) isAccessExpired(uuid string, token string, userId string) error {
	return rs._isExpired("a"+rs._getStorageKey(userId, uuid), token)
}

func (rs *redisStorage) isRefreshExpired(uuid string, token string, userId string) error {
	fmt.Println(userId)
	return rs._isExpired("r"+rs._getStorageKey(userId, uuid), token)
}

func (rs *redisStorage) deleteAllTokens(userId string) error {
	userIdUuidKeys := rs._getUserIdUuidStorageKeys(userId)
	if len(userIdUuidKeys) == 0 {
		return errNotAuthUser
	}
	rs.con.Del(context.Background(), userIdUuidKeys...)
	return nil
}

func (rs *redisStorage) _isExpired(key string, token string) error {
	tkn, err := rs.con.Get(context.Background(), key).Result()
	fmt.Println(err)
	fmt.Println(key)
	if err != nil {
		return errTokenExpired
	}
	if tkn != token {
		return errTokenInvalid
	}
	return nil
}

func (rs *redisStorage) _getUserIdUuidStorageKeys(userId string) []string {
	var keysToDelete []string
	iter := rs.con.Scan(context.Background(), 0, "[ar]"+userId+"_*", 0).Iterator()
	for iter.Next(context.Background()) {
		keysToDelete = append(keysToDelete, iter.Val())
	}
	return keysToDelete
}

func (rs *redisStorage) _getStorageKeys(userId string, uuids ...string) []string {
	var keys []string
	for _, key := range uuids {
		keys = append(keys, rs._getStorageKey(userId, key))
	}
	return keys
}

func (rs *redisStorage) _getStorageKey(userId string, uuid string) string {
	return userId + "_" + uuid
}
