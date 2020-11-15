package main

import (
	"flag"
	//"fmt"
	"net/http"
	"strings"
	"log"

	"github.com/alexanderi96/leafnet/config"
	"github.com/alexanderi96/leafnet/views"
)

func main(){

	values, err := config.ReadConfig("config.json")
	var port string

	if err != nil {
		flag.StringVar(&port, "Port", "", "Ip address")
		flag.Parse()

		if !strings.HasPrefix(port, ":") {
			port = ":" + port
			log.Println("Port is: " + port)
		}

		values.ServerPort = port
	}

	views.PopulateTemplates()

	// login-logout handlers
	http.HandleFunc("/login/", views.LoginFunc)
	http.HandleFunc("/signup/", views.SignUpFunc) 
	http.HandleFunc("/logout/", views.RequiresLogin(views.LogoutFunc)) 
	http.HandleFunc("/delete-user/", views.RequiresLogin(views.DeleteMyAccount))
	http.HandleFunc("/update-user/", views.RequiresLogin(views.UpdateAccountInfo))

	http.HandleFunc("/", views.RequiresLogin(views.HomeFunction)) // User must be logged in to access the homepage. If he isn't he will recieve a page with the project description and the possibility to login or register
	http.HandleFunc("/myprofile/", views.RequiresLogin(views.MyProfile))

	// listen to the port 8081 for any request
	log.Println("Running leafnet on ", values.ServerPort)
	log.Fatal(http.ListenAndServe(values.ServerPort, nil))
}