package handlers

import (
	"log"
	"net/http"
	"strconv"

	"forum/backend/repositories"
	util "forum/backend/util"
)

func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Method not allowed", r.Method)
		util.ErrorHandler(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println("error parsing form:", err)
		util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
		return
	}

	reactionType := r.FormValue("reaction")
	postID, _ := strconv.Atoi(r.FormValue("post_id"))

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

	check, reaction := repositories.CheckReactions(util.Database, sessionData["userId"].(int), postID)

	if !check {
		_, err := repositories.InsertRecord(util.Database, "tblReactions", []string{"user_id", "post_id", "reaction"}, sessionData["userId"].(int), postID, reactionType)
		if err != nil {
			log.Println("Failed to insert record:", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if reactionType == reaction {
		err := repositories.UpdateReactionStatus(util.Database, sessionData["userId"].(int), postID)
		if err != nil {
			log.Println("Failed to update reaction status:", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	} else {
		err := repositories.UpdateReaction(util.Database, reactionType, sessionData["userId"].(int), postID)
		if err != nil {
			log.Println("Failed to update reaction:", err)
			util.ErrorHandler(w, "An Unexpected Error Occurred. Try Again Later", http.StatusInternalServerError)
			return
		}
	}
	r.Method = http.MethodGet
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
