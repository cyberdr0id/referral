package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	IsAdmin bool `json:"isAdmin"`
	jwt.StandardClaims
}

type TokenManager struct {
	key              []byte
	expiryTimeInHour int
}

func NewTokenManager(key string, expiryTime int) *TokenManager {
	return &TokenManager{
		key:              []byte(key),
		expiryTimeInHour: expiryTime,
	}
}

func (t *TokenManager) GenerateToken(userID string, isAdmin bool) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		IsAdmin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(t.expiryTimeInHour)).Unix(),
			Subject:   userID,
		},
	}).SignedString(t.key)
}

func (t *TokenManager) ParseToken(_token string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(_token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(t.key), nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot parse token: %w", err)
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return nil, fmt.Errorf("cannot get claims from token")
	}

	return claims, nil
}
