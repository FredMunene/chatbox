package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"forum/utils"

	"github.com/google/uuid"
)

type GoogleUser struct {
	Sub, Name, Email string
}

const (
	GoogleAuthUrl     = "https://accounts.google.com/o/oauth2/v2/auth"
	GoogleTokenUrl    = "https://oauth2.googleapis.com/token"
	GoogleUserInfoUrl = "https://www.googleapis.com/oauth2/v3/userinfo"
)

func generateStateCookie(w http.ResponseWriter) string {
	//  Generate a random UUID
	state := uuid.NewString()

	//  Set cookie

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		Secure:   false, // set to true for production for HTTPS-only
		MaxAge:   3600,  // 1 hr
		SameSite: http.SameSiteLaxMode,
	})

	return state
}

func validateState(r *http.Request) error {
	// query url for state
	state := r.URL.Query().Get("state")
	//  retrieve cookie
	cookie, err := r.Cookie("oauth_state")
	if err != nil {
		log.Printf("Cookie error: %v", err)
		return err
	}

	//  check states match
	if cookie.Value != state {
		log.Printf("State don't match. Cookie value:%s, State:%s", cookie.Value, state)
		return errors.New("invalid state")
	}

	return nil
}

func swapGoogleCodeForToken(code string) (string, error) {
	data := url.Values{
		"code":          {code},
		"client_id":     {utils.GoogleClientID},
		"client_secret": {utils.GoogleClientSecret},
		"redirect_uri":  {"https://localhost:9000/auth/google/signin/callback"},
		"grant_type":    {"authorization_code"},
	}

	println("here")

	// make a POST request to Google's token url
	resp, err := http.PostForm(GoogleTokenUrl, data)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// check HTTP status code

	fmt.Println(resp.StatusCode)

	// if resp.StatusCode != http.StatusOK {
	// 	return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	// }

	var token struct {
		AccessToken string `json:"access_token"`
	}

	json.NewDecoder(resp.Body).Decode(&token)
	return token.AccessToken, nil
}

func getGoogleUserDetails(token string) (*GoogleUser, error) {
	var user GoogleUser

	request, err := http.NewRequest("GET", GoogleUserInfoUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// decode resp into GoogleUset
	json.NewDecoder(resp.Body).Decode(&user)

	return &user, nil
}
