// Package hash contains function for password hashing and comparing.
package hash

import "golang.org/x/crypto/bcrypt"

// HashPassword transform password to hash-string.
func HashPassword(password string) (string, error) {
	cost := 14

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

// CheckPasswordHash comparing password with his hash.
func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
