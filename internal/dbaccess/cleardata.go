package dbaccess

import "log"

func removeDBEntries() {
	_, err := db.Exec("TRUNCATE queue, admins RESTART IDENTITY;")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Entries wiped.")
}
