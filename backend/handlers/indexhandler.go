package handlers

import (
	"log"
	"net/http"

	"forum/backend/repositories"
	util "forum/backend/util"
)

type RequestData struct {
	ID string `json:"id"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/home" {
		util.ErrorHandler(w, "Page does not exist", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		log.Println("wrong method used", r.Method)
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := GetSessionID(r)
	if err != nil {
		log.Println("Invalid Session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sessionData, err := GetSessionData(cookie)
	if err != nil {
		log.Println("Invalid Session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Fetch user information
	_, err = repositories.GetUserByEmail(sessionData["userEmail"].(string))
	if err != nil {
		log.Printf("Invalid session token: %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	posts, err := repositories.GetPosts(util.Database)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		return
	}
	PostDetails(w, r, posts, true)
}
