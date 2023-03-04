package users_test

import (
	"context"
	"fmt"
	"io"

	"github.com/alexanderi96/golang-react-jwt/users"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var _ = Describe("User repository", func() {

	const username = "neo4j"
	const password = "neo4j"

	var ctx context.Context
	var neo4jContainer testcontainers.Container

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		neo4jContainer, err = startContainer(ctx, username, password)
		Expect(err).To(BeNil(), "Container should start")
	})

	AfterEach(func() {
		Expect(neo4jContainer.Terminate(ctx)).To(BeNil(), "Container should be terminated")
	})

	FIt("should register", func() {
		port, err := neo4jContainer.MappedPort(ctx, "7687")
		Expect(err).To(BeNil(), "Port should be resolved")
		address := fmt.Sprintf("bolt://localhost:%d", port.Int())
		driver, err := neo4j.NewDriver(address, neo4j.BasicAuth(username, password, ""))
		Expect(err).To(BeNil(), "Driver should be created")
		defer Close(driver, "Driver")
		repository := &users.UserNeo4jRepository{
			Driver: driver,
		}

		username := "asd"
		email := "asd@asd.asd"
		err = repository.RegisterUser(&users.User{
			Username: username,
			Email:    email,
			Password: "asd",
		})
		Expect(err).To(BeNil(), "User should be registered")

		session := driver.NewSession(neo4j.SessionConfig{})
		defer Close(session, "Session")
		result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			res, err := tx.Run(
				"MATCH (u:User{username:$username, email:$email}) RETURN u.username as username, u.email as email, u.password as password",
				map[string]interface{}{
					"username": username,
					"email":    email,
				})
			if err != nil {
				return nil, err
			}
			singleRecord, err := res.Single()
			if err != nil {
				return nil, err
			}
			return &users.User{
				Username: singleRecord.Values[0].(string),
				Email:    singleRecord.Values[1].(string),
				Password: singleRecord.Values[2].(string),
			}, nil
		})
		Expect(err).To(BeNil(), "Transaction should be successful")
		persistedUser := result.(*users.User)
		Expect(persistedUser.Username).To(Equal(username))
		Expect(persistedUser.Email).To(Equal(email))
		Expect(persistedUser.Password).NotTo(BeNil())
	})
})

func Close(closer io.Closer, resourceName string) {
	Expect(closer.Close()).To(BeNil(), fmt.Sprintf("%s should close", resourceName))
}

func startContainer(ctx context.Context, username, password string) (testcontainers.Container, error) {
	request := testcontainers.ContainerRequest{
		Image:        "neo4j",
		ExposedPorts: []string{"7687/tcp"},
		Env: map[string]string{
			"NEO4J_AUTH": fmt.Sprintf("#{username}/#{password}"),
		},
		WaitingFor: wait.ForLog("Bolt enabled"),
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
}
