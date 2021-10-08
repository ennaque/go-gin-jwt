package gwt

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"time"
)

type refreshTokenData struct {
	userId     string
	token      string
	uuid       string
	expire     int64
	accessUuid string
}

type accessTokenData struct {
	userId      string
	token       string
	uuid        string
	expire      int64
	refreshUuid string
}

type Settings struct {

	// SigningMethod signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SigningMethod string

	// AccessSecretKey used for signing. Required.
	AccessSecretKey []byte

	// RefreshSecretKey used for signing. Optional, AccessSecretKey is used by default.
	RefreshSecretKey []byte

	// AccessLifetime is a duration that an access token is valid. Optional, ten minutes by defaults.
	AccessLifetime time.Duration

	// RefreshLifetime is a duration that a refresh token is valid. Optional, one day by defaults.
	RefreshLifetime time.Duration

	// AuthHeadName is a string in the header. Default value is "Bearer"
	AuthHeadName string

	// Callback function that should perform the authentication of the user based on login info.
	// Must return user id as string. Required.
	Authenticator func(c *gin.Context) (string, error)

	LoginResponseFunc func(c *gin.Context, code int, accessToken string,
		accessExpire int64, refreshToken string, refreshExpire int64)

	LogoutResponseFunc func(c *gin.Context, code int)

	ErrResponseFunc func(c *gin.Context, code int, message string)

	RedisConnection *redis.Client

	storage storageInterface
}

type Gwt struct {
	settings *Settings
}

func (gwt *Gwt) GetLoginHandler() func(c *gin.Context) {
	return gwt.loginHandler
}
func (gwt *Gwt) GetRefreshHandler() func(c *gin.Context) {
	return gwt.refreshHandler
}
func (gwt *Gwt) GetLogoutHandler() func(c *gin.Context) {
	return gwt.logoutHandler
}
func (gwt *Gwt) GetAuthMiddleware() gin.HandlerFunc {
	return gwt.authMiddleware()
}
func (gwt *Gwt) GetForceLogoutHandler() func(c *gin.Context) {
	return gwt.forceLogoutHandler
}
