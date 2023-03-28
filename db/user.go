package db

import (
	"fmt"

	"github.com/alexanderi96/leafnet/types"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// ValidUser will check if the user exists in db and if exists if the username password
// combination is valid
func ValidUser(email, password string) (bool, error) {
	session := newSession()
	defer session.Close()

	query := fmt.Sprintf(`MATCH (u:User {email: '%s'}) RETURN u.password`, email)

	if res, err := session.Run(query, nil); err != nil {
		return false, err
	} else if res.Next() {
		record := res.Record()
		pwd := record.GetByIndex(0)
		if pwd != nil && pwd.(string) == password {
			return true, nil
		}
	}
	return false, nil
}

func GetUserPasswdHash(email string) (string, error) {
	session := newSession()
	defer session.Close()

	query := fmt.Sprintf(`MATCH (u:User {email: '%s'}) RETURN u.password`, email)

	if res, err := session.Run(query, nil); err != nil {
		return "", err
	} else if res.Next() {
		record := res.Record()
		if pwd, ok := record.GetByIndex(0).(string); ok {
			return pwd, nil
		}
	}
	return "", nil
}

func GetUserInfoByEmail(email string) (types.User, error) {
	session := newSession()
	defer session.Close()

	query := fmt.Sprintf(`MATCH (u:User {email: '%s'}) RETURN u.uuid as uuid, u.creation_date as creation_date, u.last_update as last_update, u.user_name as user_name, u.email as email, u.password as password, u.person as person`, email)

	if res, err := session.Run(query, nil); err != nil {
		return types.User{}, err
	} else if res.Next() {
		return checkRecordAndGetUser(res.Record()), nil
	}

	return types.User{}, nil
}

func GetUserInfoByUserName(user_name string) (types.User, error) {
	session := newSession()
	defer session.Close()

	query := fmt.Sprintf(`MATCH (u:User {user_name: '%s'}) RETURN u.uuid as uuid, u.creation_date as creation_date, u.last_update as last_update, u.user_name as user_name, u.email as email, u.password as password, u.person as person`, user_name)

	if res, err := session.Run(query, nil); err != nil {
		return types.User{}, err
	} else if res.Next() {
		return checkRecordAndGetUser(res.Record()), nil
	}

	return types.User{}, nil
}

func DeleteSelectedUser(email string, password string) error {
	session := newSession()
	defer session.Close()

	query := fmt.Sprintf(`MATCH (u:User {email: '%s', password: '%s'}) DETACH DELETE n`, email, password)

	if _, err := session.Run(query, nil); err != nil {
		return err
	} else {
		return nil
	}
}

// NewPerson crea un nuovo nodo User su neo4j, o lo aggiorna se gia presente
func NewUser(u *types.User) error {
	session := newSession()
	defer session.Close()

	if len(u.Node.UUID) == 0 {
		u.Node.CreateUUID()
	}

	query := `
		MERGE (u:User {uuid: $uuid})
		ON CREATE SET u.user_name = $user_name, u.email = $email, u.password = $password, u.person = $person, u.creation_date = timestamp(), u.last_update = timestamp()
		ON MATCH SET u.user_name = $user_name, u.email = $email, u.password = $password, u.person = $person, u.last_update = timestamp()
	`

	params := map[string]interface{}{
		"uuid":      u.Node.UUID,
		"user_name": u.UserName,
		"email":     u.Email,
		"password":  u.Password,
		"person":    "" + u.Person,
	}

	// Esecuzione della query

	if _, err := session.Run(query, params); err != nil {
		return err
	} else {
		return nil
	}
}

func checkRecordAndGetUser(record neo4j.Record) types.User {

	user := types.User{}

	if uuid, ok := record.Get("uuid"); ok && uuid != nil {
		user.Node.UUID = uuid.(string)
	}

	if creationDate, ok := record.Get("creation_date"); ok && creationDate != nil {
		user.Node.CreationDate = creationDate.(int64)
	}

	if lastUpdate, ok := record.Get("last_update"); ok && lastUpdate != nil {
		user.Node.LastUpdate = lastUpdate.(int64)
	}

	if user_name, ok := record.Get("user_name"); ok && user_name != nil {
		user.UserName = user_name.(string)
	}

	if email, ok := record.Get("email"); ok && email != nil {
		user.Email = email.(string)
	}

	if password, ok := record.Get("password"); ok && password != nil {
		user.Password = password.(string)
	}

	if person, ok := record.Get("person"); ok && person != nil {
		user.Person = person.(string)
	}
	return user
}
