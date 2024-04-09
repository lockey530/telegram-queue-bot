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
// implicit .env arguments: clearData.
func EstablishDBConnection(clearData bool) {
	log.Println("Connecting to database...")
	/*
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatalln("Error loading .env file")
		}*/

	user := os.Getenv("PGUSER")
	dbname := os.Getenv("PGDATABASE")
	password := os.Getenv("PGPASSWORD")
	port := os.Getenv("PGPORT")
	host := os.Getenv("PGHOST")

	log.Println(user)
	log.Println(dbname)
	log.Println(password)
	log.Println(port)

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname)

	if user == "" {
		log.Fatalln("User not provided in .env file.")
	} else if dbname == "" {
		log.Fatalln("Database name not provided in .env file.")
	} else if port == "" {
		log.Fatalln("Port not provided in .env file.")
	}

	// prevent db from being scoped locally - for package use.
	var err error
	db, err = sqlx.Connect(
		"postgres",
		url)
	// fmt.Sprintf("postgres://%s:%s@%s:/%s?sslmode=disable",
	// 	"joshthoo", "2100Isnotaleapyear!", "localhost", port))

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

	initSchemaIfEmpty()

	if clearData {
		removeDBEntries()
	}
}

func initSchemaIfEmpty() {
	var tableExists bool

	err := db.Get(&tableExists, checkTableExistenceQuery, "queue")
	if err != nil {
		log.Fatal("Error retrieving connection information:", err)
	}

	if !tableExists {
		db.MustExec(queueSchema)
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
		log.Println("Schema initiated for admins.")
	} else {
		log.Println("admin schema already initiated.")
	}
}
