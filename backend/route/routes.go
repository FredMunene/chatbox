package route

import (
	"net/http"

	"forum/backend/auth"
	"forum/backend/handlers"
)

func InitRoutes() *http.ServeMux {
	r := http.NewServeMux()

	fs := http.FileServer(http.Dir("./frontend"))
	r.Handle("/frontend/", http.StripPrefix("/frontend/", fs))

	uploadsFs := http.FileServer(http.Dir("./uploads"))
	r.Handle("/uploads/", http.StripPrefix("/uploads/", uploadsFs))

	//  routes
	r.HandleFunc("/signin", handlers.SigninHandler)
	r.HandleFunc("/home", handlers.HomeHandler)

	r.HandleFunc("/auth/google/signin", auth.GoogleSignIn)
	r.HandleFunc("/auth/google/signin/callback", auth.GoogleSignInCallback)


	r.HandleFunc("/auth/github/signup", auth.GithubSignUp)
	r.HandleFunc("/auth/github/callback", auth.GithubCallback)

	return r
	
}
