package auth

import (
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ComparePassword(passwordHash string, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plain))
}
