package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

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

	books, err := searchBooks(id, title, author)

	//Some error with fetching the books search
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"books": "null", "error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"books": books, "error": "null"})
}

func checkoutBookAPI(c *gin.Context) {
	id := c.DefaultQuery("id", "")

	if id == "" {
		//No parameter for id is given, return a status indicating a bad request
		c.IndentedJSON(http.StatusBadRequest, gin.H{"book": "null", "error": "No ID Supplied"})
		return
	}

	book, err := checkoutDB(id) //Call the database

	if err != nil {
		//Return an error if the database throws an error
		c.IndentedJSON(http.StatusBadRequest, gin.H{"book": "null", "error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"book": book, "error": "null"})
}

func returnBookAPI(c *gin.Context) {
	id := c.DefaultQuery("id", "")

	if id == "" {
		//No parameter for id is given, return a status indicating a bad request
		c.IndentedJSON(http.StatusBadRequest, gin.H{"book": "null", "error": "No ID Supplied"})
		return
	}

	book, err := returnBookDB(id) //Call the database

	if err != nil {
		//Return an error if the database throws an error
		c.IndentedJSON(http.StatusBadRequest, gin.H{"book": "null", "error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"book": book, "error": "null"})
}

func getRoutes(c *gin.Context) {
	routes := gin.H{
		"GET: /":               "Returns all available routes within the API",
		"POST: /createBook":    "Adds a book with JSON of values (title, author and quantity) to the Database",
		"GET: /searchBooks":    "Returns books matching query ? (title, author or quantity) in the Database",
		"PATCH: /checkoutBook": "Decreases the quantity by one with query ? (id) in the Database",
		"PATCH: /returnBook	":  "Increases the quantity by one with query ? (id) in the Database",
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
	router.GET("/help", getRoutes)

	router.POST("/createBook", createBookAPI)
	router.GET("/searchBooks", searchBooksAPI)
	router.PATCH("/checkoutBook", checkoutBookAPI)
	router.PATCH("/returnBook", returnBookAPI)

	router.Run("localhost:8000")

}
