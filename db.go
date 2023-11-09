package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func getBooksDB(db *sql.DB) {
	fmt.Println("Getting Books")
}

// Connects to the database and returns the database object
func connectDB() (*sql.DB, error) {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	databaseName := os.Getenv("DB_NAME")

	// Make a connection to the MySQL Database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, databaseName))

	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		log.Print("Connect Sucessfully to " + databaseName)
	}

	return db, nil
}
