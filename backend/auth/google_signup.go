package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	utils "forum/backend/util"
)

func GoogleSignUp(w http.ResponseWriter, r *http.Request) {
	// generate state with cookie?

	state := generateStateCookie(w, "signup")

	// set the Google OAuth 2.0 authorization URL
	redirectUrl := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid email profile&state=%s&prompt=select_account&access_type=offline",
		GoogleAuthUrl, utils.GoogleClientID, url.QueryEscape(BaseUrl+"/auth/google/signin/callback"), state)

	//  set CORS specifying domain
	w.Header().Set("Access-Control-Allow-Origin", BaseUrl)

	//  redirect user
	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the callback from the Google OAuth 2.0 server after the user has granted the necessary permissions.
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Validate the state to prevent CSRF attacks
	if err := validateState(r); err != nil {
		log.Printf("State validation failed: %v", err)
		http.Redirect(w, r, "/sign-in?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	// Get the authorization code from the query parameter
	code := r.URL.Query().Get("code")
	token, err := exchangeGithubCodeForToken(code)
	if err != nil {
		log.Printf("Token exchange failed: %v\n", err)
		http.Redirect(w, r, "/sign-in?error=token_exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	// Get the user information from the Google UserInfo endpoint
	user, err := getGoogleUserDetails(token)
	if err != nil {
		log.Printf("Failed to get user info: %v\n", err)
		http.Redirect(w, r, "/sign-in?error=user_info_failed", http.StatusTemporaryRedirect)
		return
	}

	// Attempt to create or authenticate the user
	if handleUserAuth(w, user.Email, user.Name) {
		log.Printf("User authentication successful")
		// http.Redirect(w, r, "/sign-in", http.StatusTemporaryRedirect)
		http.Redirect(w, r, "/sign-in?status=success", http.StatusTemporaryRedirect)
	} else {
		log.Printf("User authentication failed")
		http.Redirect(w, r, "/sign-in?error=auth_failed", http.StatusTemporaryRedirect)
	}
}
