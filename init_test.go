package gwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitSuccess(t *testing.T) {
	settings := getSettingsFixture()
	auth, err := Init(*settings)

	assert.Nil(t, err)
	assert.IsType(t, &Gwt{}, auth)
	assert.IsType(t, &Service{}, auth.Service)
	assert.IsType(t, &Handler{}, auth.Handler)
	assert.IsType(t, &Middleware{}, auth.Middleware)
}

func TestInitEmptyStorageError(t *testing.T) {
	settings := getSettingsFixture()
	settings.Storage = nil
	auth, err := Init(*settings)

	assert.Nil(t, auth)
	assert.Equal(t, err.Error(), ErrEmptyStorage.Error())
}

func TestInitEmptyAuthenticatorError(t *testing.T) {
	settings := getSettingsFixture()
	settings.Authenticator = nil
	auth, err := Init(*settings)

	assert.Nil(t, auth)
	assert.Equal(t, err.Error(), ErrEmptyAuthenticator.Error())
}

func TestInitEmptyGetUserFuncError(t *testing.T) {
	settings := getSettingsFixture()
	settings.GetUserFunc = nil
	auth, err := Init(*settings)

	assert.Nil(t, auth)
	assert.Equal(t, err.Error(), ErrEmptyGetUserFunc.Error())
}

func TestInitEmptyAccessSecretError(t *testing.T) {
	settings := getSettingsFixture()
	settings.AccessSecretKey = nil
	auth, err := Init(*settings)

	assert.Nil(t, auth)
	assert.Equal(t, err.Error(), ErrEmptyAccessSecretKey.Error())
}

func TestDefaultRefreshSecret(t *testing.T) {
	settings := getSettingsFixture()
	settings.AccessSecretKey = []byte("access_secret")
	settings.RefreshSecretKey = nil
	auth, err := Init(*settings)

	assert.Nil(t, err)
	assert.Equal(t, []byte("access_secret"), auth.Service.settings.RefreshSecretKey)
}

func TestDefaultSigningMethod(t *testing.T) {
	settings := getSettingsFixture()
	settings.SigningMethod = ""
	auth, err := Init(*settings)

	assert.Nil(t, err)
	assert.Equal(t, defaultSigningMethod, auth.Service.settings.SigningMethod)
}

func TestUnknownSigningMethodError(t *testing.T) {
	settings := getSettingsFixture()
	settings.SigningMethod = "Unknown"
	auth, err := Init(*settings)

	assert.Equal(t, ErrUnknownSigningMethod.Error(), err.Error())
	assert.Nil(t, auth)
}

func TestDefaultAccessLifetime(t *testing.T) {
	settings := getSettingsFixture()
	settings.AccessLifetime = 0
	auth, err := Init(*settings)

	assert.Nil(t, err)
	assert.Equal(t, defaultAccessLifetime, auth.Service.settings.AccessLifetime)
}

func TestDefaultRefreshLifetime(t *testing.T) {
	settings := getSettingsFixture()
	settings.RefreshLifetime = 0
	auth, err := Init(*settings)

	assert.Nil(t, err)
	assert.Equal(t, defaultRefreshLifetime, auth.Service.settings.RefreshLifetime)
}

func TestDefaultAuthHeaderName(t *testing.T) {
	settings := getSettingsFixture()
	settings.AuthHeadName = ""
	auth, err := Init(*settings)

	assert.Nil(t, err)
	assert.Equal(t, defaultAuthHeadName, auth.Service.settings.AuthHeadName)
}
