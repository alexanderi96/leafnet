package views

import (
	"log"
	"strconv"
	"net/http"	
	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/sessions"
	"github.com/alexanderi96/leafnet/utils"
	"github.com/alexanderi96/leafnet/types"
)

func SignUpFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Redirect(w, r, "/login/", http.StatusBadRequest)
		return
	}
	r.ParseForm()

	var u types.User

	if r.Form.Get("cicerone") != "on" {
		u = parseGlobe(r)
	} else {
		u = parseCice(r)
	}
	
	e = db.CreateUser(u)
	if e != nil {
		http.Error(w, "Unable to sign user up", http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, "/login/", 302)
	}
}

func DeleteMyAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	r.ParseForm()

	if e := db.DeleteSelectedUser(r.Form.Get("email"), r.Form.Get("password")); e != nil {
		//TODO: handle better this behaviour
		http.Redirect(w, r, "/myprofile/", http.StatusUnauthorized)
	} else {
		log.Println("sas")
		http.Redirect(w, r, "/logout/", 302)
	}
}

func UpdateAccountInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/myprofile/", http.StatusUnauthorized)
		return
	}
	
	uid := db.GetUserID(sessions.GetCurrentUser(r))
	if db.IsCicerone(uid) {
		if e = updateCice(r); e != nil {
			log.Println("(uid: ", uid, ") error updating user information, keeping old details\nError: ", e)
			http.Redirect(w, r, "/myprofile/", http.StatusInternalServerError)
		}
	} else if e = updateGlobe(r); e != nil {
		log.Println("(uid: ", uid, ") error updating user information, keeping old details\nError: ", e)
		http.Redirect(w, r, "/myprofile/", http.StatusInternalServerError)
	}

	//let's make the user log back in in order to load the new info (must find a better way)
	http.Redirect(w, r, "/logout/", 302)	
}
