package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var username string
var password string
var host string
var port string
var databaseName string

var db *sql.DB

func initDB() error {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		return err
	}

	//Grab the database connection variables from .env file
	var username string = os.Getenv("DB_USERNAME")
	var password string = os.Getenv("DB_PASSWORD")
	var host string = os.Getenv("DB_HOST")
	var port string = os.Getenv("DB_PORT")
	var databaseName string = os.Getenv("DB_NAME")

	//Open the database
	var err error
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, databaseName))

	if err != nil {
		return err
	}

	//Create table if one does not exist for books
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS books (
            id INT AUTO_INCREMENT PRIMARY KEY,
            title VARCHAR(255),
            author VARCHAR(255),
            quantity INT
        )
    `)

	//Check if there is an error creating the table
	if err != nil {
		return nil
	}

	// Set maximum open connections and maximum idle connections
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return nil
}

// Returns all books in the database
func getBooksDB() ([]Book, error) {
	conn, err := db.Begin()

	fmt.Println(err)

	if err != nil {
		return nil, err
	}
	defer conn.Rollback() // Rollback if an error occurs

	rows, err := db.Query("SELECT * FROM books")

	// Check if there was an error fetching the rows
	if err != nil {
		log.Print(err)
		return nil, err
	}

	var books []Book

	//Loop through each row and add the book to an array
	for rows.Next() {
		var book Book

		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Quantity)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
		fmt.Println(books)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := conn.Commit(); err != nil {
		return nil, err
	}

	return books, nil
}

// Adds a book record to the database
func createBookDB(newBook BookInput) (bool, error) {
	conn, err := db.Begin()
	if err != nil {
		return false, err
	}
	defer conn.Rollback() // Rollback if an error occurs

	_, err = db.Exec("INSERT INTO books (title, author, quantity) VALUES (?, ?, ?)", newBook.Title, newBook.Author, newBook.Quantity)

	if err != nil {
		//Check to see if there is an error inserting into the database
		return false, err
	}

	// Commit the transaction
	if err := conn.Commit(); err != nil {
		return false, err
	}

	return true, nil
}
