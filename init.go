package gwt

import (
	"github.com/gin-gonic/gin"
	"time"
)

var (
	userIdClaim        = "user_id"
	accessUuidClaim    = "access_uuid"
	refreshUuidClaim   = "refresh_uuid"
	expiredClaim       = "exp"
	authHeader         = "Authorization"
	userIdRequestParam = "user_id"
	UserKey            = "user"
)

var availSigningMethods = map[string]string{
	"HS256": "true",
	"HS384": "true",
	"HS512": "true",
}

type DefaultLoginResponse struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	AccessExpire  int64  `json:"access_expire"`
	RefreshExpire int64  `json:"refresh_expire"`
}

type DefaultErrResponse struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

type DefaultLogoutResponse struct{}

var (
	defaultSigningMethod     = "HS256"
	defaultAccessLifetime    = time.Minute * 10
	defaultRefreshLifetime   = time.Hour * 24
	defaultAuthHeadName      = "Bearer"
	defaultLoginResponseFunc = func(c *gin.Context, code int, accessToken string,
		accessExpire int64, refreshToken string, refreshExpire int64) {
		c.JSON(code, DefaultLoginResponse{
			AccessToken:   accessToken,
			RefreshToken:  refreshToken,
			AccessExpire:  accessExpire,
			RefreshExpire: refreshExpire,
		})
	}
	defaultErrResponseFunc = func(c *gin.Context, code int, message string) {
		c.JSON(code, DefaultErrResponse{ErrorCode: code, ErrorMessage: message})
		c.Abort()
	}
	defaultLogoutResponseFunc = func(c *gin.Context, code int) {
		c.JSON(code, DefaultLogoutResponse{})
	}
)

type Gwt struct {
	Service    *Service
	Handler    *Handler
	Middleware *Middleware
}

func Init(settings Settings) (*Gwt, error) {
	if settings.Storage == nil {
		return nil, ErrEmptyStorage
	}
	if settings.Authenticator == nil {
		return nil, ErrEmptyAuthenticator
	}
	if settings.GetUserFunc == nil {
		return nil, ErrEmptyGetUserFunc
	}
	if settings.AccessSecretKey == nil {
		return nil, ErrEmptyAccessSecretKey
	}
	if settings.RefreshSecretKey == nil {
		settings.RefreshSecretKey = settings.AccessSecretKey
	}
	if settings.SigningMethod == "" {
		settings.SigningMethod = defaultSigningMethod
	} else {
		if availSigningMethods[settings.SigningMethod] == "" {
			return nil, ErrUnknownSigningMethod
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

	return &Gwt{
		Middleware: &Middleware{settings: &settings},
		Handler:    &Handler{settings: &settings},
		Service:    &Service{settings: &settings},
	}, nil
}
