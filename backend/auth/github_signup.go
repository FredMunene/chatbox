package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"forum/backend/utils"
)

const (
	GithubAuthUrl   = "https://github.com/login/oauth/authorize"
	GithubTokenUrl  = "https://github.com/login/oauth/access_token"
	GithubUserUrl   = "https://api.github.com/user"
	GithubEmailsUrl = "https://api.github.com/user/emails"
	BaseUrl = "http://localhost:9000"
)

type GithubUser struct {
	Login, Email string
}

func GithubSignUp(w http.ResponseWriter, r *http.Request) {
	state := generateStateCookie(w, "sign-up")

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

	// log.Printf("Cookie state: %s", cookie.Value)
	// log.Printf("URL state: %s", r.URL.Query().Get("state"))

	// retrieve state from cookie
	stateParts := strings.Split(cookie.Value, ":")
	if len(stateParts) != 2 {
		log.Printf("Invalid state")
		http.Redirect(w, r, "/sign-in?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	_, flowType := stateParts[0], stateParts[1]

	// // check if state from cookie matches state from query

	// urlState := strings.Split(r.URL.Query().Get("state"), ":")

	// if urlState[0] != state {
	// 	log.Println("State don't match")
	// 	http.Redirect(w, r, "/sign-in?error=invalid_state", http.StatusTemporaryRedirect)
	// 	return
	// }

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
	case "sign-up":
		// handle user sign up
		http.Redirect(w, r, "/home?status=success", http.StatusSeeOther)

	case "signin":
		// handle user sign in
		http.Redirect(w, r, "/home?status=success", http.StatusSeeOther)

	default:
		http.Redirect(w, r, "/sign-in?error=invalid_flow", http.StatusTemporaryRedirect)
	}
}
