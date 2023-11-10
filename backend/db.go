package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"

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

// Adds a book record to the database
func createBookDB(newBook BookInput) (bool, error) {
	conn, err := db.Begin()
	if err != nil {
		return false, err
	}
	defer conn.Rollback() // Rollback if an error occurs

	_, err = conn.Exec("INSERT INTO books (title, author, quantity) VALUES (?, ?, ?)", newBook.Title, newBook.Author, newBook.Quantity)

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

// Returns a query and parameters for the given parameters to the function
func sqlSearchBuilder(id int64, title string, author string) (string, []interface{}) {
	var query string
	var params []interface{}

	query = "SELECT * FROM books"

	if id >= 0 {
		query += " WHERE id = ?"
		params = append(params, id)
	}

	if title != "" {
		if query != "SELECT * FROM books" {
			query += " AND "
		} else {
			query += " WHERE "
		}
		query += "title LIKE ?"
		params = append(params, "%"+title+"%")
	}

	if author != "" {
		if query != "SELECT * FROM books" {
			query += " AND "
		} else if title != "" {
			query += " AND "
		} else {
			query += " WHERE "
		}
		query += "author LIKE ?"
		params = append(params, "%"+author+"%")
	}
	return query, params
}

// Optionlly search the database by given parameters -> Returns the array of the books matching the query
func searchBooks(id string, title string, author string) ([]Book, error) {
	idInt, err := strconv.ParseInt(id, 10, 0) //Convert the id to string

	if err != nil {
		return nil, err
	}

	query, params := sqlSearchBuilder(idInt, title, author) //Builds the query parameters automatically

	conn, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer conn.Rollback() // Rollback if an error occurs

	rows, err := conn.Query(query, params...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Quantity)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func checkoutDB(id string) (*Book, error) {
	idInt, err := strconv.ParseInt(id, 10, 0) //Convert the id to string
	if err != nil {
		//Return an error if there was an error with the string
		return nil, err
	}

	//start the connection with the database
	conn, err := db.Begin()
	if err != nil {
		return nil, err
	}

	var (
		updatedQuantity int
		title           string
		author          string
	)

	err = conn.QueryRow("SELECT title, author, quantity FROM books WHERE id = ?", idInt).
		Scan(&title, &author, &updatedQuantity)

	if err != nil {
		return nil, err
	}

	if updatedQuantity > 0 {
		// Decrement the quantity
		_, err = db.Exec("UPDATE books SET quantity = quantity - 1 WHERE id = ?", idInt)

		if err != nil {
			return nil, err
		}

		// Commit the transaction
		if err := conn.Commit(); err != nil {
			conn.Rollback()
			return nil, err
		}

		//Return the updated book
		updatedBook := &Book{
			ID:       id,
			Title:    title,
			Author:   author,
			Quantity: updatedQuantity - 1,
		}

		return updatedBook, nil
	} else {
		return nil, errors.New("No books Available")
	}
}

func returnBookDB(id string) (*Book, error) {
	idInt, err := strconv.ParseInt(id, 10, 0) //Convert the id to string
	if err != nil {
		//Return an error if there was an error with the string
		return nil, err
	}

	//start the connection with the database
	conn, err := db.Begin()
	if err != nil {
		return nil, err
	}

	var (
		updatedQuantity int
		title           string
		author          string
	)

	err = conn.QueryRow("SELECT title, author, quantity FROM books WHERE id = ?", idInt).
		Scan(&title, &author, &updatedQuantity)

	if err != nil {
		return nil, err
	}

	// Increment the quantity
	_, err = db.Exec("UPDATE books SET quantity = quantity + 1 WHERE id = ?", idInt)

	if err != nil {
		return nil, err
	}

	// Commit the transaction
	if err := conn.Commit(); err != nil {
		conn.Rollback()
		return nil, err
	}

	//Return the updated book
	updatedBook := &Book{
		ID:       id,
		Title:    title,
		Author:   author,
		Quantity: updatedQuantity + 1,
	}

	return updatedBook, nil

}
