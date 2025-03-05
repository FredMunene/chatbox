package handlers

import (
	"log"
	"net/http"

	"forum/auth"

	"golang.org/x/crypto/bcrypt"
)

func SigninHandler(w http.ResponseWriter, r *http.Request) {
	//  read username and password from request
	// username := r.FormValue("username")
	password := r.FormValue("password")

	dbPassword, _ := bcrypt.GenerateFromPassword([]byte("hashedpassword"), bcrypt.DefaultCost)
	//  check if username and password are correct
	//  retrieve user from sqlite3 database
	//  if user is not found, return an error
	//  if user is found, retrieve hashed password

	//  compare with hashed password

	err := auth.ComparePassword(password, string(dbPassword))
	if err != nil {
		//  return an error
		log.Printf("incorrect password")
		return
	}
	// generate a token

	//  return token to user
}

func Homehandler(w http.ResponseWriter, r *http.Request) {
}
