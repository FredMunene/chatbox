package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"forum/backend/models"
	"forum/backend/repositories"
	util "forum/backend/util"
)

func GetAllPosts(db *sql.DB, tmpl *template.Template, posts []models.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fetch comments for each post
		for i, post := range posts {
			comments, err := repositories.GetComments(db, post.ID)
			if err != nil {
				log.Printf("Failed to get comments: %v", err)
				util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
				return
			}

			posts[i].Comments = comments
		}

		// Set content type to text/html and serve the index page
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)

		err := tmpl.ExecuteTemplate(w, "index.html", struct {
			Posts []models.Post
		}{Posts: posts})
		if err != nil {
			log.Printf("Failed to render template: %v", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}
	}
}

func GetAllPostsAPI(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := repositories.GetPosts(db)
		if err != nil {
			log.Printf("Failed to get posts: %v", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}
		// fetch comments for each post
		for i, post := range posts {
			comments, err := repositories.GetComments(db, post.ID)
			if err != nil {
				log.Printf("Failed to get posts: %v", err)
				util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
				return
			}

			posts[i].Comments = comments
		}

		w.Header().Set("Content-Type", "application/json")

		if err = json.NewEncoder(w).Encode(posts); err != nil {
			log.Printf("Failed to encode posts to JSON: %v", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}
	}
}

// FilterPosts - Handles filtering posts by category or user
func FilterPosts(w http.ResponseWriter, r *http.Request) {
	logged := false
	if r.URL.Path != "/filter" {
		util.ErrorHandler(w, "Page does not exist", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		log.Println("Method not allowed", r.Method)
		util.ErrorHandler(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form", err)
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		return
	}

	categories := r.Form["category"]
	filter := r.FormValue("filter")

	if len(categories) != 0 {
		posts, err := repositories.FilterPostsByCategories(util.Database, categories)
		if err != nil {
			log.Println("error filtering posts:", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}

		cookie, _ := GetSessionID(r)
		_, ok := SessionStore[cookie]
		if ok {
			logged = true
		}

		PostDetails(w, r, posts, logged)
		return
	}

	cookie, err := GetSessionID(r)
	if err != nil {
		log.Println("Invalid Session:", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sessionData, err := GetSessionData(cookie)
	if err != nil {
		log.Println("Invalid Session:", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	posts := []models.Post{}

	if filter == "created" {
		posts, err = repositories.FilterPostsByUser(util.Database, sessionData["userId"].(int))
	}
	if filter == "liked" {
		posts, err = repositories.FilterPostsByLikes(util.Database, sessionData["userId"].(int))
	}
	if err != nil {
		log.Println(err)
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		return
	}

	PostDetails(w, r, posts, true)
}
