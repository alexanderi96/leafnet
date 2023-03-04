package users_test

import (
	"io"
	"net/http/httptest"
	"strings"

	"github.com/alexanderi96/golang-react-jwt/users"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type FakeUserRepository struct {
}

func (f *FakeUserRepository) RegisterUser(user *users.User) error {
	return nil
}

var _ = Describe("Users", func() {
	It("should register", func() {
		handler := users.UserHandler{
			Path:           "/users",
			UserRepository: &FakeUserRepository{},
		}

		// {"user":{"email":"{{EMAIL}}","password":"{{PASSWORD}}", "username":"{{USERNAME}}"}}
		testResponseWriter := httptest.NewRecorder()
		requestBody := strings.NewReader("{\"user\":{\"email\":\"asd@asd.asd\",\"password\":\"asd\", \"username\":\"asd\"}}")
		handler.Register(testResponseWriter, httptest.NewRequest("POST", "/users", requestBody))

		Expect(testResponseWriter.Code).To(Equal(200))
		responseBody, _ := io.ReadAll(testResponseWriter.Body)
		Expect(string(responseBody)).To(Equal("{\"user\":{\"email\":\"asd@asd.asd\",\"password\":\"asd\", \"username\":\"asd\"}}"))

	})
})
