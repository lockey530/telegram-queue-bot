package dbaccess

import "log"

func removeDBEntries() {
	// "TRUNCATE queue, admins RESTART IDENTITY;"
	// "DROP TABLE queue, admins;"
	_, err := db.Exec("TRUNCATE queue, admins RESTART IDENTITY;")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Entries wiped.")
}
