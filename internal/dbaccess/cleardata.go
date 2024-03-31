package dbaccess

import "log"

func removeDBEntries() {
	if _, err := db.Exec("TRUNCATE queue, admins"); err != nil {
		log.Fatal(err)
	}
	log.Println("Entries wiped.")
}
