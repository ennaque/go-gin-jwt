package gwt

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Middleware struct {
	settings *Settings
}

func (mw *Middleware) GetAuthMiddleware() gin.HandlerFunc {
	return mw.authMiddleware()
}

func (mw *Middleware) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := getHeaderToken(c.Request.Header.Get(authHeader), mw.settings.AuthHeadName)
		if err != nil {
			mw.settings.ErrResponseFunc(c, http.StatusUnauthorized, err.Error())
			return
		}
		parsedToken, err := parseToken(accessToken, mw.settings.AccessSecretKey, mw.settings.SigningMethod)
		if err != nil {
			mw.settings.ErrResponseFunc(c, http.StatusBadRequest, err.Error())
			return
		}
		claims, getClaimsErr := getClaims(parsedToken, []string{accessUuidClaim, userIdClaim, expiredClaim})
		if err != nil {
			mw.settings.ErrResponseFunc(c, http.StatusBadRequest, getClaimsErr.Error())
			return
		}
		if tokenExpErr := mw.settings.Storage.HasAccessToken(claims[accessUuidClaim], accessToken, claims[userIdClaim]); tokenExpErr != nil {
			mw.settings.ErrResponseFunc(c, http.StatusUnauthorized, tokenExpErr.Error())
			return
		}
		if expErr := isExpired(claims[expiredClaim]); expErr != nil {
			mw.settings.ErrResponseFunc(c, http.StatusUnauthorized, expErr.Error())
			return
		}
		user, userErr := mw.settings.GetUserFunc(claims[userIdClaim])
		if userErr != nil {
			mw.settings.ErrResponseFunc(c, http.StatusInternalServerError, userErr.Error())
			return
		}
		c.Set(UserKey, user)
		c.Next()
	}
}
