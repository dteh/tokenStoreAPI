package store

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func insertTokens(tokens []string, db *sql.DB) {
	// return if no tokens
	if len(tokens) == 0 {
		return
	}
	log.Println("Inserting", len(tokens), "tokens")
	statement := `INSERT INTO tokens (token) VALUES `
	for i, token := range tokens {
		statement += "('" + token + "')"
		if i != len(tokens)-1 {
			statement += ","
		}
	}
	log.Println(statement)
	res, err := db.Exec(statement)
	if err != nil {
		log.Println("Error inserting tokens into db", err)
	}
	log.Println(res.RowsAffected())
}

type tokenEntry struct {
	id    int
	token string
}

func getSingleToken(db *sql.DB) (string, error) {
	rows, err := db.Query(`SELECT * FROM tokens LIMIT 1`)
	if err != nil {
		return "", err
	}
	hasNext := rows.Next()
	if !hasNext {
		return "", errors.New("no tokens found in db")
	}
	token := &tokenEntry{}
	err = rows.Scan(&token.id, &token.token)
	if err != nil {
		return "", err
	}
	rows.Close()

	idToDelete := token.id
	_, err = db.Exec(`DELETE FROM tokens WHERE id=` + fmt.Sprint(idToDelete))
	if err != nil {
		return "", err
	}

	return token.token, nil
}
