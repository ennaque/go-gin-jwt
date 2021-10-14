package gwt

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func testLoginInit(authErr string, saveTokensErr string) (*httptest.ResponseRecorder, error) {
	strgMock := new(storageMock)
	if saveTokensErr == "" {
		strgMock.On("SaveTokens", mock.Anything).Return(nil)
	} else {
		strgMock.On("SaveTokens", mock.Anything).Return(errors.New(saveTokensErr))
	}
	settings := getSettingsFixture()
	settings.Authenticator = func(c *gin.Context) (string, error) {
		if authErr == "" {
			return "1", nil
		}
		return "", errors.New(authErr)
	}
	settings.Storage = strgMock
	handler := &Handler{settings: settings}

	gin.SetMode(gin.TestMode)
	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/login", handler.GetLoginHandler())

	request, err := http.NewRequest(http.MethodPost, "/login", nil)
	router.ServeHTTP(rr, request)

	return rr, err
}

func TestLoginSuccess(t *testing.T) {
	rr, err := testLoginInit("", "")

	var res map[string]string
	json.NewDecoder(rr.Body).Decode(&res)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.IsType(t, "", res["access_expire"])
	assert.IsType(t, "", res["access_token"])
	assert.IsType(t, "", res["refresh_expire"])
	assert.IsType(t, "", res["refresh_token"])
}

func TestLoginAuthenticatorError(t *testing.T) {
	rr, err := testLoginInit("invalid credentials", "")
	var res map[string]string
	json.NewDecoder(rr.Body).Decode(&res)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "invalid credentials", res["error_message"])
}

func TestLoginSaveTokensError(t *testing.T) {
	rr, err := testLoginInit("", "save error")
	var res map[string]string
	json.NewDecoder(rr.Body).Decode(&res)

	assert.NoError(t, err)
	assert.Equal(t, "save error", res["error_message"])
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func testRefreshInit(saveTokensErr string, hasRefreshTokenErr string,
	deleteTokensErr string, refreshToken string) *httptest.ResponseRecorder {
	strgMock := new(storageMock)
	if saveTokensErr == "" {
		strgMock.On("SaveTokens", mock.Anything).Return(nil)
	} else {
		strgMock.On("SaveTokens", mock.Anything).Return(errors.New(saveTokensErr))
	}
	if hasRefreshTokenErr == "" {
		strgMock.On("HasRefreshToken", mock.Anything).Return(nil)
	} else {
		strgMock.On("HasRefreshToken", mock.Anything).Return(errors.New(hasRefreshTokenErr))
	}
	if deleteTokensErr == "" {
		strgMock.On("DeleteTokens", mock.Anything).Return(nil)
	} else {
		strgMock.On("DeleteTokens", mock.Anything).Return(errors.New(deleteTokensErr))
	}

	settings := getSettingsFixture()
	settings.Storage = strgMock

	handler := &Handler{settings: settings}

	gin.SetMode(gin.TestMode)
	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/refresh", handler.GetRefreshHandler())
	params, _ := json.Marshal(map[string]string{"refresh_token": refreshToken})
	request, _ := http.NewRequest(http.MethodPost, "/refresh", bytes.NewBuffer(params))
	request.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(rr, request)

	return rr
}

func TestRefreshSuccess(t *testing.T) {
	tService := tokenService{}
	refreshData, _ := tService._createRefreshToken(getSettingsFixture(), "1", "access", "refresh")

	rr := testRefreshInit("", "", "", refreshData.token)

	var resRefresh map[string]string
	json.NewDecoder(rr.Body).Decode(&resRefresh)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.IsType(t, "", resRefresh["access_expire"])
	assert.IsType(t, "", resRefresh["access_token"])
	assert.IsType(t, "", resRefresh["refresh_expire"])
	assert.IsType(t, "", resRefresh["refresh_token"])
}

func TestRefreshParseTokenError(t *testing.T) {
	tService := tokenService{}
	refreshData, _ := tService._createRefreshToken(getSettingsFixture(), "1", "access", "refresh")

	rr := testRefreshInit("", "", "", refreshData.token+"wrong")

	var resRefresh map[string]string
	json.NewDecoder(rr.Body).Decode(&resRefresh)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "token is not valid", resRefresh["error_message"])
}

func TestRefreshHasRefreshTokenError(t *testing.T) {
	tService := tokenService{}
	refreshData, _ := tService._createRefreshToken(getSettingsFixture(), "1", "access", "refresh")

	rr := testRefreshInit("", "expired", "", refreshData.token)

	var resRefresh map[string]string
	json.NewDecoder(rr.Body).Decode(&resRefresh)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "expired", resRefresh["error_message"])
}

func TestRefreshDeleteTokensError(t *testing.T) {
	tService := tokenService{}
	refreshData, _ := tService._createRefreshToken(getSettingsFixture(), "1", "access", "refresh")

	rr := testRefreshInit("", "", "delete error", refreshData.token)

	var resRefresh map[string]string
	json.NewDecoder(rr.Body).Decode(&resRefresh)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "delete error", resRefresh["error_message"])
}

func TestRefreshSaveTokensError(t *testing.T) {
	tService := tokenService{}
	refreshData, _ := tService._createRefreshToken(getSettingsFixture(), "1", "access", "refresh")

	rr := testRefreshInit("save error", "", "", refreshData.token)

	var resRefresh map[string]string
	json.NewDecoder(rr.Body).Decode(&resRefresh)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "save error", resRefresh["error_message"])
}

