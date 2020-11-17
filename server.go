package main

import (
	"net/http"
	"os"
	"github.com/alexanderi96/leafnet/person"
	"github.com/alexanderi96/leafnet/utils"
	"github.com/alexanderi96/leafnet/user"

	"github.com/gorilla/sessions"
)



type app struct {
	store *sessions.CookieStore
	ph *person.PersonHandler
	uh *user.UserHandler
}

func newApp() *app {
	//TODO: check session key

	return &app{
		store: sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY"))),
		ph: person.PersonHandler(),
	}
}

func (a *app) bodyGuard(handler func(w http.ResponseWriter, r *http.Request)) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !a.IsLoggedIn(r) {
			log.Print("invalid session")
			http.Redirect(w, r, "/grantAccess", 302)
			return
		}
		handler(w, r)
	}
}

func (a *app) IsLoggedIn(r *http.Request) {
	return true
}

func (a *app) grantAccess(w http.ResponseWriter, r *http.Request) {
	var credentials user.User
	err = json.Unmarshal(utils.CheckJsonAndGet(r), &user)

	if err != nil || (credentials.Username != "" && credentials.Password != "") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	for _, user := range a.uh.Store() {
		if credentials.Username == user.Username && credentials.Password == user.Password {
			a.store.Values["Timestamp"] = time.Now().UnixNano()
			a.store.Values["Id"] = user.Id
			a.store.Save(r, w)
			log.Print("user ", user.Username , " granted for 30 minutes")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("session granted for 30 minutes"))
		}
	}
}

func main() {
	
	app := newApp()
	http.HandleFunc("/grantAccess", app.grantAccess)
	http.HandleFunc("/persons", app.ph.Persons)
	http.HandleFunc("/getPersons/", ph.GetPerson)
	http.HandleFunc("/app", app.bodyGuard())
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		panic(err)
	}
}