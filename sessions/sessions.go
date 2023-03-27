package sessions

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
)

// Store the cookie store which is going to store session data in the cookie
// The Store key must be stored in the enviroment variable "SESSION_KEY"
var Store *sessions.CookieStore

// var session *sessions.Session

func Init(sessionKey string) {
	if sessionKey == "" {
		log.Fatal("Session key cannot be empty")
	}
	Store = sessions.NewCookieStore([]byte(sessionKey))
}

// IsLoggedIn will check if the user has an active session and return True
func IsLoggedIn(r *http.Request) bool {
	session, err := Store.Get(r, "session")
	if err == nil && (session.Values["loggedin"] == "true") {
		return true
	}
	return false
}

// GetCurrentUser returns the email of the logged in user
func GetCurrentUser(r *http.Request) string {
	if session, err := Store.Get(r, "session"); err == nil && session.Values["loggedin"] == "true" {
		return session.Values["email"].(string)
	}
	return ""
}
