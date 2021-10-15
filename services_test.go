package gwt

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func testInitService(deleteAllTokensErr string) *Service {
	strgMock := new(storageMock)
	if deleteAllTokensErr == "" {
		strgMock.On("DeleteAllTokens", mock.Anything).Return(nil)
	} else {
		strgMock.On("DeleteAllTokens", mock.Anything).Return(errors.New(deleteAllTokensErr))
	}
	settings := getSettingsFixture()
	settings.Storage = strgMock

	return &Service{settings: settings}
}

func TestForceLogoutUserSuccess(t *testing.T) {
	service := testInitService("")
	err := service.ForceLogoutUser("1")

	assert.Nil(t, err)
}

func TestForceLogoutUserError(t *testing.T) {
	service := testInitService("delete error")
	err := service.ForceLogoutUser("1")

	assert.Error(t, err)
}
