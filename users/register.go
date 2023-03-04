package users

import (
	"encoding/json"
	"io"
	"net/http"
)

type UserRegistration struct {
	User User `json:"user"`
}

type User struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type UserHandler struct {
	Path           string
	UserRepository UserRepository
}

func (u *UserHandler) Register(w http.ResponseWriter, request *http.Request) {
	requestBody, _ := io.ReadAll(request.Body)
	userRegistrationRequest := UserRegistration{}
	_ = json.Unmarshal(requestBody, &userRegistrationRequest)
	requestUser := userRegistrationRequest.User
	_ = u.UserRepository.RegisterUser(&requestUser)

	w.WriteHeader(201)
	w.Header().Add("Content-Type", "application/json")

	userRegistrationResponse := UserRegistration{
		User: User{
			Username: requestUser.Username,
			Email:    requestUser.Email,
		}}
	bytes, _ := json.Marshal(&userRegistrationResponse)
	_, _ = w.Write(bytes)
}
