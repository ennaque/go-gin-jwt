package gwt

import "errors"

var (
	// ErrFailedToCreateAccessToken indicates access Token failed to create, reason unknown
	ErrFailedToCreateAccessToken = errors.New("failed to create access Token")

	// ErrFailedToCreateRefreshToken indicates refresh Token failed to create, reason unknown
	ErrFailedToCreateRefreshToken = errors.New("failed to create refresh Token")

	// ErrEmptyAccessSecretKey indicates access secret key is empty
	ErrEmptyAccessSecretKey = errors.New("empty access token secret key")

	// ErrEmptyAuthenticator indicates authentication function is empty
	ErrEmptyAuthenticator = errors.New("empty authentication function")

	// ErrCannotSaveToken indicates token is failed to save, reason unknown
	ErrCannotSaveToken = errors.New("failed to save token from storage")

	// ErrCannotDeleteToken indicates token is failed to delete, reason unknown
	ErrCannotDeleteToken = errors.New("failed to delete token from storage")

	// ErrUnknownSigningMethod indicates unknown signing method provided
	ErrUnknownSigningMethod = errors.New("unknown signing method provided")

	// ErrInvalidSigningMethod indicates signing method id invalid
	ErrInvalidSigningMethod = errors.New("invalid signing method")

	// ErrTokenExpired indicates token has expired
	ErrTokenExpired = errors.New("token has expired")

	// ErrTokenInvalid indicates token is not valid
	ErrTokenInvalid = errors.New("token is not valid")

	// ErrRefreshTokenIsNotProvided indicates refresh token is not provided
	ErrRefreshTokenIsNotProvided = errors.New("refresh token is not provided")

	// ErrNoAuthHeader indicates no auth header is provided
	ErrNoAuthHeader = errors.New("no auth header provided")

	// ErrInvalidAuthHeader indicates auth header is not valid
	ErrInvalidAuthHeader = errors.New("invalid auth header")
)
