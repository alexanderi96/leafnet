package db

/*
Stores the database functions related to events like
GetEventsByID(id int)
GetEvents(city string)
DeleteAll()
*/

import (
	"database/sql"
	//"html/template"
	"log"
	//"strconv"
	//"strings"
	//"time"

	_ "github.com/mattn/go-sqlite3" //we want to use sqlite natively
	//md "github.com/shurcooL/github_flavored_markdown"
	//"github.com/alexanderi96/cicerone/types"
)

var database Database
var err error

//Database encapsulates database
type Database struct {
	db *sql.DB
}

//Begins a transaction
func (db Database) begin() (tx *sql.Tx) {
	tx, err := db.db.Begin()
	if err != nil {
		log.Println(err)
		return nil
	}
	return tx
}

func (db Database) prepare(q string) (stmt *sql.Stmt) {
	stmt, err := db.db.Prepare(q)
	if err != nil {
		log.Println(err)
		return nil
	}
	return stmt
}

func (db Database) query(q string, args ...interface{}) (rows *sql.Rows) {
	rows, err := db.db.Query(q, args...)
	if err != nil {
		log.Println(err)
		return nil
	}
	return rows
}

func init() {
	database.db, err = sql.Open("sqlite3", "./cicerone.db")
	if err != nil {
		log.Fatal(err)
	}
}

//Close function closes this database connection
func Close() {
	database.db.Close()
}

//gQuery (genericQuery) encapsulates running multiple queries which don't do much things
func gQuery(sql string, args ...interface{}) error {
	log.Print("inside task query")
	SQL := database.prepare(sql)
	tx := database.begin()
	_, err = tx.Stmt(SQL).Exec(args...)
	if err != nil {
		log.Println("taskQuery: ", err)
		tx.Rollback()
	} else {
		err = tx.Commit()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("Commit successful")
	}
	return err
}
