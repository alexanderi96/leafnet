package views

import (
	//"encoding/json"
	"log"
	//"html/template"
	"net/http"

	"time"

	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/types"
	//"github.com/alexanderi96/leafnet/sessions"
)

func AddPerson(w http.ResponseWriter, r *http.Request) {
	prepareContext(w, r)

	if r.Method == "POST" {
		r.ParseForm()

		p := types.Person{}

		uuid := r.Form.Get("uuid")
		p.Node.UUID = uuid

		firstName := r.Form.Get("first_name")
		p.FirstName = firstName

		lastName := r.Form.Get("last_name")
		p.LastName = lastName

		if bDate := r.Form.Get("birth_date"); bDate != "" {
			// birthDate, err := strconv.ParseInt(bDate, 10, 64)
			birthDate, err := time.Parse("2006-01-02", bDate)
			if err != nil {
				panic(err)
			}
			p.BirthDate = birthDate.Unix()
		}

		if dDate := r.Form.Get("death_date"); dDate != "" {
			//deathDate, err := strconv.ParseInt(dDate, 10, 64)
			deathDate, err := time.Parse("2006-01-02", dDate)
			if err != nil {
				panic(err)
			}
			p.DeathDate = deathDate.Unix()
		}

		parent1 := r.Form.Get("parent1")
		p.Parent1 = parent1

		parent2 := r.Form.Get("parent2")
		p.Parent2 = parent2

		bio := r.Form.Get("bio")
		p.Bio = bio

		// Salva i dati nel database
		err := db.ManagePerson(&p)
		if err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, "/", http.StatusFound)
	} else if r.Method == "GET" {
		uuid := r.URL.Query().Get("uuid")
		var err error
		c.Person, err = db.GetPerson(uuid)
		if err != nil {
			// gestisci il caso in cui non ci sia la persona con l'uuid specificato
			log.Println(err)
			c.Person = types.Person{}
		}

		c.Page.IsOwner = c.Person.Node.Owner == "" || c.Person.Node.Owner == c.User.Email
		c.Page.IsDisabled = !c.Page.IsOwner

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
	prepareContext(w, r)

	peopleTemplate.Execute(w, c)
}

// func GetPeople(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
//     json.NewEncoder(w).Encode(db.GetPersons())
// }
