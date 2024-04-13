package dbaccess

import "log"

func dropDB() {
	// "TRUNCATE queue, admins RESTART IDENTITY;"
	// "DROP TABLE queue, admins;"
	_, err := db.Exec("DROP TABLE queue, admins;")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("DB dropped.")
}
