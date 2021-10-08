package gwt

import "errors"

var (
	// errFailedToCreateAccessToken indicates access Token failed to create, reason unknown
	errFailedToCreateAccessToken = errors.New("failed to create access Token")

	// errFailedToCreateRefreshToken indicates refresh Token failed to create, reason unknown
	errFailedToCreateRefreshToken = errors.New("failed to create refresh Token")

	// errEmptyAccessSecretKey indicates access secret key is empty
	errEmptyAccessSecretKey = errors.New("empty access token secret key")

	// errEmptyGetUserFunc indicates get user func is empty
	errEmptyGetUserFunc = errors.New("empty get user by id func")

	// errEmptyAuthenticator indicates authentication function is empty
	errEmptyAuthenticator = errors.New("empty authentication function")

	// errEmptyRedisConnection indicates redis connection is not provided
	errEmptyRedisConnection = errors.New("empty redis connection")

	// errCannotSaveToken indicates token is failed to save, reason unknown
	errCannotSaveToken = errors.New("failed to save token from storage")

	// errCannotDeleteToken indicates token is failed to delete, reason unknown
	errCannotDeleteToken = errors.New("failed to delete token from storage")

	// errUnknownSigningMethod indicates unknown signing method provided
	errUnknownSigningMethod = errors.New("unknown signing method provided")

	// errInvalidSigningMethod indicates signing method id invalid
	errInvalidSigningMethod = errors.New("invalid signing method")

	// errTokenExpired indicates token has expired
	errTokenExpired = errors.New("token has expired")

	// errTokenInvalid indicates token is not valid
	errTokenInvalid = errors.New("token is not valid")

	// errRefreshTokenIsNotProvided indicates refresh token is not provided
	errRefreshTokenIsNotProvided = errors.New("refresh token is not provided")

	// errNoAuthHeader indicates no auth header is provided
	errNoAuthHeader = errors.New("no auth header provided")

	// errInvalidAuthHeader indicates auth header is not valid
	errInvalidAuthHeader = errors.New("invalid auth header")

	// errUserIdIsNotProvided indicates user id is not provided
	errUserIdIsNotProvided = errors.New("user id is not provided")

	// errNotAuthUser indicates user is not authenticated
	errNotAuthUser = errors.New("user is not authenticated")
)
