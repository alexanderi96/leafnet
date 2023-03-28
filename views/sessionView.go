package views

import (
	"errors"
	"log"
	"net/http"

	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/sessions"
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
	if err != nil {
		WriteError(w, err)
		return
	}

	//If there is no error, then remove session
	log.Println("Removing session for user: ", session.Values["email"])
	if session.Values["loggedin"] != "false" {
		session.Values["loggedin"] = "false"
		session.Save(r, w)
	}

	http.Redirect(w, r, "/login", http.StatusFound) //redirect to login irrespective of error or not
}

// LoginFunc implements the login functionality, will add a cookie to the cookie store for managing authentication
func LoginFunc(w http.ResponseWriter, r *http.Request) {

	prepareContext(w, r)
	session, err := sessions.Store.Get(r, "session")

	if err != nil {
		WriteError(w, err)
		return
	}

	switch r.Method {
	case "GET":
		log.Print("New access to the login page")
		if err := templates["login"].Execute(w, c); err != nil {
			WriteError(w, err)
			return
		}

	case "POST":
		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		log.Print("Attempting to login with email: ", email)

		if hash, err := db.GetUserPasswdHash(email); err != nil {
			WriteError(w, err)
			return

		} else if hash == "" {
			log.Print("User not found")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		} else if res, err := utils.CheckStrHash(password, hash); err != nil || !res {
			log.Print("Wrong password")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		session.Values["loggedin"] = "true"
		session.Values["email"] = email
		session.Options.MaxAge = 3600 // imposto l'età massima della sessione a 1 ora

		if err := session.Save(r, w); err != nil {
			WriteError(w, err)
			return
		}

		log.Print("user ", email, " is authenticated")
		http.Redirect(w, r, "/", http.StatusFound)
		return

	default:
		WriteError(w, errors.New("invalid method"))
	}
}

func SignUpFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	r.ParseForm()

	u, err := parseUser(r)
	log.Println("Attempting to sign up with Email: ", u.Email, " and UserName: ", u.UserName)

	if err != nil {
		WriteError(w, err)
		return
	} else if exists, err := checkIfUserExists(u); err != nil {
		WriteError(w, err)
		return
	} else if exists {
		WriteError(w, errors.New("user already exists"))
	}

	if err = db.NewUser(&u); err != nil {
		WriteError(w, err)
		return
	}

	log.Println("User ", u.UserName, " has been created")
	http.Redirect(w, r, "/login", http.StatusFound)
}

// TODO: add ability to filter displayed events
func HomeFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		prepareContext(w, r)
		log.Println("Attempting to access home page")
		if err := templates["home"].Execute(w, c); err != nil {
			WriteError(w, err)
			return
		}
	}
}

func GraphFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		prepareContext(w, r)
		uuid := r.URL.Query().Get("uuid")

		//TODO: fix this
		if uuid != "" {
			log.Println("Attempting to access graph page with uuid: ", uuid)
			if ppl, err := db.FetchAncestors(uuid); err == nil {
				log.Println("Fetched ", len(ppl), " persons")
				c.Persons = ppl
			} else {
				WriteError(w, err)
			}
		} else if ppl, err := db.GetPersons(); err == nil {
			c.Persons = ppl
		} else {
			WriteError(w, err)
		}

		log.Println("Attempting to access graph page")
		if err := templates["graph"].Execute(w, c); err != nil {
			WriteError(w, err)
			return
		}
	}
}
