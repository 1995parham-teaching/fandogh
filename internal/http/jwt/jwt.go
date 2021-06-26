package jwt

import (
	"fmt"
	"time"

	"github.com/1995parham/fandogh/internal/model"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type Config struct {
	AccessTokenSecret string
}

// NewAccessToken creates new access token for given user.
func NewAccessToken(u model.User, cfg Config) (string, error) {
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
	encodedToken, err := token.SignedString([]byte(cfg.AccessTokenSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign a token: %w", err)
	}

	return encodedToken, nil
}
