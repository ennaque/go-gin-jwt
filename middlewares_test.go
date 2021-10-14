package gwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func testAuthMiddlewareInit(hasAccessTokenErr string, accessToken string,
	provideAccessToken bool, userFuncErr bool) *httptest.ResponseRecorder {
	strgMock := new(storageMock)
	if hasAccessTokenErr == "" {
		strgMock.On("HasAccessToken", mock.Anything).Return(nil)
	} else {
		strgMock.On("HasAccessToken", mock.Anything).Return(errors.New(hasAccessTokenErr))
	}

	settings := getSettingsFixture()
	settings.Storage = strgMock
	if userFuncErr {
		settings.GetUserFunc = func(userId string) (interface{}, error) {
			return nil, errors.New("get user error")
		}
	}

	mw := &Middleware{settings: settings}

	gin.SetMode(gin.TestMode)
	rr := httptest.NewRecorder()
	router := gin.Default()
	router.Use(mw.GetAuthMiddleware()).POST("/test-auth", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{})
	})
	request, _ := http.NewRequest(http.MethodPost, "/test-auth", nil)
	request.Header.Add("Content-Type", "application/json")
	if provideAccessToken {
		request.Header.Add("Authorization", "Bearer "+accessToken)
	}

	router.ServeHTTP(rr, request)

	return rr
}

func TestAuthMiddlewareSuccess(t *testing.T) {
	tService := tokenService{}
	accessData, _ := tService._createAccessToken(getSettingsFixture(), "1", "access", "refresh")
	rr := testAuthMiddlewareInit("", accessData.token, true, false)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestNoAuthHeaderError(t *testing.T) {
	tService := tokenService{}
	accessData, _ := tService._createAccessToken(getSettingsFixture(), "1", "access", "refresh")
	rr := testAuthMiddlewareInit("", accessData.token, false, false)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestParseTokenError(t *testing.T) {
	tService := tokenService{}
	accessData, _ := tService._createAccessToken(getSettingsFixture(), "1", "access", "refresh")
	rr := testAuthMiddlewareInit("", accessData.token+"wrong", true, false)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHasAccessTokenError(t *testing.T) {
	tService := tokenService{}
	accessData, _ := tService._createAccessToken(getSettingsFixture(), "1", "access", "refresh")
	rr := testAuthMiddlewareInit("", accessData.token+"wrong", true, false)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAccessTokenExpiredError(t *testing.T) {
	tService := tokenService{}
	settings := getSettingsFixture()
	settings.AccessLifetime = time.Nanosecond
	accessData, _ := tService._createAccessToken(settings, "1", "access", "refresh")
	rr := testAuthMiddlewareInit("", accessData.token, true, false)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestGetUserError(t *testing.T) {
	tService := tokenService{}
	accessData, _ := tService._createAccessToken(getSettingsFixture(), "1", "access", "refresh")
	rr := testAuthMiddlewareInit("", accessData.token, true, true)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
