package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTManager interface {
	GenerateToken(id string) string
}

type jwtManager struct {
	key string
	ttl time.Duration
}

func NewJWTManager(key string, ttl time.Duration) JWTManager {
	return &jwtManager{
		key: key,
		ttl: ttl,
	}
}

func (m *jwtManager) GenerateToken(id string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": id,
			"exp": time.Now().Add(m.ttl).Unix(),
		})

	tokenString, err := token.SignedString([]byte(m.key))
	if err != nil {
		panic("failed to sign jwt: " + err.Error())
	}

	return tokenString
}
