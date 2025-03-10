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
	// r.HandleFunc("/signin", handlers.SigninHandler)
	r.HandleFunc("/home", handlers.IndexHandler)
	r.HandleFunc("/", handlers.HomeHandler)
	r.HandleFunc("/signin", handlers.LoginHandler)
	r.HandleFunc("/signup", handlers.SignupHandler)
	r.HandleFunc("/logout", handlers.LogoutHandler)

	r.HandleFunc("/auth/google/signin", auth.GoogleSignIn)
	r.HandleFunc("/auth/google/signup", auth.GoogleSignUp)
	r.HandleFunc("/auth/google/signin/callback", auth.GoogleSignInCallback)

	r.HandleFunc("/auth/github/signup", auth.GithubSignUp)
	r.HandleFunc("/auth/github/signin", auth.GitHubSignIn)
	r.HandleFunc("/auth/github/callback", auth.GithubCallback)

	// r.HandleFunc("")

	return r
}
