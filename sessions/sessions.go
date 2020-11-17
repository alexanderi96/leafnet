package sessions

import (
	"net/http"
	"os"
	"log"
	"time"
	"github.com/gorilla/sessions"
)

//Store the cookie store which is going to store session data in the cookie
//The Store key must be stored in the enviroment variable "SESSION_KEY"
var Store = sessions.NewCookieStore([]byte(os.Getenv("APP_SESSION_KEY")))
var session *sessions.Session

func BodyGuard(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isLoggedIn(r) {
			log.Print("invalid session")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid session"))
			return
		}
		handler(w, r)
	}
}

func isLoggedIn(r *http.Request) bool {
	session, err := Store.Get(r, "session")
	if err == nil && session.Values["Timestamp"] != nil && (time.Now().UnixNano() - 1800000 > session.Values["Timestamp"].(int64)) {
		return true
	}
	return false
}

//GetCurrentUser returns the email of the logged in user
func GetCurrentUserId(r *http.Request) (id string) {
	session, err := Store.Get(r, "session")
	if err == nil {
		return session.Values["Id"].(string)
	}
	return
}
