package store

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var createDBQuery = `CREATE TABLE tokens (
	"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
	"token" TEXT NOT NULL
);`

func createStoreTables(db *sql.DB) {
	// Create token table
	log.Println("Creating token table")
	statement, err := db.Prepare(createDBQuery)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
}

func CreateOrGetDB(filename string) (*sql.DB, error) {
	dbExists := true
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			dbExists = false
			_, err := os.Create(filename)
			if err != nil {
				return nil, err
			}
		}
	}

	DB, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}

	if !dbExists {
		log.Println("DB is new, creating store tables")
		createStoreTables(DB)
	} else {
		log.Println("DB is not new, skipping table creation")
	}

	return DB, nil
}
