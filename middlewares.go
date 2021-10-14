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
		service := &tokenService{}
		accessToken, getErr := getHeaderToken(c.Request.Header.Get(authHeader), mw.settings.AuthHeadName)
		if getErr != nil {
			mw.settings.ErrResponseFunc(c, http.StatusUnauthorized, getErr.Error())
			return
		}
		parsedToken, parseErr := service.parseToken(accessToken, mw.settings.AccessSecretKey, mw.settings.SigningMethod)
		if parseErr != nil {
			mw.settings.ErrResponseFunc(c, http.StatusBadRequest, parseErr.Error())
			return
		}
		claims, getClaimsErr := service.getClaims(parsedToken, []string{accessUuidClaim, userIdClaim, expiredClaim})
		if getClaimsErr != nil {
			mw.settings.ErrResponseFunc(c, http.StatusBadRequest, getClaimsErr.Error())
			return
		}
		if tokenExpErr := mw.settings.Storage.HasAccessToken(claims[accessUuidClaim], accessToken, claims[userIdClaim]); tokenExpErr != nil {
			mw.settings.ErrResponseFunc(c, http.StatusUnauthorized, tokenExpErr.Error())
			return
		}
		if expErr := service.isExpired(claims[expiredClaim]); expErr != nil {
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
