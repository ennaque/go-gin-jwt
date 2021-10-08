package gwt

type storageInterface interface {
	deleteTokens(userId string, uuid ...string) error
	saveTokens(accessTokenData *accessTokenData, refreshTokenData *refreshTokenData) error
	isRefreshExpired(uuid string, token string, userId string) error
	isAccessExpired(uuid string, token string, userId string) error
	deleteAllTokens(userId string) error
}
