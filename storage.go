package gwt

import (
	"context"
	"time"
)

func deleteTokensFromStorage(settings *Settings, uuid ...string) error {
	redis := settings.RedisConnection
	if err := redis.Del(context.Background(), uuid...).Err(); err != nil {
		return ErrCannotDeleteToken
	}
	return nil
}

func saveTokensIntoStorage(settings *Settings, accessExpire int64, accessUuid string,
	refreshExpire int64, refreshUuid string, userId string) error {
	redisClient := settings.RedisConnection
	pipe := redisClient.TxPipeline()
	pipe.Set(context.Background(), accessUuid, userId, time.Unix(accessExpire, 0).Sub(time.Now()))
	pipe.Set(context.Background(), refreshUuid, userId, time.Unix(refreshExpire, 0).Sub(time.Now()))
	_, err := pipe.Exec(context.Background())
	if err != nil {
		return ErrCannotSaveToken
	}
	return nil
}

func isTokenExpired(settings *Settings, uuid string) error {
	redis := settings.RedisConnection
	_, err := redis.Get(context.Background(), uuid).Result()
	if err != nil {
		return ErrTokenExpired
	}
	return nil
}
