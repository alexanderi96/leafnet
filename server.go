package main

import (
	"encoding/json"
	"net/http"
	"log"
	"fmt"
	"time"
	"sync"
	"strings"
	"os"
	"io/ioutil"
	"math/rand"
)

type Person struct {
	id 			string`json:"id"`
	name 		string`json:"name"`
	surname 	string`json:"surname"`
	birthdate 	int64 `json:"birthdate"`
}

type personHandler struct {
	sync.Mutex
	store map[string]Person
}

func (h *personHandler) persons(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
}

func (h *personHandler) get(w http.ResponseWriter, r *http.Request) {
	persons := make([]Person, len(h.store))

	h.Lock()
	i := 0
	for _, person := range h.store {
		persons[i] = person
		i++
	}
	h.Unlock()


	jsonBytes, err := json.Marshal(persons)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *personHandler) getRandomPerson(w http.ResponseWriter, r *http.Request) {
	ids := make([]string, len(h.store))
	h.Lock()

	i:=0
	for _, id := range h.store {

		ids[i] = id.id
		i++
	}
	defer h.Unlock()

	var target string
	if len(ids) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if len(ids) == 1 {
		target = ids[0]
	} else {
		rand.Seed(time.Now().UnixNano())
		target = ids[rand.Intn(len(ids))]
	}
	log.Print(target)
	w.Header().Add("location", fmt.Sprintf("/persons/%s", target))
	w.WriteHeader(http.StatusFound)

}
func (h *personHandler) getPerson(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if parts[2] == "random" {
		h.getRandomPerson(w, r)
		return
	}

	h.Lock()
	person, ok := h.store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(person)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *personHandler) post(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ct := r.Header.Get("content-type")
	if ct != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
	}

	var person Person
	err = json.Unmarshal(bodyBytes, &person)
	log.Print(person)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	person.id = fmt.Sprintf("%d", time.Now().UnixNano())

	h.Lock()
	h.store[person.id] = person
	defer h.Unlock()

}

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

func newPersonHandler() *personHandler {
	return &personHandler{
		store: map[string]Person{
			"stego": Person{
				id: "stego",
				name: "alessandro",
				surname: "ianne",
				birthdate: 837095400,
			},
		},
	}
}

func main() {
	admin := newAdminPortal()
	personHandler := newPersonHandler()
	http.HandleFunc("/persons", personHandler.persons)
	http.HandleFunc("/persons/", personHandler.getPerson)
	http.HandleFunc("/admin", admin.handler)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		panic(err)
	}
}