package db

import (
	"fmt"
	"log"

	"github.com/alexanderi96/leafnet/config"
	"github.com/alexanderi96/leafnet/types"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

var driver neo4j.Driver

func init() {
	// Connessione al database
	log.Println("Initiating db")

	var err error
	driver, err = neo4j.NewDriver("bolt://"+config.Config["neo4j_endpoint"]+":"+config.Config["neo4j_port"]+"/"+config.Config["neo4j_schema"]+"/", neo4j.BasicAuth(config.Config["neo4j_username"], config.Config["neo4j_password"], ""), func(c *neo4j.Config) { c.Encrypted = false })
	if err != nil {
		log.Fatalf("Error creating driver: %v", err)
	}
	//defer driver.Close()
}

func newSession() neo4j.Session {
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}
	//defer session.Close()
	return session
}

func DeleteRelation(from, relation, to string) (err error) {
	session := newSession()
	defer session.Close()

	query := `
		OPTIONAL MATCH (form:Person {uuid: $form})
		OPTIONAL MATCH (to:Person {uuid: $to})
		OPTIONAL MATCH (form)-[r:$relation]->(to)
		DELETE r
		`

	params := map[string]interface{}{
		"form":     from,
		"relation": relation,
		"to":       to,
	}

	_, err = session.Run(query, params)
	if err != nil {
		log.Println("Error running query: ", err)
	}

	return
}

// ManagePerson crea un nuovo nodo Person su neo4j, o lo aggiorna se gia presente
func ManagePerson(p *types.Person) error {
	session := newSession()
	defer session.Close()

	query := `
		MERGE (p:Person {uuid: $uuid})
		WITH p

		OPTIONAL MATCH (p1:Person {uuid: p.parent1})
		OPTIONAL MATCH (p2:Person {uuid: p.parent2})
		OPTIONAL MATCH (p1)-[r1:PARENT_OF]->(p)
		OPTIONAL MATCH (p2)-[r2:PARENT_OF]->(p)
		DELETE r1, r2

		SET p.last_update = timestamp(), p.first_name = $first_name, p.last_name = $last_name, p.birth_date = $birth_date, p.death_date = $death_date, p.parent1 = $parent1, p.parent2 = $parent2, p.bio = $bio
	`

	if len(p.UUID) == 0 {
		p.CreateUUID()
		query += `, p.creation_date = timestamp(), p.owner = $owner`
	}

	query += `

		FOREACH (g1 IN CASE WHEN $parent1 <> "" THEN [1] ELSE [] END |
		  MERGE (p1:Person {uuid: $parent1})
		  MERGE (p1)-[:PARENT_OF]->(p)
		)

		FOREACH (g2 IN CASE WHEN $parent2 <> "" THEN [1] ELSE [] END |
		  MERGE (p2:Person {uuid: $parent2})
		  MERGE (p2)-[:PARENT_OF]->(p)
		)
	`

	params := map[string]interface{}{
		"uuid":       p.Node.UUID,
		"owner":      p.Node.Owner,
		"first_name": p.FirstName,
		"last_name":  p.LastName,
		"birth_date": p.BirthDate,
		"death_date": p.DeathDate,
		"parent1":    "" + p.Parent1,
		"parent2":    "" + p.Parent2,
		"bio":        p.Bio,
	}

	// Esecuzione della query
	result, err := session.Run(query, params)
	if err != nil {
		log.Println(err)
		return err
	}

	// Risultato della query
	if result.Err() == nil {
		return nil
	} else {
		return result.Err()
	}
}

// GetPerson recupera una persona dal database
func GetPerson(uuid string) (types.Person, error) {
	session := newSession()
	defer session.Close()

	query := fmt.Sprintf(`MATCH (p:Person {uuid: '%s'}) RETURN p.uuid as uuid, p.creation_date as creation_date, p.last_update as last_update, p.owner as owner, p.first_name as first_name, p.last_name as last_name, p.birth_date as birth_date, p.death_date as death_date, p.parent1 as parent1, p.parent2 as parent2, p.bio as bio`, uuid)
	result, err := session.Run(query, nil)

	if err != nil {
		return types.Person{}, err
	}

	if result.Next() {
		return checkRecordAndGetPerson(result.Record()), nil
	} else {
		return types.Person{}, fmt.Errorf("person not found")
	}
}

func GetPersons() []types.Person {
	session := newSession()
	defer session.Close()

	result, err := session.Run(`MATCH (p:Person) RETURN p.uuid as uuid, p.creation_date as creation_date, p.last_update as last_update, p.owner as owner, p.first_name as first_name, p.last_name as last_name, p.birth_date as birth_date, p.death_date as death_date, p.parent1 as parent1, p.parent2 as parent2, p.bio as bio`, nil)
	if err != nil {
		log.Fatalf("Error running query: %s\n", err)
	}

	persons := []types.Person{}

	for result.Next() {
		record := result.Record()

		persons = append(persons, checkRecordAndGetPerson(record))
	}

	return persons
}

func DeletePerson(uuid string) (err error) {
	session := newSession()
	defer session.Close()

	_, err = session.Run(fmt.Sprintf(`MATCH (p:Person{uuid: '%s'}) DETACH DELETE p`, uuid), nil)
	if err != nil {
		log.Println("Error running query: ", err)
	}

	return
}

func checkRecordAndGetPerson(record neo4j.Record) types.Person {

	person := types.Person{}

	if uuid, ok := record.Get("uuid"); ok && uuid != nil {
		person.Node.UUID = uuid.(string)
	}

	if creationDate, ok := record.Get("creation_date"); ok && creationDate != nil {
		person.Node.CreationDate = creationDate.(int64)
	}

	if lastUpdate, ok := record.Get("last_update"); ok && lastUpdate != nil {
		person.Node.LastUpdate = lastUpdate.(int64)
	}

	if owner, ok := record.Get("owner"); ok && owner != nil {
		person.Node.Owner = owner.(string)
	}

	if firstName, ok := record.Get("first_name"); ok && firstName != nil {
		person.FirstName = firstName.(string)
	}

	if lastName, ok := record.Get("last_name"); ok && lastName != nil {
		person.LastName = lastName.(string)
	}

	if birthDate, ok := record.Get("birth_date"); ok && birthDate != nil {
		person.BirthDate = birthDate.(int64)
	}

	if deathDate, ok := record.Get("death_date"); ok && deathDate != nil {
		person.DeathDate = deathDate.(int64)
	}

	if parent1, ok := record.Get("parent1"); ok && parent1 != nil {
		person.Parent1 = parent1.(string)
	}

	if parent2, ok := record.Get("parent2"); ok && parent2 != nil {
		person.Parent2 = parent2.(string)
	}

	if bio, ok := record.Get("bio"); ok && bio != nil {
		person.Bio = bio.(string)
	}

	return person
}
