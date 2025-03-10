package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"forum/backend/repositories"
	"forum/backend/utils"

	"forum/backend/handlers"
	"forum/backend/util"
)

const (
	GithubAuthUrl   = "https://github.com/login/oauth/authorize"
	GithubTokenUrl  = "https://github.com/login/oauth/access_token"
	GithubUserUrl   = "https://api.github.com/user"
	GithubEmailsUrl = "https://api.github.com/user/emails"
	BaseUrl         = "http://localhost:9000"
)

type GithubUser struct {
	Login, Email string
}

func GithubSignUp(w http.ResponseWriter, r *http.Request) {
	state := generateStateCookie(w, "signup")

	params := url.Values{
		"client_id":    {utils.GithubClientID},
		"redirect_uri": {BaseUrl + "/auth/github/callback"},
		"scope":        {"read:user user:email"},
		"state":        {state},
		"prompt":       {"consent"},
	}

	redirectUrl := fmt.Sprintf("%s?%s", GithubAuthUrl, params.Encode())
	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
}

func GithubCallback(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		log.Printf("Cookie error: %v", err)
		http.Redirect(w, r, "/sign-in?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}
	// retrieve state from cookie
	stateParts := strings.Split(cookie.Value, ":")
	if len(stateParts) != 2 {
		log.Printf("Invalid state")
		http.Redirect(w, r, "/sign-in?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	_, flowType := stateParts[0], stateParts[1]

	if err := validateState(r); err != nil {
		log.Println("State don't match")
		http.Redirect(w, r, "/sign-in?error=invalid_state", http.StatusTemporaryRedirect)
		return

	}

	code := r.URL.Query().Get("code")
	token, err := exchangeGithubCodeForToken(code)
	if err != nil {
		log.Printf("Github token exchange failed: %v", err)
		http.Redirect(w, r, "/sign-in?error=token_exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	user, err := getGithubUserDetails(token)
	if err != nil {
		log.Printf("Retrieving user details failed: %v", err)
		return
	}

	if user != nil {
		fmt.Println(user.Email)
	}

	switch flowType {
	case "signup":
		// handle user sign up
		if handleUserAuth(w, user.Email, user.Login) {
			log.Printf("User signup successful")
			http.Redirect(w, r, "/sign-up?status=success", http.StatusTemporaryRedirect)
		} else {

			log.Printf("User signup failed")
			http.Redirect(w, r, "/sign-up?status=auth_failed", http.StatusSeeOther)

		}

	case "signin":
		// handle user sign in
		var userID int

		err := util.Database.QueryRow("SELECT id FROM Users WHERE email = ?", user.Email).scan(&userID)
		if err != nil {
			http.Redirect(w, r, "/sign-inp?error=no_account", http.StatusTemporaryRedirect)
		}

		sessionToken := handlers.CreateSession()
		if userID != 0 {
			handlers.DeleteSession(userID)
		}

		handlers.EnableCors(w)
		handlers.SetSessionCookie(w,sessionToken)
		handlers.SetSessionData()
		handlers.SetSessionData()

		expiryTime := time.Now().Add(24* time.Hour)

		err = repositories.StoreSession(userID, sessionToken, expiryTime)
		if err != nil {
			log.Printf("Failed to store session token:%v", err)
			http.Redirect(w,r,"/sign-in?error=session_error", http.StatusTemporaryRedirect)
			return
		}

		http.Redirect(w, r, "/home?status=success", http.StatusSeeOther)

	default:
		http.Redirect(w, r, "/sign-in?error=invalid_flow", http.StatusTemporaryRedirect)
	}
}
