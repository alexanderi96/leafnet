package views

import (
	"log"
	"net/http"

	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/types"
)

func UserPage(w http.ResponseWriter, r *http.Request) {
	prepareContext(w, r)

	if r.Method == "POST" {
		r.ParseForm()

		// u := types.User{}

		// uuid := r.Form.Get("uuid")
		// p.Node.UUID = uuid

		// firstName := r.Form.Get("first_name")
		// p.FirstName = firstName

		// lastName := r.Form.Get("last_name")
		// p.LastName = lastName

		// if bDate := r.Form.Get("birth_date"); bDate != "" {
		// 	// birthDate, err := strconv.ParseInt(bDate, 10, 64)
		// 	birthDate, err := time.Parse("2006-01-02", bDate)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	p.BirthDate = birthDate.Unix()
		// }

		// if dDate := r.Form.Get("death_date"); dDate != "" {
		// 	//deathDate, err := strconv.ParseInt(dDate, 10, 64)
		// 	deathDate, err := time.Parse("2006-01-02", dDate)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	p.DeathDate = deathDate.Unix()
		// }

		// parent1 := r.Form.Get("parent1")
		// p.Parent1 = parent1

		// parent2 := r.Form.Get("parent2")
		// p.Parent2 = parent2

		// bio := r.Form.Get("bio")
		// p.Bio = bio

		// // Salva i dati nel database
		// err := db.ManagePerson(&p)
		// if err != nil {
		// 	log.Println(err)
		// }

		http.Redirect(w, r, "/", http.StatusFound)
	} else if r.Method == "GET" {
		email := r.URL.Query().Get("email")

		var err error
		c.User, err = db.GetUserInfo(email)
		if err != nil {
			// gestisci il caso in cui non ci sia la persona con l'uuid specificato
			log.Println(err)
			c.User = types.User{}
		}
		log.Println("Viewing: ", c.User)

		userPagetemplate.Execute(w, c)
	}
}
