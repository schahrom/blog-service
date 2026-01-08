package lib

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type TokenManager interface {
	NewJWT(userId int64) (string, error)
	Parse(token string) (string, error)
}
type Manager struct {
	signingKey string
	ttl        time.Duration
}

func NewTokenManagerImpl(signingKey string, ttl time.Duration) *Manager {
	return &Manager{signingKey: signingKey, ttl: ttl}
}

func (m *Manager) NewJWT(userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(m.ttl).Unix(),
		Subject:   strconv.Itoa(int(userId)),
	})
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		fmt.Printf("Error parsing token: %v\n", err)
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Printf("Error parsing claims: %v\n", err)
	}
	return claims["sub"].(string), nil
}
