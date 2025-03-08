package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error while hashing password: %w", err)
	}
	return bytes, nil
}

func HashMatchesPassword(hash, pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}
