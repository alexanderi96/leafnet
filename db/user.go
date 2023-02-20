package db

import (
	"log"
	"fmt"

	"github.com/alexanderi96/leafnet/types"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

//ValidUser will check if the user exists in db and if exists if the username password
//combination is valid
func ValidUser(email, password string) bool {
	session := newSession()
	defer session.Close()

	query := fmt.Sprintf(`MATCH (u:User {email: '%s'}) RETURN u.password`, email)
	result, err := session.Run(query, nil)
	
	if err != nil {
		log.Println(err)
		return false
	}

	if result.Next() {
		record := result.Record()
		pwd := record.GetByIndex(0)
		
		if pwd != nil && pwd.(string)==password {
			return true
		}
	}

	return false
}

func GetUserInfo(email string) (u types.User, e error){
		session := newSession()
	defer session.Close()

	query := fmt.Sprintf(`MATCH (u:User {email: '%s'}) RETURN u.uuid, u.creation_date, u.last_update, u.user_name, u.email, u.password, u.person`, email)
	result, err := session.Run(query, nil)
	if err != nil {
		log.Println(err)
		return types.User{}, err
	}

	if result.Next() {
		log.Println(result)
		return *checkRecordAndGetUser(result.Record()), nil
	}

	return types.User{}, nil
}

func DeleteSelectedUser(email string, password string) error{
	session := newSession()
	defer session.Close()

	query := fmt.Sprintf(`MATCH (u:User {email: '%s', password: '%s'}) DETACH DELETE n`, email, password)

	_, err := session.Run(query, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// NewPerson crea un nuovo nodo User su neo4j, o lo aggiorna se gia presente
func NewUser(u *types.User) error {
	session := newSession()
	defer session.Close()
	
	if len(u.Node.UUID) == 0 {
		u.Node.CreateUUID()
	}

	query := fmt.Sprintf(`
		MERGE (u:User {uuid: $uuid})
		ON CREATE SET u.user_name = $user_name, u.email = $email, u.password = $password, u.person = $person, u.creation_date = timestamp(), u.last_update = timestamp()
		ON MATCH SET u.user_name = $user_name, u.email = $email, u.password = $password, u.person = $person, u.last_update = timestamp()
	`)

	params := map[string]interface{}{
		"uuid":        u.Node.UUID,
		"user_name": u.UserName,
		"email":  u.Email,
		"password": u.Password,
		"person":   "" + u.Person,
	}

	log.Println("before execute NewUser")
	
	// Esecuzione della query
	result, err := session.Run(query, params)
	if err != nil {
		log.Println(err)
		return err
	}
	
	// Risultato della query
	if result.Err() == nil {
		log.Println(result)
		return nil
	} else {
		log.Println(result.Err())
		return result.Err()
	}
}

func checkRecordAndGetUser(record neo4j.Record) (user *types.User) {

	user = &types.User{}
	
	uuid := record.GetByIndex(0)

	if uuid != nil {
		user.UUID = uuid.(string)
	}

	creationDate := record.GetByIndex(1)

	if creationDate != nil {
		user.Node.CreationDate = creationDate.(int64)
	}

	lastUpdate := record.GetByIndex(2)

	if lastUpdate != nil {
		user.Node.LastUpdate = lastUpdate.(int64)
	}

	user_name := record.GetByIndex(3)
	
	if user_name != nil {
		user.UserName = user_name.(string)
	}

	email := record.GetByIndex(4)
	
	if email != nil {
		user.Email = email.(string)
	}

	password := record.GetByIndex(5)
	
	if password != nil {
		user.Password = password.(string)
	}

	person := record.GetByIndex(6)

	if person != nil {
		user.Person = person.(string)
	}

	return user
}