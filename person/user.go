package person

import (
	"encoding/json"
	"net/http"
	"fmt"
	"time"
	"sync"
	"strings"
	"io/ioutil"
)

type User struct {
	Id string
	Username string
	Password string
	Email string
	Person Person
}

type UserHandler struct {
	sync.Mutex
	Store map[string]User
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		Store: map[string]User{
			"stego": User{
				Id: "stego",
				Username: "stego",
				Password: "ciao",
				Email: "alessandro.ianne96@gmail.com",
			},
		},
	}
}

func (h *UserHandler) Users(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) get(w http.ResponseWriter, r *http.Request) {
	users := make([]User, len(h.Store))

	h.Lock()
	i := 0
	for _, user := range h.Store {
		users[i] = user
		i++
	}
	h.Unlock()


	jsonBytes, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	h.Lock()
	user, ok := h.Store[parts[2]]
	h.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *UserHandler) post(w http.ResponseWriter, r *http.Request) {
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

	var user User
	err = json.Unmarshal(bodyBytes, &user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	user.Id = fmt.Sprintf("%d", time.Now().UnixNano())

	h.Lock()
	h.Store[user.Id] = user
	defer h.Unlock()

	w.WriteHeader(http.StatusOK)
	jsonBytes, _ := json.Marshal(user.Id)
	w.Write(jsonBytes)
}
