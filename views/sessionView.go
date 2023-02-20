package views

import(
	"log"
	"net/http"

	"github.com/alexanderi96/leafnet/types"
	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/sessions"
)

//RequiresLogin is a middleware which will be used for each httpHandler to check if there is any active session
func RequiresLogin(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !sessions.IsLoggedIn(r) {
			log.Print("invalid session")
			http.Redirect(w, r, "/login", 302)
			return
		}
		handler(w, r)
	}
}


//LogoutFunc Implements the logout functionality. WIll delete the session information from the cookie store
func LogoutFunc(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Store.Get(r, "session")
	log.Println("Logout function")
	if err == nil { //If there is no error, then remove session
		if session.Values["loggedin"] != "false" {
			session.Values["loggedin"] = "false"
			session.Save(r, w)
		}
	}
	http.Redirect(w, r, "/login", 302) //redirect to login irrespective of error or not
}

//LoginFunc implements the login functionality, will add a cookie to the cookie store for managing authentication
func LoginFunc(w http.ResponseWriter, r *http.Request) {
	session, err := sessions.Store.Get(r, "session")

	if err != nil {
	    http.Error(w, err.Error(), http.StatusInternalServerError)
	    return
	}

	switch r.Method {
	case "GET":
		log.Print("Inside GET")
		loginTemplate.Execute(w, nil)
	case "POST":
		log.Print("Inside POST")
		r.ParseForm()
		email := r.Form.Get("email")
		password := r.Form.Get("password")

		if (email != "" && password != "") && db.ValidUser(email, password) {
			session.Values["loggedin"] = "true"
			session.Values["email"] = email
			session.Save(r, w)
			log.Print("user ", email, " is authenticated")
			http.Redirect(w, r, "/", 302)
			return
		}
		log.Print("Invalid user " + email)
		loginTemplate.Execute(w, nil)
	default:
		http.Redirect(w, r, "/login", http.StatusUnauthorized)
	}
}

func SignUpFunc(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Redirect(w, r, "/login", http.StatusBadRequest)
		return
	}
	r.ParseForm()

	u := parseUser(r)
	
	e = db.NewUser(&u)
	if e != nil {
		http.Error(w, "Unable to sign user up", http.StatusInternalServerError)
	} else {
		http.Redirect(w, r, "/login", 302)
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
		http.Redirect(w, r, "/logout/", 302)
	}
}

func parseUser(r *http.Request) (u types.User) {
	r.ParseForm()

	u = types.User{
		UserName: r.Form.Get("user_name"),
		Email: r.Form.Get("email"),
		Password: r.Form.Get("password"),
		Person: r.Form.Get("person")}
	return
}