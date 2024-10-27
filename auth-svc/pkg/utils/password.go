package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) string {
	hashBytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashBytes)
}

func CompareHashAndPassword(hashedPass string, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), password)
	return err == nil
}
