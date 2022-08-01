//Started from https://www.youtube.com/watch?v=bj77B59nkTQ
package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
	"net/http"
	"strconv"
)

//Declare backing object
type book struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

//Create a data set to work with
var books = []book{
	{Id: 1, Title: "Hell Divers", Author: "Nicholas Sansbury Smith", Quantity: 3},
	{Id: 2, Title: "Harry Potter and the Prisoner of Azkaban", Author: "J. K. Rowling", Quantity: 5},
	{Id: 3, Title: "A Clash of Kings", Author: "George R. R. Martin", Quantity: 8},
}

//Helper Methods
func bookById(idAsString string) (*book, error) {
	//Convert to integer or record error in doing so
	id, conversionError := strconv.Atoi(idAsString)

	if conversionError != nil {
		return nil, errors.New("Unable to convert Id")
	}

	for i, existingBook := range books {
		if existingBook.Id == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("Book not found")
}

//Route methods
func checkoutBook(context *gin.Context) {
	idAsString, ok := context.GetQuery("id")

	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Id not valid."})
		return
	}

	book, error := bookById(idAsString)

	if error != nil || book.Quantity <= 0 {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": error.Error()})
		return
	}

	book.Quantity -= 1

	context.IndentedJSON(http.StatusOK, book)
}

func createBook(context *gin.Context) {
	var newBook book

	if error := context.BindJSON(&newBook); error != nil { //Error present when binding JSON
		return //Automagically passes status message as response thanks to BindJSON
	}

	//Add to data set since JSON has been bound
	books = append(books, newBook)

	//Return successful creation
	context.IndentedJSON(http.StatusCreated, newBook)
}

func deleteBook(context *gin.Context) {
	foundBook, error := bookById(context.Param("id")) //Get the book

	if error != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": error.Error()})
		return
	}

	indexToRemove := slices.IndexFunc(books, func(b book) bool { return b.Id == foundBook.Id })
	books = slices.Delete(books, indexToRemove, indexToRemove+1)

	context.IndentedJSON(http.StatusOK, books)
}

func getBookById(context *gin.Context) {
	book, error := bookById(context.Param("id")) //Get the book

	if error != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": error.Error()})
		return
	}

	context.IndentedJSON(http.StatusOK, book)
}

func getBooks(context *gin.Context) {
	//Return all books
	context.IndentedJSON(http.StatusOK, books)
}

func returnBook(context *gin.Context) {
	idAsString, ok := context.GetQuery("id")

	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Id not valid."})
		return
	}

	book, error := bookById(idAsString)

	if error != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": error.Error()})
		return
	}

	book.Quantity += 1

	context.IndentedJSON(http.StatusOK, book)
}

func main() {
	//Setup API and endpoints
	router := gin.Default()
	router.DELETE("/sellAll/:id", deleteBook)
	router.GET("/books", getBooks)
	router.GET("/books/:id", getBookById)
	router.PATCH("/checkout", checkoutBook) //Example of query string
	router.PATCH("/return", returnBook)
	router.POST("/books", createBook)
	router.Run("localhost:8080")
}
