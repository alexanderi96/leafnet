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
	log.Println("Initiating session")
	session, err := driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		log.Fatalf("Error creating session: %v", err)
	}
	log.Println("Session initiated")
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

		SET p.first_name = $first_name, p.last_name = $last_name, p.birth_date = $birth_date, p.death_date = $death_date, p.parent1 = $parent1, p.parent2 = $parent2, p.bio = $bio, p.last_update = timestamp()
	`

	if len(p.UUID) == 0 {
		p.CreateUUID()
		query += `, p.creation_date = timestamp()`
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

	query := fmt.Sprintf(`MATCH (p:Person {uuid: '%s'}) RETURN p.uuid, p.creation_date as creation_date, p.last_update as last_update, p.first_name, p.last_name, p.birth_date, p.death_date, p.parent1, p.parent2, p.bio`, uuid)
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

	result, err := session.Run(`MATCH (p:Person) RETURN p.uuid as uuid, p.creation_date as creation_date, p.last_update as last_update, p.first_name as first_name, p.last_name as last_name, p.birth_date as birth_date, p.death_date as death_date, p.parent1 as parent1, p.parent2 as parent2, p.bio as bio`, nil)
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

func checkRecordAndGetPerson(record neo4j.Record) (person types.Person) {

	person = types.Person{}

	if uuid, ok := record.GetByIndex(0).(string); ok {
		person.Node.UUID = uuid
	}

	if creationDate, ok := record.GetByIndex(1).(int64); ok {
		person.Node.CreationDate = creationDate
	}

	if lastUpdate, ok := record.GetByIndex(2).(int64); ok {
		person.Node.LastUpdate = lastUpdate
	}

	if firstName, ok := record.GetByIndex(3).(string); ok {
		person.FirstName = firstName
	}

	if lastName, ok := record.GetByIndex(4).(string); ok {
		person.LastName = lastName
	}

	if birthDate, ok := record.GetByIndex(5).(int64); ok {
		person.BirthDate = birthDate
	}

	if deathDate, ok := record.GetByIndex(6).(int64); ok {
		person.DeathDate = deathDate
	}

	if parent1, ok := record.GetByIndex(7).(string); ok {
		person.Parent1 = parent1
	}

	if parent2, ok := record.GetByIndex(8).(string); ok {
		person.Parent2 = parent2
	}

	if bio, ok := record.GetByIndex(9).(string); ok {
		person.Bio = bio
	}

	return person
}
