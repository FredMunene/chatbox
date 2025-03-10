package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func ComparePassword(userPassword string, dbPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(userPassword))
	if err != nil {
		return err
	}
	return nil
}
