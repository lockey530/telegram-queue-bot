package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func EstablishDBConnection() {
	log.Println("Connecting to database...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}

	user := os.Getenv("POSTGRES_USER")
	dbname := os.Getenv("POSTGRES_DBNAME")
	port := os.Getenv("POSTGRES_PORT")

	if user == "" {
		log.Fatalln("User not provided in .env file.")
	} else if dbname == "" {
		log.Fatalln("Database name not provided in .env file.")
	} else if port == "" {
		log.Fatalln("Port not provided in .env file.")
	}
	// password := os.Getenv("POSTGRES_PASSWORD")
	db, err := sqlx.Connect(
		"postgres",
		fmt.Sprintf("user=%s dbname=%s port=%s sslmode=disable",
			user, dbname, port))
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

	fmt.Printf("Database Name: %s\n", databaseName)
	fmt.Printf("User Name: %s\n", userName)
}
