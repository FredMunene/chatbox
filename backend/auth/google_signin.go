package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"forum/backend/handlers"
	"forum/backend/repositories"
	"forum/backend/util"
	utils "forum/backend/util"
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

	switch flowType {
	case "signup":
		// handle user sign up
		if handleUserAuth(w, user.Email, user.Name) {
			log.Printf("User Authentication successful")
			http.Redirect(w, r, "/sign-up?status=success", http.StatusTemporaryRedirect)
		} else {

			log.Printf("User signup failed")
			http.Redirect(w, r, "/sign-up?status=auth_failed", http.StatusSeeOther)

		}

	case "signin":
		// Check if the user exists in the database
		var userID int
		err = util.Database.QueryRow("SELECT id FROM tblUsers WHERE email = ?", user.Email).Scan(&userID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Redirect(w, r, "/sign-up?error=no_account", http.StatusTemporaryRedirect)
				return
			}
			log.Printf("Database error: %v", err)
			http.Redirect(w, r, "/sign-in?error=database_error", http.StatusTemporaryRedirect)
			return
		}

		sessionToken := handlers.CreateSession()

		// Delete any existing sessions for this user
		if userID != 0 {
			handlers.DeleteSession(userID)
		}
		err = repositories.DeleteSessionByUser(userID)
		if err != nil {
			log.Printf("Failed to delete session token: %v", err)
			http.Redirect(w, r, "/sign-in?error=session_error", http.StatusTemporaryRedirect)
			return
		}

		// Enable CORS
		handlers.EnableCors(w)

		// Set session cookie and data
		handlers.SetSessionCookie(w, sessionToken)
		handlers.SetSessionData(sessionToken, "userId", userID)
		handlers.SetSessionData(sessionToken, "userEmail", user.Email)

		// Store session with expiry time
		expiryTime := time.Now().Add(24 * time.Hour)
		err = repositories.StoreSession(userID, sessionToken, expiryTime)
		if err != nil {
			log.Printf("Failed to store session token: %v", err)
			http.Redirect(w, r, "/sign-in?error=session_error", http.StatusTemporaryRedirect)
			return
		}

	}

	http.Redirect(w, r, "/home?status=success", http.StatusSeeOther)
}
