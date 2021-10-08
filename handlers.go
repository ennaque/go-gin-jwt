package gwt

import (
	"fmt"
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

	gwt.settings.LoginResponseFunc(c, http.StatusOK, accessData.token,
		accessData.expire, refreshData.token, refreshData.expire)
}
func (gwt *Gwt) refreshHandler(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBind(&mapToken); err != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusBadRequest, errRefreshTokenIsNotProvided.Error())
		return
	}
	parsedToken, err := parseToken(mapToken[refreshTokenRequestParam], gwt.settings.RefreshSecretKey, gwt.settings.SigningMethod)
	if err != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusBadRequest, err.Error())
		return
	}
	claims, getClaimsErr := getClaims(parsedToken, []string{refreshUuidClaim, accessUuidClaim, userIdClaim})
	if err != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusBadRequest, getClaimsErr.Error())
		return
	}
	fmt.Println(claims)
	if tokenExpErr := gwt.settings.storage.isRefreshExpired(claims[refreshUuidClaim], mapToken[refreshTokenRequestParam], claims[userIdClaim]); tokenExpErr != nil {
		fmt.Println(tokenExpErr)
		gwt.settings.ErrResponseFunc(c, http.StatusUnauthorized, tokenExpErr.Error())
		return
	}
	if deleteRefreshErr := gwt.settings.storage.deleteTokens(claims[userIdClaim], claims[refreshUuidClaim],
		claims[accessUuidClaim]); deleteRefreshErr != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusInternalServerError, err.Error())
		return
	}
	accessData, refreshData, er := getTokens(gwt.settings, claims[userIdClaim])
	if er != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusInternalServerError, err.Error())
		return
	}

	gwt.settings.LoginResponseFunc(c, http.StatusOK, accessData.token,
		accessData.expire, refreshData.token, refreshData.expire)
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
	if tokenExpErr := gwt.settings.storage.isAccessExpired(claims[accessUuidClaim], accessToken, claims[userIdClaim]); tokenExpErr != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusUnauthorized, tokenExpErr.Error())
		return
	}
	if deleteRefreshErr := gwt.settings.storage.deleteTokens(claims[userIdClaim], claims[accessUuidClaim],
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
		if tokenExpErr := gwt.settings.storage.isAccessExpired(claims[accessUuidClaim], accessToken, claims[userIdClaim]); tokenExpErr != nil {
			gwt.settings.ErrResponseFunc(c, http.StatusUnauthorized, tokenExpErr.Error())
			return
		}
		user, userErr := gwt.settings.GetUserFunc(claims[userIdClaim])
		if userErr != nil {
			gwt.settings.ErrResponseFunc(c, http.StatusInternalServerError, userErr.Error())
			return
		}
		c.Set(UserKey, user)
		c.Next()
	}
}

func (gwt *Gwt) forceLogoutHandler(c *gin.Context) {
	mapUserId := map[string]string{}
	if err := c.ShouldBind(&mapUserId); err != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusBadRequest, errUserIdIsNotProvided.Error())
		return
	}
	if err := gwt.ForceLogoutUser(mapUserId[userIdRequestParam]); err != nil {
		gwt.settings.ErrResponseFunc(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
