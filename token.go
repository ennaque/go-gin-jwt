package gwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
	"math/big"
	"time"
)

type tokenService struct{}

func (ts *tokenService) isExpired(expireStr string) error {
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

func (ts *tokenService) getTokens(settings *Settings, userId string) (*accessTokenData, *refreshTokenData, error) {
	refreshUuid := uuid.NewV4().String()
	accessUuid := uuid.NewV4().String()

	accessData, accessError := ts._createAccessToken(settings, userId, accessUuid, refreshUuid)
	if accessError != nil {
		return nil, nil, accessError
	}
	refreshData, refreshError := ts._createRefreshToken(settings, userId, accessUuid, refreshUuid)
	if refreshError != nil {
		return nil, nil, refreshError
	}
	return accessData, refreshData, nil
}

func (ts *tokenService) getClaims(token *jwt.Token, claimNames []string) (map[string]string, error) {
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

func (ts *tokenService) parseToken(tkn string, secret []byte, signingMethod string) (*jwt.Token, error) {
	token, err := jwt.Parse(tkn, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(signingMethod) != token.Method {
			return nil, ErrInvalidSigningMethod
		}

		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return token, nil
}

func (ts *tokenService) _createAccessToken(settings *Settings, userId string,
	accessUuid string, refreshUuid string) (*accessTokenData, error) {
	td := &accessTokenData{}
	td.expire = time.Now().Add(settings.AccessLifetime).Unix()
	td.uuid = accessUuid
	td.refreshUuid = refreshUuid
	td.userId = userId

	var err error
	td.token, err = ts._createToken(
		settings.SigningMethod,
		jwt.MapClaims{accessUuidClaim: td.uuid, userIdClaim: td.userId,
			expiredClaim: td.expire, refreshUuidClaim: td.refreshUuid},
		settings.AccessSecretKey)
	if err != nil {
		return nil, ErrFailedToCreateAccessToken
	}
	return td, nil
}

func (ts *tokenService) _createRefreshToken(settings *Settings, userId string,
	accessUuid string, refreshUuid string) (*refreshTokenData, error) {
	td := &refreshTokenData{}
	td.expire = time.Now().Add(settings.RefreshLifetime).Unix()
	td.uuid = refreshUuid
	td.accessUuid = accessUuid
	td.userId = userId

	var err error
	td.token, err = ts._createToken(
		settings.SigningMethod,
		jwt.MapClaims{refreshUuidClaim: td.uuid,
			userIdClaim: td.userId, expiredClaim: td.expire, accessUuidClaim: td.accessUuid},
		settings.RefreshSecretKey)
	if err != nil {
		return nil, ErrFailedToCreateRefreshToken
	}
	return td, nil
}

func (ts *tokenService) _createToken(signingMethod string, claims jwt.MapClaims, secretKey []byte) (string, error) {
	rt := jwt.NewWithClaims(jwt.GetSigningMethod(signingMethod), claims)
	return rt.SignedString(secretKey)
}
