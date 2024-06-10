package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

type TokensManager interface {
	NewAccessToken(id string) (string, error)
	NewRefreshToken() (string, error)
}

type tokensManager struct {
	accessTokenKey string
	accessTokenTTL time.Duration
}

func NewTokensManager(accessTokenKey string, accessTokenTTL time.Duration) TokensManager {
	return &tokensManager{
		accessTokenKey: accessTokenKey,
		accessTokenTTL: accessTokenTTL,
	}
}

func (m *tokensManager) NewAccessToken(id string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": id,
			"exp": time.Now().Add(m.accessTokenTTL).Unix(),
		})

	tokenString, err := token.SignedString([]byte(m.accessTokenKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *tokensManager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
