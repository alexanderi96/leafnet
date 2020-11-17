package main

import (
	"encoding/json"
	"net/http"
	"time"
	"log"
	"github.com/alexanderi96/leafnet/person"
	"github.com/alexanderi96/leafnet/utils"
	"github.com/alexanderi96/leafnet/sessions"

)



type app struct {
	ph *person.PersonHandler
	uh *person.UserHandler
}

func newApp() *app {
	return &app{
		ph: person.NewPersonHandler(),
		uh: person.NewUserHandler(),
	}
}

func (a *app) grantAccess(w http.ResponseWriter, r *http.Request) {
	var credentials person.User
	err := json.Unmarshal(utils.CheckJsonAndGet(w, r), &credentials)

	if err != nil || (credentials.Username == "" && credentials.Password == "") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	session, err := sessions.Store.Get(r, "session")
	for _, user := range a.uh.Store {
		if credentials.Username == user.Username && credentials.Password == user.Password {
			session.Values["Timestamp"] = time.Now().UnixNano()
			session.Values["Id"] = user.Id
			session.Save(r, w)
			log.Print("user ", user.Username , " granted for 30 minutes")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("session granted for 30 minutes"))
		}
	}
}

func main() {
	
	app := newApp()
	http.HandleFunc("/grantAccess", app.grantAccess)

	// Use GET to get a list of users/persons, 
	// use POST instead, if you want to save one
	http.HandleFunc("/users", sessions.BodyGuard(app.uh.Users))
	http.HandleFunc("/persons", sessions.BodyGuard(app.ph.Persons))
	
	http.HandleFunc("/WhoAmI", sessions.BodyGuard(app.uh.WhoAmI))



	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		panic(err)
	}
}