package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
)

type TokensManager interface {
	NewAccessToken(claims Claims) (string, error)
	GetAccessTokenClaims(accessToken string) (Claims, error)
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

type Claims struct {
	UserID    string
	CompanyID string
}

func (m *tokensManager) NewAccessToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":        claims.UserID,
			"company_id": claims.CompanyID,
			"exp":        time.Now().Add(m.accessTokenTTL).Unix(),
		})

	tokenString, err := token.SignedString([]byte(m.accessTokenKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *tokensManager) GetAccessTokenClaims(accessToken string) (Claims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.accessTokenKey), nil
	})
	if err != nil {
		return Claims{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, fmt.Errorf("error getting claims from access token")
	}

	return Claims{
		UserID:    claims["sub"].(string),
		CompanyID: claims["company_id"].(string),
	}, nil
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
