package auth

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"log"
	"net/http"

	"forum/backend/util"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func ComparePassword(userPassword string, dbPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(userPassword))
	if err != nil {
		return err
	}
	return nil
}

func handleUserAuth(w http.ResponseWriter, email, username string) bool {
	var userID int

	err := util.Database.QueryRow(
		"SELECT id from tblUsers WHERE email = ?", email).Scan(&userID)

	//  Create a new user if not found
	if errors.Is(err, sql.ErrNoRows) {
		res, err := util.Database.Exec(
			"INSERT INTO tblUsers(username, email) VALUES(?,?)",
			username, email,
		)
		if err != nil {
			log.Printf("User creation failed: %v", err)
			return false

		}
		id, _ := res.LastInsertId()
		userID = int(id)
	} else if err != nil {
		log.Printf("Database error: %v", err)
		return false
	}

	setSessionCookie(w, userID)
	return true
}

// SetSessionCookie sets a session cookie for the given user ID.
func setSessionCookie(w http.ResponseWriter, userID int) {
	token := generateSessionToken()
	_, err := util.Database.Exec(
		"INSERT INTO tblSessions(user_id, session_token) VALUES(?, ?)",
		userID, token,
	)
	if err != nil {
		log.Printf("Session creation failed: %v", err)
		return
	}

	// Set the session cookie for the user
	http.SetCookie(w, &http.Cookie{
		Name:     "forum_session",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   86400, // 1 day
	})
}

func generateSessionToken() string {
	// Generate a new UUID (UUIDv4)
	u := uuid.New()

	// Convert the UUID to its byte representation
	uBytes, err := u.MarshalBinary()
	if err != nil {
		log.Printf("Error marshaling UUID:", err)
		return ""
	}

	// Encode the bytes in URL-safe base64 format
	return base64.URLEncoding.EncodeToString(uBytes)
}
