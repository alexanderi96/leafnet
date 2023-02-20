package views

import (
	"encoding/json"
	"log"
	//"html/template"
	"net/http"
	"strconv"
	//"time"

	"github.com/alexanderi96/leafnet/types"
	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/sessions"
)

func AddPerson(w http.ResponseWriter, r *http.Request) {
	c.User, e = db.GetUserInfo(sessions.GetCurrentUser(r))

	if e != nil   {
		log.Println("Internal server error retriving context")
		http.Redirect(	w, r, "/", http.StatusInternalServerError)
	} else if r.Method == "POST" {

		uuid := r.FormValue("uuid")
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")

		date := r.FormValue("birth_date")
		birthDate, err := strconv.ParseInt(date, 10, 64)
		if err != nil {
		    panic(err)
		}

		dDate := r.FormValue("death_date")
		deathDate, err := strconv.ParseInt(dDate, 10, 64)
		if err != nil {
		    panic(err)
		}

		parent1 := r.FormValue("parent1")
		parent2 := r.FormValue("parent2")

		bio := r.FormValue("bio")

		p := types.Person{
			Node: types.Node{UUID: uuid},
			FirstName:      firstName,
			LastName:       lastName,
			BirthDate:      birthDate,
			DeathDate:      deathDate,
			Parent1:        parent1,
			Parent2:        parent2,
			Bio:            bio,
		}

		log.Println(("saved person : "), p)

		// Salva i dati nel database o nella memoria

		db.NewPerson(&p)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if r.Method == "GET" {
		uuid := r.URL.Query().Get("uuid")

		var err error 
		c.Person, err = db.GetPerson(uuid)
		if err != nil {
			// gestisci il caso in cui non ci sia la persona con l'uuid specificato
			log.Println(err)
			c.Person = types.Person{}
		}
		log.Println("person to edit: ", c.Person)

		managePersonTemplate.Execute(w, c)
	}
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		uuid := r.URL.Query().Get("uuid")
		
		log.Println("person to delete: ", uuid)

		db.DeletePerson(uuid)

		log.Println("deleted")

		http.Redirect(w, r, "/view", http.StatusSeeOther)
	}
}

func ViewPeople(w http.ResponseWriter, r *http.Request) {
	c.User, e = db.GetUserInfo(sessions.GetCurrentUser(r))

	if e != nil   {
		log.Println("Internal server error retriving context")
		http.Redirect(	w, r, "/", http.StatusInternalServerError)
	}

	c.Persons = db.GetPersons()
	peopleTemplate.Execute(w, c)
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(db.GetPersons())
}