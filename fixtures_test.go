package gwt

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"strconv"
	"time"
)

type storageMock struct {
	mock.Mock
}

func (m *storageMock) DeleteTokens(userId string, uuid ...string) error {
	args := m.Called()
	return args.Error(0)
}
func (m *storageMock) SaveTokens(userId string, accessUuid string, refreshUuid string, accessExpire int64,
	refreshExpire int64, accessToken string, refreshToken string) error {
	args := m.Called()
	return args.Error(0)
}
func (m *storageMock) HasRefreshToken(uuid string, token string, userId string) error {
	args := m.Called()
	return args.Error(0)
}
func (m *storageMock) HasAccessToken(uuid string, token string, userId string) error {
	args := m.Called()
	return args.Error(0)
}
func (m *storageMock) DeleteAllTokens(userId string) error {
	args := m.Called()
	return args.Error(0)
}

func getSettingsFixture() *Settings {
	return &Settings{
		SigningMethod:    "HS256",
		AccessSecretKey:  []byte("super_secret"),
		RefreshSecretKey: []byte("super_secret"),
		AccessLifetime:   time.Minute * 1,
		RefreshLifetime:  time.Minute * 2,
		AuthHeadName:     "Bearer",
		Authenticator: func(c *gin.Context) (string, error) {
			return "1", nil
		},
		GetUserFunc: func(userId string) (interface{}, error) {
			return userId, nil
		},
		LoginResponseFunc: func(c *gin.Context, code int, accessToken string,
			accessExpire int64, refreshToken string, refreshExpire int64) {
			c.JSON(code, gin.H{
				"access_token":   accessToken,
				"refresh_token":  refreshToken,
				"access_expire":  strconv.FormatUint(uint64(accessExpire), 10),
				"refresh_expire": strconv.FormatUint(uint64(refreshExpire), 10),
			})
		},
		LogoutResponseFunc: func(c *gin.Context, code int) {
			c.JSON(code, gin.H{})
		},
		ErrResponseFunc: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"error_code":    code,
				"error_message": message,
			})
			c.Abort()
		},
		Storage: &storageMock{},
	}
}
