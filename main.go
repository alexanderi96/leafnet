package main

import (
	// "encoding/json"
	"log"
	//"html/template"
	"net/http"
	// "strconv"
	//"time"

	"github.com/alexanderi96/leafnet/views"
	// "github.com/alexanderi96/leafnet/types"
	// "github.com/alexanderi96/leafnet/db"
)

func main() {

	views.PopulateTemplates()

	// login-logout handlers
	http.HandleFunc("/login", views.LoginFunc)
	http.HandleFunc("/signup", views.SignUpFunc) 
	http.HandleFunc("/logout", views.RequiresLogin(views.LogoutFunc)) 
	http.HandleFunc("/delete-user", views.RequiresLogin(views.DeleteMyAccount))
	//http.HandleFunc("/update-user", views.RequiresLogin(views.UpdateAccountInfo))


	http.HandleFunc("/", views.RequiresLogin(views.HomeFunc))

	http.HandleFunc("/manageperson", views.RequiresLogin(views.AddPerson))
	http.HandleFunc("/view", views.RequiresLogin(views.ViewPeople))
	http.HandleFunc("/delete", views.RequiresLogin(views.DeletePerson))

	// Api
	// http.HandleFunc("/getpeople", views.GetPeople)

	log.Println("Server in ascolto sulla porta 8080")
	http.ListenAndServe(":8080", nil)
}
