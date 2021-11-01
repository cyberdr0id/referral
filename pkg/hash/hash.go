// Package hash contains function for password hashing and comparing.
package hash

import "golang.org/x/crypto/bcrypt"

// HashPassword transform password to hash-string.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash comparing password with his hash.
func CheckPassowrdHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
