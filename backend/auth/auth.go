package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func ComparePassword(userPswd string, dbPsswd string) error {
	err := bcrypt.CompareHashAndPassword([]byte(dbPsswd), []byte(userPswd))
	if err != nil {
		return err
	}
	return nil
}
