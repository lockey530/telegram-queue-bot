package dbaccess

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func removeDBEntries(db *sqlx.DB) {
	if _, err := db.Exec(wipeData); err != nil {
		log.Fatal(err)
	}
	log.Println("Entries wiped.")
}
