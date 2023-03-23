package main

import (
	// "encoding/json"

	"log"

	//"html/template"
	"net/http"
	// "strconv"
	//"time"

	"github.com/alexanderi96/leafnet/config"
	"github.com/alexanderi96/leafnet/views"
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
	http.HandleFunc("/graph", views.RequiresLogin(views.GraphFunc))

	http.HandleFunc("/manageperson", views.RequiresLogin(views.AddPerson))
	http.HandleFunc("/view", views.RequiresLogin(views.ViewPeople))
	http.HandleFunc("/delete", views.RequiresLogin(views.DeletePerson))

	// Api
	// http.HandleFunc("/getpeople", views.GetPeople)

	log.Println("Server in ascolto sulla porta " + config.Config["leafnet_port"])
	http.ListenAndServe(":"+config.Config["leafnet_port"], nil)
}
