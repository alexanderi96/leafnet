package main

import (
	"net/http"
	"os"
	"github.com/alexanderi96/leafnet/person"
)



type adminPortal struct {
	password string
}

func newAdminPortal() *adminPortal {
	password := os.Getenv("ADMIN_PASSWORD")
	if password == "" {
		panic("required env var ADMIN_PASSWORD not set")
	}

	return &adminPortal{password: password}
}

func (a *adminPortal) handler(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != "admin" || pass != a.password {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - unauthorized"))
		return
	}
	w.Write([]byte("<html><h1>super secret admin portal</h1></html>"))

}

func main() {
	admin := newAdminPortal()
	personHandler := person.NewPersonHandler()
	http.HandleFunc("/persons", personHandler.Persons)
	http.HandleFunc("/persons/", personHandler.GetPerson)
	http.HandleFunc("/admin", admin.handler)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		panic(err)
	}
}