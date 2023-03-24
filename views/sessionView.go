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

// LogoutFunc Implements the logout functionality. WIll delete the session information from the cookie store
func LogoutFunc(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Store.Get(r, "session")
	if err == nil { //If there is no error, then remove session
		log.Println("removing session for user: ", session.Values["email"])
		if session.Values["loggedin"] != "false" {
			session.Values["loggedin"] = "false"
			session.Save(r, w)
		}
	}

	loginTemplate.Execute(w, http.StatusFound) //redirect to login irrespective of error or not
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
		loginTemplate.Execute(w, nil)
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
		http.Redirect(w, r, "/login", http.StatusUnauthorized)

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
		http.Redirect(w, r, "/login", http.StatusAccepted)
	}
}

func DeleteMyAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	r.ParseForm()

	if e := db.DeleteSelectedUser(r.Form.Get("email"), r.Form.Get("password")); e != nil {
		log.Print("Error Deleting Account ", e)
		//TODO: handle better this behaviour
		http.Redirect(w, r, "/myprofile", http.StatusUnauthorized)
	} else {
		log.Println("sas")
		http.Redirect(w, r, "/logout", http.StatusAccepted)
	}
}

func parseUser(r *http.Request) (u types.User) {
	r.ParseForm()

	hashedPwd, err := utils.EncryptStr(r.Form.Get("password"))
	if err != nil {
		log.Print("Error encrypting password: ", err)
		return
	}

	u = types.User{
		UserName: r.Form.Get("user_name"),
		Email:    r.Form.Get("email"),
		Password: hashedPwd,
		Person:   r.Form.Get("person")}
	return
}
