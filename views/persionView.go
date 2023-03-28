package views

import (
	"log"
	"net/http"
	"time"

	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/types"
)

func AddPerson(w http.ResponseWriter, r *http.Request) {
	prepareContext(w, r)

	if r.Method == "POST" {
		r.ParseForm()

		p := types.Person{}

		uuid := r.Form.Get("uuid")
		p.Node.UUID = uuid

		if owner := r.Form.Get("owner"); owner == "" {
			p.Node.Owner = c.User.UserName
		} else {
			p.Node.Owner = owner
		}

		firstName := r.Form.Get("first_name")
		p.FirstName = firstName

		lastName := r.Form.Get("last_name")
		p.LastName = lastName

		if bDate := r.Form.Get("birth_date"); bDate != "" {
			birthDate, err := time.Parse("2006-01-02", bDate)
			if err != nil {
				panic(err)
			}
			p.BirthDate = birthDate.Unix()
		}

		if dDate := r.Form.Get("death_date"); dDate != "" {
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

		// Save the person to the database
		if err := db.ManagePerson(&p); err != nil {
			WriteError(w, err)
			return
		}

		http.Redirect(w, r, "/view", http.StatusFound)
	} else if r.Method == "GET" {
		uuid := r.URL.Query().Get("uuid")

		if p, err := db.GetPerson(uuid); err != nil {
			log.Println(err)
			c.Person = types.Person{}
		} else {
			c.Person = p
		}

		c.Page.IsOwner = c.Person.Node.Owner == "" || c.Person.Node.Owner == c.User.UserName
		c.Page.IsDisabled = !c.Page.IsOwner

		if err := templates["manageperson"].Execute(w, c); err != nil {
			WriteError(w, err)
			return
		}
	}
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		uuid := r.URL.Query().Get("uuid")

		if err := db.DeletePerson(uuid); err != nil {
			WriteError(w, err)
			return
		}
		log.Println("Deleted person: ", uuid)

		http.Redirect(w, r, "/view", http.StatusSeeOther)
	}
}

func ViewPeople(w http.ResponseWriter, r *http.Request) {
	prepareContext(w, r)

	if err := templates["people"].Execute(w, c); err != nil {
		WriteError(w, err)
		return
	}
}
