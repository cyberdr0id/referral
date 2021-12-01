package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var key = []byte("Str0ngP@$$w0rd?##")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateToken(username string) (string, error) {
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 5).Unix(),
		},
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(_token string) (string, error) {
	token, err := jwt.ParseWithClaims(_token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(key), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok {
		return "", fmt.Errorf("cannot get claims from token")
	}

	return claims.Subject, nil
}
