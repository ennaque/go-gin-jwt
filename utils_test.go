package gwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetHeaderToken(t *testing.T) {
	token := "token_string"
	res, err := getHeaderToken("Bearer "+token, "Bearer")
	assert.Equal(t, res, token)
	assert.Nil(t, err)

	res, err = getHeaderToken("", "Bearer")
	assert.Equal(t, err, ErrNoAuthHeader)
	assert.Empty(t, res)

	res, err = getHeaderToken("Bearer token token", "Bearer")
	assert.Equal(t, err, ErrInvalidAuthHeader)
	assert.Empty(t, res)
}
