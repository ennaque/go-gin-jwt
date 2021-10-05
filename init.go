package gwt

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

var (
	userIdClaim      = "user_id"
	accessUuidClaim  = "access_uuid"
	refreshUuidClaim = "refresh_uuid"
	expiredClaim     = "exp"
	authHeader       = "Authorization"
)

var availSigningMethods = map[string]string{
	"HS256": "true",
	"HS384": "true",
	"HS512": "true",
}

var (
	defaultSigningMethod     = "HS256"
	defaultAccessLifetime    = time.Minute * 10
	defaultRefreshLifetime   = time.Hour * 24
	defaultAuthHeadName      = "Bearer"
	defaultLoginResponseFunc = func(c *gin.Context, code int, accessToken string,
		accessExpire int64, refreshToken string, refreshExpire int64) {
		c.JSON(code, gin.H{
			"access_token":   accessToken,
			"refresh_token":  refreshToken,
			"access_expire":  strconv.FormatUint(uint64(accessExpire), 10),
			"refresh_expire": strconv.FormatUint(uint64(refreshExpire), 10),
		})
	}
	defaultErrResponseFunc = func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"error_code":    code,
			"error_message": message,
		})
		c.Abort()
	}
	defaultLogoutResponseFunc = func(c *gin.Context, code int) {
		c.JSON(code, gin.H{})
	}
)

func Init(settings Settings) (*Gwt, error) {
	if err := initStorage(&settings); err != nil {
		return nil, err
	}
	if settings.Authenticator == nil {
		return nil, errEmptyAuthenticator
	}
	if settings.AccessSecretKey == nil {
		return nil, errEmptyAccessSecretKey
	}
	if settings.RefreshSecretKey == nil {
		settings.RefreshSecretKey = settings.AccessSecretKey
	}
	if settings.SigningMethod == "" {
		settings.SigningMethod = defaultSigningMethod
	} else {
		if availSigningMethods[settings.SigningMethod] == "" {
			return nil, errUnknownSigningMethod
		}
	}
	if settings.AccessLifetime == 0 {
		settings.AccessLifetime = defaultAccessLifetime
	}
	if settings.RefreshLifetime == 0 {
		settings.RefreshLifetime = defaultRefreshLifetime
	}
	if settings.AuthHeadName == "" {
		settings.AuthHeadName = defaultAuthHeadName
	}
	if settings.LoginResponseFunc == nil {
		settings.LoginResponseFunc = defaultLoginResponseFunc
	}
	if settings.ErrResponseFunc == nil {
		settings.ErrResponseFunc = defaultErrResponseFunc
	}
	if settings.LogoutResponseFunc == nil {
		settings.LogoutResponseFunc = defaultLogoutResponseFunc
	}

	return &Gwt{settings: &settings}, nil
}

func initStorage(settings *Settings) error {
	if settings.RedisConnection == nil {
		return errEmptyRedisConnection
	}
	settings.storage = &redisStorage{con: settings.RedisConnection}
	return nil
}
