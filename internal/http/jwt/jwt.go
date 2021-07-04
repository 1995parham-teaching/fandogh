package jwt

import (
	"fmt"
	"time"

	"github.com/1995parham/fandogh/internal/http/common"
	"github.com/1995parham/fandogh/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Config struct {
	AccessTokenSecret string
}

type JWT struct {
	Config
}

func (j JWT) Middleware() echo.MiddlewareFunc {
	// nolint: exhaustivestruct
	return middleware.JWTWithConfig(middleware.JWTConfig{
		ContextKey:    common.UserContextKey,
		SigningKey:    []byte(j.Config.AccessTokenSecret),
		SigningMethod: jwt.SigningMethodHS256.Name,
		Claims:        &jwt.StandardClaims{},
		TokenLookup:   "header:Authorization",
	})
}

// NewAccessToken creates new access token for given user.
func (j JWT) NewAccessToken(u model.User) (string, error) {
	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Audience:  "user",
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		Id:        uuid.New().String(),
		IssuedAt:  time.Now().Unix(),
		Issuer:    "fandogh",
		NotBefore: time.Now().Unix(),
		Subject:   u.Email,
	})

	// generate encoded token and send it as response
	encodedToken, err := token.SignedString([]byte(j.AccessTokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign a token: %w", err)
	}

	return encodedToken, nil
}
