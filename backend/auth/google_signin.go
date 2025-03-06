package auth

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"forum/backend/utils"
)

func GoogleSignIn(w http.ResponseWriter, r *http.Request) {
	// generate state with cookie?

	state := generateStateCookie(w, "signin")

	// set the Google OAuth 2.0 authorization URL
	redirectUrl := fmt.Sprintf(
		"%s?client_id=%s&redirect_uri=%s&response_type=code&scope=openid email profile&state=%s&prompt=select_account&access_type=offline",
		GoogleAuthUrl, utils.GoogleClientID, url.QueryEscape(BaseUrl+"/auth/google/signin/callback"), state)

	//  set CORS specifying domain
	w.Header().Set("Access-Control-Allow-Origin", BaseUrl)

	//  redirect user
	http.Redirect(w, r, redirectUrl, http.StatusTemporaryRedirect)
}

func GoogleSignInCallback(w http.ResponseWriter, r *http.Request) {
	//  validate state

	if err := validateState(r); err != nil {
		log.Println("Invalid state")
		http.Redirect(w, r, "/sign-in?error=invalid_state", http.StatusTemporaryRedirect)
		return
	}

	//  get authorization code
	code := r.URL.Query().Get("code")

	// use code to get access token
	token, err := swapGoogleCodeForToken(code)
	if err != nil {
		log.Printf("Google token exchange failed: %v", err)
		http.Redirect(w, r, "/sign-in?error=token_exchange_failed", http.StatusTemporaryRedirect)
		return
	}

	//  retrieve user info from Google
	user, err := getGoogleUserDetails(token)
	if err != nil {
		log.Printf("Failed to get user details from Google: %v", err)
		http.Redirect(w, r, "/sign-in?error=retrieving_user_info_failed", http.StatusTemporaryRedirect)
	}

	println("User email:"+user.Email, user.Sub, "Username:"+user.Name)

	http.Redirect(w, r, "/home?status=success", http.StatusSeeOther)
}
