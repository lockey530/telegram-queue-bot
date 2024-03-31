package dbaccess

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

// Must be called before any interaction with the DB to initialize the db connection.
func EstablishDBConnection(clearData bool) {
	log.Println("Connecting to database...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	user := os.Getenv("POSTGRES_USER")
	dbname := os.Getenv("POSTGRES_DBNAME")
	port := os.Getenv("POSTGRES_PORT")

	log.Println(user)
	log.Println(dbname)
	log.Println(port)

	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		"postgres",
		"postgres",
		"postgres",
		"5432",
		"postgres")

	if user == "" {
		log.Fatalln("User not provided in .env file.")
	} else if dbname == "" {
		log.Fatalln("Database name not provided in .env file.")
	} else if port == "" {
		log.Fatalln("Port not provided in .env file.")
	}

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

	initSchemaIfEmpty(db)

	if clearData {
		removeDBEntries(db)
	}
}

func initSchemaIfEmpty(db *sqlx.DB) {
	var tableExists bool

	err := db.Get(&tableExists, checkExistenceQuery, "queue")
	if err != nil {
		log.Fatal("Error retrieving connection information:", err)
	}

	if !tableExists {
		db.MustExec(queueSchema)
		log.Println("schema initiated for the queue.")
		db.MustExec(adminSchema)
		log.Println("Schema initiated for admins.")
	} else {
		log.Println("users schema already initiated.")
	}

	err = db.Get(&tableExists, checkExistenceQuery, "admins")
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
