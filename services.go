package gwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"strings"
	"time"
)

func getTokens(settings *Settings, userId string) (*accessTokenData, *refreshTokenData, error) {
	refreshUuid := uuid.NewV4().String()
	accessUuid := uuid.NewV4().String()

	accessData, accessError := createAccessToken(settings, userId, accessUuid, refreshUuid)
	if accessError != nil {
		return nil, nil, accessError
	}
	refreshData, refreshError := createRefreshToken(settings, userId, accessUuid, refreshUuid)
	if refreshError != nil {
		return nil, nil, refreshError
	}
	if saveErr := saveTokensIntoStorage(settings, accessData.AccessExpire, accessData.AccessUuid,
		refreshData.RefreshExpire, refreshData.RefreshUuid, userId); saveErr != nil {
		return nil, nil, saveErr
	}
	return accessData, refreshData, nil
}

func createAccessToken(settings *Settings, userId string,
	accessUuid string, refreshUuid string) (*accessTokenData, error) {
	td := &accessTokenData{}
	td.AccessExpire = time.Now().Add(settings.AccessLifetime).Unix()
	td.AccessUuid = accessUuid
	td.RefreshUuid = refreshUuid
	td.userId = userId

	var err error
	td.AccessToken, err = createToken(
		settings.SigningMethod,
		jwt.MapClaims{accessUuidClaim: td.AccessUuid, userIdClaim: td.userId,
			expiredClaim: td.AccessExpire, refreshUuidClaim: td.RefreshUuid},
		settings.AccessSecretKey)
	if err != nil {
		return nil, ErrFailedToCreateAccessToken
	}
	return td, nil
}

func createRefreshToken(settings *Settings, userId string,
	accessUuid string, refreshUuid string) (*refreshTokenData, error) {
	td := &refreshTokenData{}
	td.RefreshExpire = time.Now().Add(settings.RefreshLifetime).Unix()
	td.RefreshUuid = refreshUuid
	td.AccessUuid = accessUuid
	td.userId = userId

	var err error
	td.RefreshToken, err = createToken(
		settings.SigningMethod,
		jwt.MapClaims{refreshUuidClaim: td.RefreshUuid,
			userIdClaim: td.userId, expiredClaim: td.RefreshExpire, accessUuidClaim: td.AccessUuid},
		settings.RefreshSecretKey)
	if err != nil {
		return nil, ErrFailedToCreateRefreshToken
	}
	return td, nil
}

func createToken(signingMethod string, claims jwt.MapClaims, secretKey []byte) (string, error) {
	rt := jwt.NewWithClaims(jwt.GetSigningMethod(signingMethod), claims)
	return rt.SignedString(secretKey)
}

func parseToken(tkn string, secret []byte, signingMethod string) (*jwt.Token, error) {
	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(signingMethod) != token.Method {
			return nil, ErrInvalidSigningMethod
		}

		return secret, nil
	})
	if err != nil {
		return nil, ErrTokenExpired
	}
	if !token.Valid {
		return nil, ErrTokenInvalid
	}
	return token, nil
}

func getClaims(token *jwt.Token, claimNames []string) (map[string]string, error) {
	res := map[string]string{}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		for _, el := range claimNames {
			val, ok := claims[el].(string)
			if !ok {
				return nil, ErrTokenInvalid
			}
			res[el] = val
		}
		return res, nil
	}
	return nil, ErrTokenInvalid
}

func getHeaderToken(headerString string, authHeadName string) (string, error) {
	if headerString == "" {
		return "", ErrNoAuthHeader
	}

	parts := strings.SplitN(headerString, " ", 2)
	if !(len(parts) == 2 && parts[0] == authHeadName) {
		return "", ErrInvalidAuthHeader
	}

	return parts[1], nil
}
