package views

/*Holds the fetch related view handlers*/

import (
	"html/template"
	"log"
	"net/http"

	// "strings"
	// "strconv"
	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/sessions"
	"github.com/alexanderi96/leafnet/types"
)

// const (
// 	token = "abcd"
// )

var templates *template.Template
var homeTemplate *template.Template
var loginTemplate *template.Template
var userPagetemplate *template.Template
var peopleTemplate *template.Template
var managePersonTemplate *template.Template
var graphTemplate *template.Template

var c types.Context
var e error

func prepareContext(w http.ResponseWriter, r *http.Request) {
	//load user info
	if c.User, e = db.GetUserInfo(sessions.GetCurrentUser(r)); e != nil {
		log.Println("Internal server error retriving user info")
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}

	//load persons
	c.Persons = db.GetPersons()
}

// func setCookie(w http.ResponseWriter) {
// 	c.CSRFToken = token

// 	cookie := http.Cookie{
// 		Name:     "csrftoken",
// 		Value:    token,
// 		Path:     "/",
// 		MaxAge:   3600,
// 		HttpOnly: true,
// 		Secure:   true,
// 		SameSite: http.SameSiteLaxMode,
// 	}

// 	log.Println("Setting cookie: ", cookie)

// 	http.SetCookie(w, &cookie)
// }

// func checkCookie(r *http.Request) bool {
// 	cookie, err := r.Cookie("csrftoken")
// 	log.Println("Checking cookie: ", cookie)
// 	if err != nil {
// 		log.Println("Error getting cookie: ", err)
// 		return false
// 	}
// 	if cookie.Value != token {
// 		log.Println("Cookie is not valid")
// 		return false
// 	}
// 	if time.Now().After(cookie.Expires) {
// 		log.Println("Cookie is expired: ", time.Now(), "Cookie expiration: ", cookie.Expires)
// 		return false
// 	}

// 	return true
// }

// TODO: add ability to filter displayed events
func HomeFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		prepareContext(w, r)

		homeTemplate.Execute(w, c)
	}
}

func GraphFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		prepareContext(w, r)

		graphTemplate.Execute(w, c)
	}
}

// func MyProfile(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "GET" {
// 		prepareContext(w, r)

// 		userPagetemplate.Execute(w, c)
// 	}
// }
