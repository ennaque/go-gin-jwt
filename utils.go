package gwt

import (
	"strings"
)

func getHeaderToken(headerString string, authHeadName string) (string, error) {
	if headerString == "" {
		return "", ErrNoAuthHeader
	}

	parts := strings.SplitN(headerString, " ", -1)
	if !(len(parts) == 2 && parts[0] == authHeadName) {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}
