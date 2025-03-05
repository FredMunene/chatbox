package main

import (
	"fmt"
	"log"
	"net/http"

	"forum/auth"
	"forum/handlers"
)

func main() {
	// load environment variables

	http.HandleFunc("/signin", handlers.SigninHandler)
	http.HandleFunc("/auth/google/signin", auth.GoogleSignIn)
	http.HandleFunc("/auth/google/signin/callback", auth.GoogleSignInCallback)
	http.HandleFunc("/home", handlers.Homehandler)
	fmt.Println("Server running at: http://localhost:9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
