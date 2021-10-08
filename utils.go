package gwt

import (
	"strings"
)

func getHeaderToken(headerString string, authHeadName string) (string, error) {
	if headerString == "" {
		return "", errNoAuthHeader
	}

	parts := strings.SplitN(headerString, " ", 2)
	if !(len(parts) == 2 && parts[0] == authHeadName) {
		return "", errInvalidAuthHeader
	}

	return parts[1], nil
}
