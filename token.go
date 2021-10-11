package gwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"math/big"
	"time"
)

func isExpired(expireStr string) error {
	expireFloat, _, err := big.ParseFloat(expireStr, 10, 0, big.ToNearestEven)
	if err != nil {
		return ErrTokenInvalid
	}
	expire, _ := expireFloat.Int64()
	if expire > time.Now().Unix() {
		return nil
	}
	return ErrTokenExpired
}

func getTokens(settings *Settings, userId string) (*accessTokenData, *refreshTokenData, error) {
	refreshUuid := uuid.NewV4().String()
	accessUuid := uuid.NewV4().String()

	accessData, accessError := _createAccessToken(settings, userId, accessUuid, refreshUuid)
	if accessError != nil {
		return nil, nil, accessError
	}
	refreshData, refreshError := _createRefreshToken(settings, userId, accessUuid, refreshUuid)
	if refreshError != nil {
		return nil, nil, refreshError
	}
	if saveErr := settings.Storage.SaveTokens(userId, accessUuid, refreshUuid, accessData.expire,
		refreshData.expire, accessData.token, refreshData.token); saveErr != nil {
		return nil, nil, saveErr
	}
	return accessData, refreshData, nil
}

func getClaims(token *jwt.Token, claimNames []string) (map[string]string, error) {
	res := map[string]string{}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		for _, el := range claimNames {
			res[el] = fmt.Sprint(claims[el])
		}
		return res, nil
	}
	return nil, ErrTokenInvalid
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

func _createAccessToken(settings *Settings, userId string,
	accessUuid string, refreshUuid string) (*accessTokenData, error) {
	td := &accessTokenData{}
	td.expire = time.Now().Add(settings.AccessLifetime).Unix()
	td.uuid = accessUuid
	td.refreshUuid = refreshUuid
	td.userId = userId

	var err error
	td.token, err = _createToken(
		settings.SigningMethod,
		jwt.MapClaims{accessUuidClaim: td.uuid, userIdClaim: td.userId,
			expiredClaim: td.expire, refreshUuidClaim: td.refreshUuid},
		settings.AccessSecretKey)
	if err != nil {
		return nil, ErrFailedToCreateAccessToken
	}
	return td, nil
}

func _createRefreshToken(settings *Settings, userId string,
	accessUuid string, refreshUuid string) (*refreshTokenData, error) {
	td := &refreshTokenData{}
	td.expire = time.Now().Add(settings.RefreshLifetime).Unix()
	td.uuid = refreshUuid
	td.accessUuid = accessUuid
	td.userId = userId

	var err error
	td.token, err = _createToken(
		settings.SigningMethod,
		jwt.MapClaims{refreshUuidClaim: td.uuid,
			userIdClaim: td.userId, expiredClaim: td.expire, accessUuidClaim: td.accessUuid},
		settings.RefreshSecretKey)
	if err != nil {
		return nil, ErrFailedToCreateRefreshToken
	}
	return td, nil
}

func _createToken(signingMethod string, claims jwt.MapClaims, secretKey []byte) (string, error) {
	rt := jwt.NewWithClaims(jwt.GetSigningMethod(signingMethod), claims)
	return rt.SignedString(secretKey)
}
