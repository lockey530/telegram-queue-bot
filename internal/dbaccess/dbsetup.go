package dbaccess

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

// Must be called before any interaction with the DB to initialize the db connection.
func EstablishDBConnection(toReset bool) {
	log.Println("Connecting to database...")
	/*
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalln("Error loading .env file")
		}*/

	user := os.Getenv("POSTGRES_USER")
	dbname := os.Getenv("POSTGRES_DB")
	password := os.Getenv("POSTGRES_PASSWORD")
	port := os.Getenv("PGPORT")
	host := os.Getenv("PGHOST")

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname)

	// prevent db from being scoped locally - for package use.
	var err error
	db, err = sqlx.Connect(
		"postgres",
		url)

	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Connected to the database")

	var databaseName, userName string
	err = db.QueryRow("SELECT current_database(), current_user").Scan(&databaseName, &userName)
	if err != nil {
		log.Fatal("Error retrieving connection information:", err)
	}

	log.Printf("Database Name: %s\n", databaseName)
	log.Printf("User Name: %s\n", userName)

	if toReset {
		_, err := db.Exec(`
			DROP TABLE IF EXISTS 
				queue, admins;
		`)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("DB dropped.")
	}
	initSchemaIfEmpty()
}

func initSchemaIfEmpty() {
	var tableExists bool

	err := db.Get(&tableExists, checkTableExistenceQuery, "queue")
	if err != nil {
		log.Fatal("Error retrieving connection information:", err)
	}

	if !tableExists {
		db.MustExec(queueSchema)
		db.MustExec(indexQueueSchema)
		log.Println("schema initiated for the queue.")
	} else {
		log.Println("queue schema already initiated.")
	}

	err = db.Get(&tableExists, checkTableExistenceQuery, "admins")
	if err != nil {
		log.Fatal("Error retrieving connection information:", err)
	}

	if !tableExists {
		db.MustExec(adminSchema)
		log.Println("schema initiated for admins.")
	} else {
		log.Println("admin schema already initiated.")
	}
}
