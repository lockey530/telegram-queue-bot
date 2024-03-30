package db

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func removeDBEntries(db *sqlx.DB) {
	_, err := db.Exec(wipeData)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Entries wiped.")
}
