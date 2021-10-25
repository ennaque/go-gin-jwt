# gin-jwt ![tests](https://github.com/ennaque/go-gin-jwt/workflows/tests/badge.svg) [![codecov](https://codecov.io/gh/ennaque/go-gin-jwt/branch/master/graph/badge.svg?token=WZMWD36EKQ)](https://codecov.io/gh/ennaque/go-gin-jwt) [![codebeat badge](https://codebeat.co/badges/e1ea8bb5-f305-4394-bab2-308efe3f718d)](https://codebeat.co/projects/github-com-ennaque-go-gin-jwt-master)
jwt package for gin go applications

# Usage

Download using [go module](https://blog.golang.org/using-go-modules):

```sh
go get github.com/ennaque/go-gin-jwt@v1.0.2
```

Import it in your code:

```go
import gwt "github.com/ennaque/go-gin-jwt"
import gwtstorage "github.com/ennaque/go-gin-jwt/storage"
```

# Example

```go
package main

import (
	gwt "github.com/ennaque/go-gin-jwt"
	gwtstorage "github.com/ennaque/go-gin-jwt/storage"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	router := gin.Default()

	// GetDBConnectionData() - user func, must return dns string
	// postgres is not required, other drivers can be used here
	db, err := gorm.Open(postgres.Open(GetDBConnectionData()))
	if err != nil {
		panic("db con failed")
	}

	// init gorm storage
	gs, err := gwtstorage.InitGormStorage(db, "jwt1234_")
	if err != nil {
		panic(err)
	}

	// init redis storage
	// GetRedisOptions() - user func, must return *redis.Options
	rs := gwtstorage.InitRedisStorage(redis.NewClient(GetRedisOptions()))

	auth, _ := gwt.Init(gwt.Settings{
		Authenticator: func(c *gin.Context) (string, error) { // required
			// LoginCredentials - your login credentials model, can be differ
			var loginCredentials LoginCredentials
			if err := c.ShouldBind(&loginCredentials); err != nil {
				return "", errors.New("bad request")
			}
			// GetUserByCredentials - user func, must return user model
			user, err := GetUserByCredentials(&loginCredentials)
			if err != nil {
				return "", errors.New("unauthorized")
			}
			return user.GetId(), nil
		},
		AccessSecretKey:  []byte("access_super_secret"), // required
		RefreshSecretKey: []byte("refresh_super_secret"), // optional, default - AccessSecretKey
		Storage:          gs, // required, use gorm or redis storage
		// Storage: rs,
		GetUserFunc: func(userId string) (interface{}, error) { // required
			return GetUserById(userId)
		},
		AccessLifetime:  time.Minute * 15, // optional, default - time.Minute * 15
		RefreshLifetime: time.Hour * 24, // optional, default - time.Hour * 24
		SigningMethod:   "HS256", // optional, default - HS256
		AuthHeadName:    "Bearer", // optional, default - Bearer
	})

	a := router.Group("auth")
	{
		a.POST("/logout", auth.Handler.GetLogoutHandler())
		a.POST("/login", auth.Handler.GetLoginHandler())
		a.POST("/refresh", auth.Handler.GetRefreshHandler())
		a.POST("/force-logout", auth.Handler.GetForceLogoutHandler())
	}

	router.Group("/api").Use(auth.Middleware.GetAuthMiddleware()).GET("/get-user-id", func(c *gin.Context) {
		user, _ := c.Get("user")
		c.JSON(http.StatusOK, gin.H{
			"userId": user.(*models.User).ID,
		})
	})

	err := router.Run(":8000")
	if err != nil {
		panic("err")
	}
}
```

## Get tokens

```sh
curl -X POST -d "username=<user_name>&password=<password>" http://localhost:8000/auth/login
```
username and password params may differ depending on your login credentials

Response `200 OK`:
```sh
{
    "access_expire": "1633653988",
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjlkODFmNjRkLWY0ZWYtNDA2NC04YTY3LTRjNjMzY2MxNjExOCIsImV4cCI6MTYzMzY1Mzk4OCwicmVmcmVzaF91dWlkIjoiOTU3NWU5ZDEtNWFjOS00YmIzLTkwOGItODA3MmJkNDdmOTM2IiwidXNlcl9pZCI6IjI5In0.0CfHPjkVFiQixa4SdE5EUhu23imNri02QMFsDDXJHzg",
    "refresh_expire": "1633739788",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjlkODFmNjRkLWY0ZWYtNDA2NC04YTY3LTRjNjMzY2MxNjExOCIsImV4cCI6MTYzMzczOTc4OCwicmVmcmVzaF91dWlkIjoiOTU3NWU5ZDEtNWFjOS00YmIzLTkwOGItODA3MmJkNDdmOTM2IiwidXNlcl9pZCI6IjI5In0.UvPTvVaNkAgFVTrAEoaUK1n4iIYFGh1yNqPzzNbtUUM"
}
```

## Refresh token

```sh
curl -X POST -d "refresh_token=<refresh_token>" http://localhost:8000/auth/refresh
```

Response `200 OK`:
```sh
{
    "access_expire":"1633659261",
    "access_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImJiNjBhYzlmLTQ4ZGEtNDlhZC04NTM1LTU5MTJhY2MwZDIwNyIsImV4cCI6MTYzMzY1OTI2MSwicmVmcmVza
F91dWlkIjoiNDkxMWYxZjUtYjk5Ni00ZTEwLWE4NGEtNDg3NGVmNjMzZDc4IiwidXNlcl9pZCI6IjI5In0.tupNFRnANQmOScjWzlnWXzncX0Kxs7M40rsbFs0Vpg-70Ucc7R7vX2e7uAFf1fiAMODfGS5d3PRK3Nwk4RoPzg",
    "refresh_expire":"1633831941",
    "refresh_token":"eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImJiNjBhYzlmLTQ4ZGEtNDlhZC04NTM1LTU5MTJhY2MwZDIwNyIsImV4cCI6MTYzMzgzMTk0MSwicmVmcmVz
aF91dWlkIjoiNDkxMWYxZjUtYjk5Ni00ZTEwLWE4NGEtNDg3NGVmNjMzZDc4IiwidXNlcl9pZCI6IjI5In0.lj2nS6-M4GT-T9PHj9ijNY4g6h5hyP0xdVTHCw1M-07aL4zp7HpFrXFrT-V6RWpofaGvM79o64f8WECEqRPjig"
}
```

## Logout

```sh
curl -X POST -H "Authorization: Bearer <access_token>" http://localhost:8000/auth/logout
```

Response `200 OK`:
```sh
{}
```

## Force logout user

This endpoint should be used only by authorized user.

```sh
curl -X POST -H "Authorization: Bearer <access_token>" -d "user_id=<user_id_to_logout>" http://localhost:8000/auth/force-logout
```
Response `200 OK`:
```sh
{}
```
Additionaly there is a public method ```gwt.Service.ForceLogoutUser(userId)```
