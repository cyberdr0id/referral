package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Claims presents a type for storing necessary user information.
type Claims struct {
	IsAdmin bool `json:"isAdmin"`
	jwt.StandardClaims
}

// TokenManager presents a type for token management, it's contains key for sign token
// and expiration time.
type TokenManager struct {
	key              []byte
	expiryTimeInHour int
}

// NewTokenManager creates a new instance of TokenManager.
func NewTokenManager(key, expiryTime string) (*TokenManager, error) {
	et, err := strconv.Atoi(expiryTime)
	if err != nil {
		return &TokenManager{}, fmt.Errorf("cannot convert expity time of JWT: %w", err)
	}

	return &TokenManager{
		key:              []byte(key),
		expiryTimeInHour: et,
	}, nil
}

// GenerateToken generates JWT token
func (t *TokenManager) GenerateToken(userID string, isAdmin bool) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		IsAdmin: isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(t.expiryTimeInHour)).Unix(),
			Subject:   userID,
		},
	}).SignedString(t.key)
}

// ParseToken gets the user claims from JWT token.
func (t *TokenManager) ParseToken(_token string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(_token, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return t.key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("cannot get claims from token")
	}

	return claims, nil
}
