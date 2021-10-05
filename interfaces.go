package gwt

type storageInterface interface {
	deleteTokensFromStorage(uuid ...string) error
	saveTokensIntoStorage(accessTokenData *accessTokenData, refreshTokenData *refreshTokenData) error
	isTokenExpired(uuid string, token string) error
}
