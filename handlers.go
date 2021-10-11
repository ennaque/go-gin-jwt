package gwt

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	settings *Settings
}

func (handler *Handler) GetLoginHandler() func(c *gin.Context) {
	return handler.loginHandler
}
func (handler *Handler) GetRefreshHandler() func(c *gin.Context) {
	return handler.refreshHandler
}
func (handler *Handler) GetLogoutHandler() func(c *gin.Context) {
	return handler.logoutHandler
}
func (handler *Handler) GetForceLogoutHandler() func(c *gin.Context) {
	return handler.forceLogoutHandler
}

func (handler *Handler) loginHandler(c *gin.Context) {
	userId, err := handler.settings.Authenticator(c)
	if err != nil {
		handler.settings.ErrResponseFunc(c, http.StatusUnauthorized, err.Error())
		return
	}
	accessData, refreshData, er := getTokens(handler.settings, userId)
	if er != nil {
		handler.settings.ErrResponseFunc(c, http.StatusInternalServerError, err.Error())
		return
	}

	handler.settings.LoginResponseFunc(c, http.StatusOK, accessData.token,
		accessData.expire, refreshData.token, refreshData.expire)
}
func (handler *Handler) refreshHandler(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBind(&mapToken); err != nil {
		handler.settings.ErrResponseFunc(c, http.StatusBadRequest, ErrRefreshTokenIsNotProvided.Error())
		return
	}
	parsedToken, err := parseToken(mapToken[refreshTokenRequestParam], handler.settings.RefreshSecretKey, handler.settings.SigningMethod)
	if err != nil {
		handler.settings.ErrResponseFunc(c, http.StatusBadRequest, err.Error())
		return
	}
	claims, getClaimsErr := getClaims(parsedToken, []string{refreshUuidClaim, accessUuidClaim, userIdClaim, expiredClaim})
	if err != nil {
		handler.settings.ErrResponseFunc(c, http.StatusBadRequest, getClaimsErr.Error())
		return
	}
	if tokenExpErr := handler.settings.Storage.HasRefreshToken(claims[refreshUuidClaim], mapToken[refreshTokenRequestParam], claims[userIdClaim]); tokenExpErr != nil {
		handler.settings.ErrResponseFunc(c, http.StatusUnauthorized, tokenExpErr.Error())
		return
	}
	if expErr := isExpired(claims[expiredClaim]); expErr != nil {
		handler.settings.ErrResponseFunc(c, http.StatusUnauthorized, expErr.Error())
		return
	}
	if deleteRefreshErr := handler.settings.Storage.DeleteTokens(claims[userIdClaim], claims[refreshUuidClaim],
		claims[accessUuidClaim]); deleteRefreshErr != nil {
		handler.settings.ErrResponseFunc(c, http.StatusInternalServerError, err.Error())
		return
	}
	accessData, refreshData, er := getTokens(handler.settings, claims[userIdClaim])
	if er != nil {
		handler.settings.ErrResponseFunc(c, http.StatusInternalServerError, err.Error())
		return
	}

	handler.settings.LoginResponseFunc(c, http.StatusOK, accessData.token,
		accessData.expire, refreshData.token, refreshData.expire)
}
func (handler *Handler) logoutHandler(c *gin.Context) {
	accessToken, getErr := getHeaderToken(c.Request.Header.Get(authHeader), handler.settings.AuthHeadName)
	if getErr != nil {
		handler.settings.ErrResponseFunc(c, http.StatusUnauthorized, getErr.Error())
		return
	}
	parsedToken, parseErr := parseToken(accessToken, handler.settings.AccessSecretKey, handler.settings.SigningMethod)
	if parseErr != nil {
		handler.settings.ErrResponseFunc(c, http.StatusBadRequest, parseErr.Error())
		return
	}
	claims, getClaimsErr := getClaims(parsedToken, []string{refreshUuidClaim, accessUuidClaim, userIdClaim, expiredClaim})
	if getClaimsErr != nil {
		handler.settings.ErrResponseFunc(c, http.StatusBadRequest, getClaimsErr.Error())
		return
	}
	if tokenExpErr := handler.settings.Storage.HasAccessToken(claims[accessUuidClaim], accessToken, claims[userIdClaim]); tokenExpErr != nil {
		handler.settings.ErrResponseFunc(c, http.StatusUnauthorized, tokenExpErr.Error())
		return
	}
	if expErr := isExpired(claims[expiredClaim]); expErr != nil {
		handler.settings.ErrResponseFunc(c, http.StatusUnauthorized, expErr.Error())
		return
	}
	if deleteRefreshErr := handler.settings.Storage.DeleteTokens(claims[userIdClaim], claims[accessUuidClaim],
		claims[refreshUuidClaim]); deleteRefreshErr != nil {
		handler.settings.ErrResponseFunc(c, http.StatusInternalServerError, deleteRefreshErr.Error())
		return
	}

	handler.settings.LogoutResponseFunc(c, http.StatusOK)
}

func (handler *Handler) forceLogoutHandler(c *gin.Context) {
	mapUserId := map[string]string{}
	if err := c.ShouldBind(&mapUserId); err != nil {
		handler.settings.ErrResponseFunc(c, http.StatusBadRequest, ErrUserIdIsNotProvided.Error())
		return
	}
	if err := handler.settings.Storage.DeleteAllTokens(mapUserId[userIdRequestParam]); err != nil {
		handler.settings.ErrResponseFunc(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
