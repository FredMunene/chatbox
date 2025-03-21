package handlers

import (
	"log"
	"net/http"

	"forum/backend/repositories"
	util "forum/backend/util"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		util.ErrorHandler(w, "Page does not exist", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		log.Println("Method not allowed", r.Method)
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, _ := GetSessionID(r)
	_, ok := SessionStore[cookie]
	if ok {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	// Load posts
	posts, err := repositories.GetPosts(util.Database)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		return
	}

	PostDetails(w, r, posts, false)
}
