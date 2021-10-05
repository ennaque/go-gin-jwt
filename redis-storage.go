package gwt

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type redisStorage struct {
	con *redis.Client
}

func (rs *redisStorage) deleteTokensFromStorage(uuid ...string) error {
	if err := rs.con.Del(context.Background(), uuid...).Err(); err != nil {
		return errCannotDeleteToken
	}
	return nil
}

func (rs *redisStorage) saveTokensIntoStorage(access *accessTokenData,
	refresh *refreshTokenData) error {
	pipe := rs.con.TxPipeline()
	pipe.Set(context.Background(), access.uuid, access.token, time.Unix(access.expire, 0).Sub(time.Now()))
	pipe.Set(context.Background(), refresh.uuid, refresh.token, time.Unix(refresh.expire, 0).Sub(time.Now()))
	_, err := pipe.Exec(context.Background())
	if err != nil {
		return errCannotSaveToken
	}
	return nil
}

func (rs *redisStorage) isTokenExpired(uuid string, token string) error {
	tkn, err := rs.con.Get(context.Background(), uuid).Result()
	if err != nil {
		return errTokenExpired
	}
	if tkn != token {
		return errTokenInvalid
	}
	return nil
}
