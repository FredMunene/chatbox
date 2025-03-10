package auth

import (
	"fmt"
	"net/http"
	"net/url"

	"forum/backend/util"
)

func GitHubSignIn(w http.ResponseWriter, r *http.Request) {
	state := generateStateCookie(w, "signin")

	params := url.Values{
		"client_id":    {util.GithubClientID},
		"redirect_uri": {BaseUrl + "/auth/github/callback"},
		"scope":        {"user:email"},
		"state":        {state},
		"prompt":       {"consent"},
	}

	redirectURL := fmt.Sprintf("%s?%s", GithubAuthUrl, params.Encode())
	w.Header().Set("Access-Control-Allow-Origin", BaseUrl)
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
