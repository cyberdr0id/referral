package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"

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

type jwtConfig struct {
	Key        string `envconfig:"JWT_KEY"`
	ExpiryTime string `envconfig:"JWT_EXPIRY_TIME"`
}

// NewTokenManager creates a new instance of TokenManager.
func NewTokenManager() (*TokenManager, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("unable to load JWT token config: %w", err)
	}

	et, err := strconv.Atoi(config.ExpiryTime)
	if err != nil {
		return &TokenManager{}, fmt.Errorf("cannot convert expity time of JWT: %w", err)
	}

	return &TokenManager{
		key:              []byte(config.Key),
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

func loadConfig() (*jwtConfig, error) {
	var c jwtConfig

	err := envconfig.Process("jwt", &c)
	if err != nil {
		return nil, fmt.Errorf("unable load JWT config: %w", err)
	}

	return &c, nil
}
