package types

import (
	"github.com/google/uuid"
)

type Node struct {
	UUID         string `json:"uuid"`
	CreationDate int64  `json:"creation_date"`
	LastUpdate   int64  `json:"last_update"`
	Owner        string `json:"owner"` // Defines the user who created the node
}

type Person struct {
	Node
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BirthDate int64  `json:"birth_date,omitempty"`
	DeathDate int64  `json:"death_date,omitempty"`
	Parent1   string `json:"parent1,omitempty"`
	Parent2   string `json:"parent2,omitempty"`
	//Events   []Multimedia `json:"events,omitempty"`
	Bio string `json:"bio,omitempty"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        int64  `json:"date"`
	Location    string `json:"location"`
}

type Multimedia struct {
	Node
	Info
	Attachment []byte
}

func (n *Node) CreateUUID() {
	if len(n.UUID) == 0 {
		n.UUID = uuid.New().String()
	}
}
