package person

import (
	"encoding/json"
	"net/http"
	"fmt"
	"time"
	"sync"
	"strings"
	"math/rand"
	"github.com/alexanderi96/leafnet/utils"
)

type Tree struct {	
	Id string`json:"id"`
	Father *Person`json:"father"`
	Mother *Person`json:"mother"`
}

type Person struct {
	Id 			string`json:"id"`
	Name 		string`json:"name"`
	Surname 	string`json:"surname"`
	Birthdate 	int64 `json:"birthdate"`
	Tree		Tree`json:"tree"`
}

type PersonHandler struct {
	sync.Mutex
	Store map[string]Person
}

func NewPersonHandler() *PersonHandler {
	return &PersonHandler{
		Store: map[string]Person{
			"stego": Person{
				Id: "stego",
				Name: "alessandro",
				Surname: "ianne",
				Birthdate: 837095400,
			},
		},
	}
}

func (h *PersonHandler) Persons(w http.ResponseWriter, r *http.Request) {
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

func (h *PersonHandler) get(w http.ResponseWriter, r *http.Request) {
	persons := make([]Person, len(h.Store))

	h.Lock()
	i := 0
	for _, person := range h.Store {
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

func (h *PersonHandler) getRandomPerson(w http.ResponseWriter, r *http.Request) {
	ids := make([]string, len(h.Store))
	h.Lock()

	i:=0
	for _, id := range h.Store {

		ids[i] = id.Id
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

	w.Header().Add("location", fmt.Sprintf("/persons/%s", target))
	w.WriteHeader(http.StatusFound)

}

func (h *PersonHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
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
	person, ok := h.Store[parts[2]]
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

func (h *PersonHandler) post(w http.ResponseWriter, r *http.Request) {
	var person Person
	err := json.Unmarshal(utils.CheckJsonAndGet(w, r), &person)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	person.Id = fmt.Sprintf("%d", time.Now().UnixNano())

	h.Lock()
	h.Store[person.Id] = person
	defer h.Unlock()

	w.WriteHeader(http.StatusOK)
	jsonBytes, _ := json.Marshal(person.Id)
	w.Write(jsonBytes)
}