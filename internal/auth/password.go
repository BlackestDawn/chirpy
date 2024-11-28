package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	if len(password) > 72 {
		password = password[:72]
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
	if len(password) > 72 {
		password = password[:72]
	}
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
