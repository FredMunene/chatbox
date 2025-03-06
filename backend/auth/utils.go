package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"forum/backend/utils"

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

func generateStateCookie(w http.ResponseWriter, flowType string) string {
	//  Generate a random UUID
	state := uuid.NewString()

	//  Set cookie

	stateData := fmt.Sprintf("%s:%s", state, flowType)

	log.Printf("Generated state before assigning to cookie: %s", state)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    stateData,
		Path:     "/",
		Domain:   "",
		HttpOnly: true,
		Secure:   false, // set to true for production for HTTPS-only
		MaxAge:   3600,  // 1 hr
		SameSite: http.SameSiteLaxMode,
	})

	log.Printf("Generated state: %s", stateData)

	return stateData
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

	stateParts := strings.Split(cookie.Value, ":")
	if len(stateParts) != 2 {
		log.Printf("Invalid state cookie value: %s", cookie.Value)
		return errors.New("invalid state cookie value")
	}

	//  check states match
	if stateParts[0] != state {
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

func exchangeGithubCodeForToken(code string) (string, error) {
	data := url.Values{
		"client_id":     {utils.GithubClientID},
		"client_secret": {utils.GithubClientSecret},
		"code":          {code},
		"redirect_uri":  {RedirectBaseUrl + "/auth/github/callback"}, 
	}

	resp, err := http.PostForm(GithubTokenUrl, data)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		return "", err
	}

	defer resp.Body.Close()

	var token struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}

	json.NewDecoder(resp.Body).Decode(&token)
	return token.AccessToken, nil
}

func getGithubUserDetails(token string) (*GithubUser, error) {
	//  make a GET request to githubuserurl
	req, err := http.NewRequest("GET", GithubUserUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept","application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user GithubUser

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	if user.Email == "" {
		req, err := http.NewRequest("GET", GithubEmailsUrl, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "token "+token)
		req.Header.Set("Accept","application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		var emails []struct {
			Email string `json:"email"`
			Primary bool `json:"primary"`
			Verified bool `json:"verified"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
			return nil, err
		}

		for _, email := range emails {
			if email.Primary && email.Verified{
				user.Email = email.Email
				break
			}
		}

		}

		return &user, nil
}
