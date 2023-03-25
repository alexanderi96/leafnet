package views

import (
	"log"
	"net/http"

	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/sessions"
	"github.com/alexanderi96/leafnet/types"
	"github.com/alexanderi96/leafnet/utils"
)

// RequiresLogin is a middleware which will be used for each httpHandler to check if there is any active session
func RequiresLogin(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessions.IsLoggedIn(r) { // || !checkCookie)r) {}
			log.Print("invalid session")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		handler(w, r)
	}
}

// LogoutFunc Implements the logout functionality. WIll delete the session information from the session store
func LogoutFunc(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Store.Get(r, "session")
	if err == nil { //If there is no error, then remove session
		log.Println("removing session for user: ", session.Values["email"])
		if session.Values["loggedin"] != "false" {
			session.Values["loggedin"] = "false"
			session.Save(r, w)
		}
	}

	http.Redirect(w, r, "/login", http.StatusFound) //redirect to login irrespective of error or not
}

// LoginFunc implements the login functionality, will add a cookie to the cookie store for managing authentication
func LoginFunc(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Store.Get(r, "session")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch r.Method {
	case "GET":
		log.Print("New access to the login page")
		templates["login"].Execute(w, c)
	case "POST":
		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		log.Print("Attempting to login with email: ", email)
		if (email != "" && password != "") && utils.CheckStrHash(password, db.GetUserPasswdHash(email)) {
			session.Values["loggedin"] = "true"
			session.Values["email"] = email
			session.Options.MaxAge = 3600 // imposto l'età massima della sessione a 1 ora

			if err := session.Save(r, w); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// setCookie(w)
			log.Print("user ", email, " is authenticated")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		log.Print("Invalid user " + email)
		http.Redirect(w, r, "/login", 401)

	default:
		http.Redirect(w, r, "/login", http.StatusMethodNotAllowed)
	}
}

func SignUpFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	r.ParseForm()

	u := parseUser(r)
	if u == (types.User{}) {
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}

	log.Println("Attempting to sign up with email: ", u.Email)

	e = db.NewUser(&u)
	if e != nil {
		http.Error(w, "Unable to sign user up", http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// TODO: add ability to filter displayed events
func HomeFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		prepareContext(w, r)

		templates["home"].Execute(w, c)
	}
}

func GraphFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		prepareContext(w, r)

		templates["graph"].Execute(w, c)
	}
}
