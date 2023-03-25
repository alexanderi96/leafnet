package views

import (
	"log"
	"net/http"

	"github.com/alexanderi96/leafnet/db"
	"github.com/alexanderi96/leafnet/types"
	"github.com/alexanderi96/leafnet/utils"
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
		// email := r.URL.Query().Get("email")

		// user, err := db.GetUserInfo(email)
		// //log.Println(user)
		// c.User = user
		// if err != nil {
		// 	// gestisci il caso in cui non ci sia la persona con l'uuid specificato
		// 	log.Println(err)
		// 	c.User = types.User{}
		// }
		// //log.Println("Viewing: ", c.User)

		templates["userprofile"].Execute(w, c)
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
