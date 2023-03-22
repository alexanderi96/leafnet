package sessions

import (
	"net/http"
	"log"
	"github.com/gorilla/sessions"
	"github.com/alexanderi96/leafnet/config"
)

//Store the cookie store which is going to store session data in the cookie
//The Store key must be stored in the enviroment variable "SESSION_KEY"
var Store *sessions.CookieStore
var session *sessions.Session

func init() {
	Store = sessions.NewCookieStore([]byte(config.Config["leafnet_session_key"]))
}

//IsLoggedIn will check if the user has an active session and return True
func IsLoggedIn(r *http.Request) bool {
	session, err := Store.Get(r, "session")
	log.Println(session.Values["loggedin"])
	if err == nil && (session.Values["loggedin"] == "true") {
		return true
	}
	return false
}

//GetCurrentUser returns the email of the logged in user
func GetCurrentUser(r *http.Request) string {
	session, err := Store.Get(r, "session")
	if err == nil {
		return session.Values["email"].(string)
	}
	return ""
}