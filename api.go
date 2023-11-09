package main

import (
	// "errors"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func getBooksAPI(c *gin.Context) {
	getBooksDB(db)
	c.IndentedJSON(http.StatusOK, gin.H{"Message": "Returns all Books"})
}

func createBookAPI(c *gin.Context) {
	c.IndentedJSON(http.StatusCreated, gin.H{"Message": "Create a Book"})
}

func searchBooksAPI(c *gin.Context) {
	id := c.DefaultQuery("id", "-1")
	title := c.DefaultQuery("title", "")
	author := c.DefaultQuery("author", "")

	c.IndentedJSON(http.StatusOK, gin.H{"id": id, "title": title, "author": author})
}

func checkoutBookAPI(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"Message": "Checkout Book"})
}

func returnBookAPI(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"Message": "Return a book"})
}

func startAPI() {
	router := gin.Default()

	router.GET("/books", getBooksAPI)
	router.POST("/createBook", createBookAPI)
	router.GET("/searchBooks", searchBooksAPI)
	router.PATCH("/checkout", checkoutBookAPI)
	router.PATCH("/returnBook", returnBookAPI)

	router.Run("localhost:8000")

	db, err = connectDB()

	//Database has failed to connect, end the API
	if err != nil {
		return
	}
}
