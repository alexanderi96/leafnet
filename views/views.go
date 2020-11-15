package views

/*Holds the fetch related view handlers*/

import (
	"html/template"
	"log"
	"net/http"
	"time"
	"strings"
	"strconv"
	"github.com/alexanderi96/leafnet/sessions"
	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/types"
)

const (
	token = "abcd"
)

var templates *template.Template
var homeTemplate *template.Template
var loginTemplate *template.Template
var profileTemplate *template.Template
var eventTemplate *template.Template

var c types.Context
var e error

//TODO: add ability to filter displayed events
func HomeFunction(w http.ResponseWriter, r *http.Request) {
	
	if r.Method == "GET" {

		c.Utente, e = db.GetUserInfo(sessions.GetCurrentUser(r))
		log.Println(e)
		c.Events, e = db.GetEvents()

		if e != nil   {
			log.Println("Internal server error retriving context")
			http.Redirect(	w, r, "/", http.StatusInternalServerError)
		} else {
			c.CSRFToken = token 
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie := http.Cookie{Name: "csrftoken", Value: token, Expires: expiration}
			http.SetCookie(w, &cookie)

			homeTemplate.Execute(w, c)
		}
	}
}

func MyProfile(w http.ResponseWriter, r *http.Request) {
		
	if r.Method == "GET" {

		if c.Utente, e = db.GetUserInfo(sessions.GetCurrentUser(r)); e != nil {
			log.Println("Internal server error retriving user info")
			http.Redirect(	w, r, "/", http.StatusInternalServerError)
		} else {
			c.CSRFToken = token 
			expiration := time.Now().Add(365 * 24 * time.Hour)
			cookie := http.Cookie{Name: "csrftoken", Value: token, Expires: expiration}
			http.SetCookie(w, &cookie)

			profileTemplate.Execute(w, c)
		}
	}
}
