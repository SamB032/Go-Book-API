package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Grab all books from the database and return them
func getBooksAPI(c *gin.Context) {
	var books []Book

	books, err := getBooksDB()

	if err != nil {
		//Return the error to the api if it cant connect
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	}
	c.IndentedJSON(http.StatusOK, books)
}

func createBookAPI(c *gin.Context) {
	var newBook BookInput

	//Return if an error is caused by binding the input to the json
	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"added": "false", "bookAdded": "null", "error": err.Error()})
		return
	}

	added, err := createBookDB(newBook)

	if !added || err != nil {
		//Return an error if we cant add it
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"added": "false", "bookAdded": "null", "error": err.Error()})
	}

	//Return routes within api
	c.IndentedJSON(http.StatusCreated, gin.H{"Added": "true", "bookAdded": newBook, "error": "null"})
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

func getRoutes(c *gin.Context) {
	routes := gin.H{
		"GET: /":            "Returns all available routes within the API",
		"GET: /book":        "Returns all books in Database",
		"POST: /createBook": "Adds a book with JSON of values (title, author and quantity) to the Database",
		"GET: /searchBooks": "Returns books matching query ? (title, author or quantity) in the Database",
		"PATCH: /checkout":  "Decreases the quantity by one with query ? (id) in the Database",
		"PATCH: /checkin":   "Increases the quantity by one with query ? (id) in the Database",
	}

	c.IndentedJSON(http.StatusOK, routes)
}

func main() {
	var err error = initDB()

	if err != nil {
		//Database has failed to connect, end the API
		log.Fatal(err)
		return
	}

	router := gin.Default()

	router.GET("/", getRoutes)
	router.GET("/books", getBooksAPI)
	router.POST("/createBook", createBookAPI)
	router.GET("/searchBooks", searchBooksAPI)
	router.PATCH("/checkout", checkoutBookAPI)
	router.PATCH("/returnBook", getRoutes)

	router.Run("localhost:8000")

}
