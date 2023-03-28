package jwt

import (
	"fmt"
	"time"

	"github.com/1995parham-teaching/fandogh/internal/http/common"
	"github.com/1995parham-teaching/fandogh/internal/model"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type Config struct {
	AccessTokenSecret string
}

type JWT struct {
	Config
}

func (j JWT) Middleware() echo.MiddlewareFunc {
	// nolint: exhaustruct
	return echojwt.WithConfig(echojwt.Config{
		ContextKey:    common.UserContextKey,
		SigningKey:    []byte(j.AccessTokenSecret),
		SigningMethod: jwt.SigningMethodHS256.Name,
		NewClaimsFunc: func(_ echo.Context) jwt.Claims { return new(jwt.RegisteredClaims) },
		TokenLookup:   "header:Authorization",
	})
}

// NewAccessToken creates new access token for given user.
func (j JWT) NewAccessToken(u model.User) (string, error) {
	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Audience:  jwt.ClaimStrings{"user"},
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		ID:        uuid.New().String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "fandogh",
		NotBefore: jwt.NewNumericDate(time.Now()),
		Subject:   u.Email,
	})

	// generate encoded token and send it as response
	encodedToken, err := token.SignedString([]byte(j.AccessTokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign a token: %w", err)
	}

	return encodedToken, nil
}
