package users

import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type UserRepository interface {
	RegisterUser(user *User) error
}

type UserNeo4jRepository struct {
	Driver neo4j.Driver
}

func (u *UserNeo4jRepository) RegisterUser(user *User) error {
	panic("implement me")
	//https://www.youtube.com/live/Ar025DqCGew?feature=share&t=2288
}
