package views

/*Holds the fetch related view handlers*/

import (
	"html/template"
	"log"
	"net/http"
	"time"

	// "strings"
	// "strconv"
	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/sessions"
	"github.com/alexanderi96/leafnet/types"
)

const (
	token = "abcd"
)

var templates *template.Template
var homeTemplate *template.Template
var loginTemplate *template.Template
var profileTemplate *template.Template
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

func setCookie(w http.ResponseWriter) {
	c.CSRFToken = token
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "csrftoken", Value: token, Expires: expiration}
	http.SetCookie(w, &cookie)
}

// TODO: add ability to filter displayed events
func HomeFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		prepareContext(w, r)
		setCookie(w)

		homeTemplate.Execute(w, c)
	}
}

func GraphFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		prepareContext(w, r)
		setCookie(w)

		graphTemplate.Execute(w, c)
	}
}

func MyProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		prepareContext(w, r)
		setCookie(w)

		profileTemplate.Execute(w, c)
	}
}
