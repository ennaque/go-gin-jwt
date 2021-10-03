package gwt

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (gwt *Gwt) loginHandler(c *gin.Context) {
	userId, err := gwt.settings.Authenticator(c)
	if err != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusUnauthorized, err.Error())
		return
	}

	accessData, refreshData, er := getTokens(gwt.settings, userId)
	if er != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusInternalServerError, err.Error())
		return
	}

	gwt.settings.LoginResponseFunc(c, http.StatusOK, accessData.AccessToken,
		accessData.AccessExpire, refreshData.RefreshToken, refreshData.RefreshExpire)
}
func (gwt *Gwt) refreshHandler(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusBadRequest, ErrRefreshTokenIsNotProvided.Error())
		return
	}
	parsedToken, err := parseToken(mapToken["refresh_token"], gwt.settings.RefreshSecretKey, gwt.settings.SigningMethod)
	if err != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusBadRequest, err.Error())
		return
	}
	claims, getClaimsErr := getClaims(parsedToken, []string{refreshUuidClaim, accessUuidClaim, userIdClaim})
	if err != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusBadRequest, getClaimsErr.Error())
		return
	}
	if tokenExpErr := isTokenExpired(gwt.settings, claims[refreshUuidClaim]); tokenExpErr != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusUnauthorized, tokenExpErr.Error())
		return
	}
	if deleteRefreshErr := deleteTokensFromStorage(gwt.settings, claims[refreshUuidClaim],
		claims[accessUuidClaim]); deleteRefreshErr != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusInternalServerError, err.Error())
		return
	}
	accessData, refreshData, er := getTokens(gwt.settings, claims[userIdClaim])
	if er != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusInternalServerError, err.Error())
		return
	}

	gwt.settings.LoginResponseFunc(c, http.StatusOK, accessData.AccessToken,
		accessData.AccessExpire, refreshData.RefreshToken, refreshData.RefreshExpire)
}
func (gwt *Gwt) logoutHandler(c *gin.Context) {
	accessToken, err := getHeaderToken(c.Request.Header.Get(authHeader), gwt.settings.AuthHeadName)
	if err != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusUnauthorized, err.Error())
		return
	}
	parsedToken, err := parseToken(accessToken, gwt.settings.AccessSecretKey, gwt.settings.SigningMethod)
	if err != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusBadRequest, err.Error())
		return
	}
	claims, getClaimsErr := getClaims(parsedToken, []string{refreshUuidClaim, accessUuidClaim, userIdClaim})
	if err != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusBadRequest, getClaimsErr.Error())
		return
	}
	if tokenExpErr := isTokenExpired(gwt.settings, claims[accessUuidClaim]); tokenExpErr != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusUnauthorized, tokenExpErr.Error())
		return
	}
	if deleteRefreshErr := deleteTokensFromStorage(gwt.settings, claims[accessUuidClaim],
		claims[refreshUuidClaim]); deleteRefreshErr != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusInternalServerError, err.Error())
		return
	}
	gwt.settings.LogoutResponseFunc(c, http.StatusOK)
}
func (gwt *Gwt) authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := getHeaderToken(c.Request.Header.Get(authHeader), gwt.settings.AuthHeadName)
		if err != nil {
			gwt.settings.ErrResponseFunc(c, http.StatusUnauthorized, err.Error())
			return
		}
		parsedToken, err := parseToken(accessToken, gwt.settings.AccessSecretKey, gwt.settings.SigningMethod)
		if err != nil {
			gwt.settings.ErrResponseFunc(c, http.StatusBadRequest, err.Error())
			return
		}
		claims, getClaimsErr := getClaims(parsedToken, []string{accessUuidClaim, userIdClaim})
		if err != nil {
			gwt.settings.ErrResponseFunc(c, http.StatusBadRequest, getClaimsErr.Error())
			return
		}
		if tokenExpErr := isTokenExpired(gwt.settings, claims[accessUuidClaim]); tokenExpErr != nil {
			gwt.settings.ErrResponseFunc(c, http.StatusUnauthorized, tokenExpErr.Error())
			return
		}
		c.Set(userIdClaim, claims[userIdClaim])
		c.Next()
	}
}
