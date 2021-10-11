package gwt

type StorageInterface interface {
	DeleteTokens(userId string, uuid ...string) error
	SaveTokens(userId string, accessUuid string, refreshUuid string, accessExpire int64,
		refreshExpire int64, accessToken string, refreshToken string) error
	HasRefreshToken(uuid string, token string, userId string) error
	HasAccessToken(uuid string, token string, userId string) error
	DeleteAllTokens(userId string) error
}