func TestRefreshExpiredError(t *testing.T) {
	tService := tokenService{}
	settings := getSettingsFixture()
	settings.RefreshLifetime = time.Nanosecond
	refreshData, _ := tService._createRefreshToken(settings, "1", "access", "refresh")

	rr := testRefreshInit("save error", "", "", refreshData.token)

	var resRefresh map[string]string
	json.NewDecoder(rr.Body).Decode(&resRefresh)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "token has expired", resRefresh["error_message"])
}

func testLogoutInit(hasAccessTokenError string, deleteTokensError string,
	accessToken string, provideHeader bool) *httptest.ResponseRecorder {
	strgMock := new(storageMock)
	if hasAccessTokenError == "" {
		strgMock.On("HasAccessToken", mock.Anything).Return(nil)
	} else {
		strgMock.On("HasAccessToken", mock.Anything).Return(errors.New(hasAccessTokenError))
	}
	if deleteTokensError == "" {
		strgMock.On("DeleteTokens", mock.Anything).Return(nil)
	} else {
		strgMock.On("DeleteTokens", mock.Anything).Return(errors.New(deleteTokensError))
	}

	settings := getSettingsFixture()
	settings.Storage = strgMock

	handler := &Handler{settings: settings}
	gin.SetMode(gin.TestMode)
	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/logout", handler.GetLogoutHandler())
	request, _ := http.NewRequest(http.MethodPost, "/logout", nil)
	request.Header.Add("Content-Type", "application/json")
	if provideHeader {
		request.Header.Add("Authorization", "Bearer "+accessToken)
	}
	router.ServeHTTP(rr, request)
	return rr
}

func TestLogoutSuccess(t *testing.T) {
	tService := tokenService{}
	accessData, _ := tService._createAccessToken(getSettingsFixture(), "1", "access", "refresh")
	rr := testLogoutInit("", "", accessData.token, true)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestLogoutNoHeaderError(t *testing.T) {
	tService := tokenService{}
	accessData, _ := tService._createAccessToken(getSettingsFixture(), "1", "access", "refresh")
	rr := testLogoutInit("", "", accessData.token, false)
	var res map[string]string
	json.NewDecoder(rr.Body).Decode(&res)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "no auth header provided", res["error_message"])
}

func TestLogoutInvalidTokenError(t *testing.T) {
	tService := tokenService{}
	accessData, _ := tService._createAccessToken(getSettingsFixture(), "1", "access", "refresh")
	rr := testLogoutInit("", "", accessData.token+"wrong", true)
	var res map[string]string
	json.NewDecoder(rr.Body).Decode(&res)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "token is not valid", res["error_message"])
}

func TestLogoutHasAccessTokenError(t *testing.T) {
	tService := tokenService{}
	accessData, _ := tService._createAccessToken(getSettingsFixture(), "1", "access", "refresh")
	rr := testLogoutInit("expired", "", accessData.token, true)
	var res map[string]string
	json.NewDecoder(rr.Body).Decode(&res)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "expired", res["error_message"])
}

func TestLogoutExpiredError(t *testing.T) {
	tService := tokenService{}
	settings := getSettingsFixture()
	settings.AccessLifetime = time.Nanosecond
	accessData, _ := tService._createAccessToken(settings, "1", "access", "refresh")
	rr := testLogoutInit("", "", accessData.token, true)
	var res map[string]string
	json.NewDecoder(rr.Body).Decode(&res)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Equal(t, "token has expired", res["error_message"])
}

func TestLogoutDeleteTokensError(t *testing.T) {
	tService := tokenService{}
	accessData, _ := tService._createAccessToken(getSettingsFixture(), "1", "access", "refresh")
	rr := testLogoutInit("", "delete error", accessData.token, true)
	var res map[string]string
	json.NewDecoder(rr.Body).Decode(&res)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "delete error", res["error_message"])
}

func testForceLogoutInit(deleteAllTokensErr string, provideParams bool) *httptest.ResponseRecorder {
	strgMock := new(storageMock)
	if deleteAllTokensErr == "" {
		strgMock.On("DeleteAllTokens", mock.Anything).Return(nil)
	} else {
		strgMock.On("DeleteAllTokens", mock.Anything).Return(errors.New(deleteAllTokensErr))
	}
	settings := getSettingsFixture()
	settings.Storage = strgMock

	handler := &Handler{settings: settings}
	gin.SetMode(gin.TestMode)
	rr := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/force-logout", handler.GetForceLogoutHandler())

	var params []byte
	if provideParams {
		params, _ = json.Marshal(map[string]string{"user_id": "1"})
	}
	request, _ := http.NewRequest(http.MethodPost, "/force-logout", bytes.NewBuffer(params))
	request.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(rr, request)
	return rr
}

func TestForceLogoutSuccess(t *testing.T) {
	rr := testForceLogoutInit("", true)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestForceLogoutInvalidParamsError(t *testing.T) {
	rr := testForceLogoutInit("", false)

	var res map[string]string
	json.NewDecoder(rr.Body).Decode(&res)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "user id is not provided", res["error_message"])
}

func TestForceLogoutDeleteTokensError(t *testing.T) {
	rr := testForceLogoutInit("delete error", true)

	var res map[string]string
	json.NewDecoder(rr.Body).Decode(&res)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "delete error", res["error_message"])
}
