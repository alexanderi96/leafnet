package db

import (
	"fmt"
	"log"

	"github.com/alexanderi96/leafnet/types"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

var driver neo4j.Driver

func init() {
	// Connessione al database
	log.Println("Initiating db")

	var err error
	driver, err = neo4j.NewDriver("bolt://192.168.1.157:7687/leafnet/", neo4j.BasicAuth("neo4j", ".M4nD0rl423", ""), func(c *neo4j.Config) { c.Encrypted = false })
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

	query := fmt.Sprintf(`
		OPTIONAL MATCH (form:Person {uuid: $form})
		OPTIONAL MATCH (to:Person {uuid: $to})
		OPTIONAL MATCH (form)-[r:$relation]->(to)
		DELETE r
		`)

	params := map[string]interface{}{
		"form":     from,
		"relation": relation,
		"to":       to,
	}

	_, err = session.Run(query, params)
	if err != nil {
		log.Println("Error running query: %s\n", err)
	}

	return
}

// NewPerson crea un nuovo nodo Person su neo4j, o lo aggiorna se gia presente
func NewPerson(p *types.Person) error {
	session := newSession()
	defer session.Close()

	if len(p.UUID) == 0 {
		p.CreateUUID()
	}

	query := fmt.Sprintf(`
		MERGE (p:Person {uuid: $uuid})
		WITH p

		OPTIONAL MATCH (p1:Person {uuid: p.parent1})
		OPTIONAL MATCH (p2:Person {uuid: p.parent2})
		OPTIONAL MATCH (p1)-[r1:PARENT_OF]->(p)
		OPTIONAL MATCH (p2)-[r2:PARENT_OF]->(p)
		DELETE r1, r2

		SET p.first_name = $first_name, p.last_name = $last_name, p.birth_date = $birth_date, p.death_date = $death_date, p.parent1 = $parent1, p.parent2 = $parent2, p.bio = $bio, p.last_update = timestamp()

		FOREACH (g1 IN CASE WHEN $parent1 <> "" THEN [1] ELSE [] END |
		  MERGE (p1:Person {uuid: $parent1})
		  MERGE (p1)-[:PARENT_OF]->(p)
		)

		FOREACH (g2 IN CASE WHEN $parent2 <> "" THEN [1] ELSE [] END |
		  MERGE (p2:Person {uuid: $parent2})
		  MERGE (p2)-[:PARENT_OF]->(p)
		)
	`)

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

	log.Println("before execute NewPerson")

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
		return types.Person{}, fmt.Errorf("Person not found")
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
		log.Println("Error running query: %s\n", err)
	}

	return
}

func checkRecordAndGetPerson(record neo4j.Record) (person types.Person) {

	person = types.Person{}

	uuid := record.GetByIndex(0)

	if uuid != nil {
		person.Node.UUID = uuid.(string)
	}

	creationDate := record.GetByIndex(1)

	if creationDate != nil {
		person.Node.CreationDate = creationDate.(int64)
	}

	lastUpdate := record.GetByIndex(2)

	if lastUpdate != nil {
		person.Node.LastUpdate = lastUpdate.(int64)
	}

	firstName := record.GetByIndex(3)

	if firstName != nil {
		person.FirstName = firstName.(string)
	}

	lastName := record.GetByIndex(4)

	if lastName != nil {
		person.LastName = lastName.(string)
	}

	birthDate := record.GetByIndex(5)

	if birthDate != nil {
		person.BirthDate = birthDate.(int64)
	}

	deathDate := record.GetByIndex(6)

	if deathDate != nil {
		person.DeathDate = deathDate.(int64)
	}

	parent1 := record.GetByIndex(7)

	if parent1 != nil {
		person.Parent1 = parent1.(string)
	}

	parent2 := record.GetByIndex(8)

	if parent2 != nil {
		person.Parent2 = parent2.(string)
	}

	bio := record.GetByIndex(9)

	if bio != nil {
		person.Bio = bio.(string)
	}

	return person
}
