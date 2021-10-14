package gwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIsExpired(t *testing.T) {
	service := &tokenService{}
	err := service.isExpired(fmt.Sprint(time.Now().Add(time.Minute).Unix()))
	assert.Nil(t, err)

	err = service.isExpired(fmt.Sprint(time.Now().Add(-time.Minute).Unix()))
	assert.Equal(t, err, ErrTokenExpired)

	err = service.isExpired("wrong")
	assert.Error(t, err, ErrTokenInvalid)
}

func Test_CreateRefreshToken(t *testing.T) {
	service := &tokenService{}
	settingsFixture := getSettingsFixture()
	data, err := service._createRefreshToken(settingsFixture, "1", "auuid", "ruuid")
	assert.Nil(t, err)
	assert.IsType(t, &refreshTokenData{}, data)
	assert.Equal(t, "ruuid", data.uuid)
	assert.Equal(t, "auuid", data.accessUuid)
	assert.Equal(t, "1", data.userId)
	assert.IsType(t, "", data.token)
	assert.Equal(t, time.Now().Add(settingsFixture.RefreshLifetime).Unix(), data.expire)
}

func Test_CreateAccessToken(t *testing.T) {
	service := &tokenService{}
	settingsFixture := getSettingsFixture()
	data, err := service._createAccessToken(settingsFixture, "1", "auuid", "ruuid")
	assert.Nil(t, err)
	assert.IsType(t, &accessTokenData{}, data)
	assert.Equal(t, "auuid", data.uuid)
	assert.Equal(t, "ruuid", data.refreshUuid)
	assert.Equal(t, "1", data.userId)
	assert.IsType(t, "", data.token)
	assert.Equal(t, time.Now().Add(settingsFixture.AccessLifetime).Unix(), data.expire)
}

func TestParseToken(t *testing.T) {
	service := &tokenService{}
	settingsFixture := getSettingsFixture()
	token, tokenErr := service._createAccessToken(settingsFixture, "1", "auuid", "ruuid")
	assert.Nil(t, tokenErr)

	data, err := service.parseToken(token.token, settingsFixture.AccessSecretKey, "wrong_sign_method")
	assert.Nil(t, data)
	assert.Equal(t, err, ErrTokenInvalid)

	tkn, tknErr := service.parseToken(token.token, settingsFixture.AccessSecretKey, settingsFixture.SigningMethod)
	assert.Nil(t, tknErr)
	assert.IsType(t, &jwt.Token{}, tkn)
	assert.Equal(t, "auuid", tkn.Claims.(jwt.MapClaims)[accessUuidClaim])
	assert.Equal(t, "ruuid", tkn.Claims.(jwt.MapClaims)[refreshUuidClaim])
	assert.Equal(t, "1", tkn.Claims.(jwt.MapClaims)[userIdClaim])
	assert.Equal(t, float64(token.expire), tkn.Claims.(jwt.MapClaims)[expiredClaim])
}

func TestGetClaims(t *testing.T) {
	service := &tokenService{}
	settingsFixture := getSettingsFixture()
	token, _ := service._createAccessToken(settingsFixture, "1", "auuid", "ruuid")

	tkn, _ := service.parseToken(token.token, settingsFixture.AccessSecretKey, settingsFixture.SigningMethod)

	claims, claimsErr := service.getClaims(tkn, []string{accessUuidClaim, refreshUuidClaim, userIdClaim, expiredClaim})
	assert.Nil(t, claimsErr)
	assert.Equal(t, "auuid", claims[accessUuidClaim])
	assert.Equal(t, "ruuid", claims[refreshUuidClaim])
	assert.Equal(t, "1", claims[userIdClaim])
	assert.Equal(t, fmt.Sprint(float64(time.Now().Add(settingsFixture.AccessLifetime).Unix())), claims[expiredClaim])
}

func TestGetTokens(t *testing.T) {
	service := &tokenService{}
	settingsFixture := getSettingsFixture()

	access, refresh, err := service.getTokens(settingsFixture, "1")
	assert.Nil(t, err)
	assert.Equal(t, "1", access.userId)
	assert.Equal(t, "1", refresh.userId)
}
