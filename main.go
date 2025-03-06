package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/backend/auth"
	"forum/backend/handlers"
	"forum/backend/utils"
)

func main() {
	// load environment variables

	utils.LoadEnvVariables(".env")

	http.HandleFunc("/signin", handlers.SigninHandler)
	http.HandleFunc("/auth/google/signin", auth.GoogleSignIn)
	http.HandleFunc("/auth/google/signin/callback", auth.GoogleSignInCallback)
	http.HandleFunc("/auth/github/sign-up", auth.GithubSignUp)
	http.HandleFunc("/auth/github/callback", auth.GithubCallback)
	http.HandleFunc("/home", handlers.HomeHandler)
	fmt.Println("Server running at: http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
