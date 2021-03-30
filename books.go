package main
import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/newrelic/go-agent/v3/newrelic"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)
// Book struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}
// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
// Init books var as a slice Book struct
var books []Book
// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}
// Get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // Gets params
	// Loop through books and find one with the id from the params
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}
// Add new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(100000000)) // Mock ID - not safe
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}
// Update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}
// Delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}
// Main function
func main() {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("learning-newrelic"),
		newrelic.ConfigLicense("734c6e76737384effd95dd3fd19f055928a0NRAL"),
		newrelic.ConfigDistributedTracerEnabled(true),
	)
	if err != nil {
		panic(err)
	}
	// Init router
	r := mux.NewRouter()
	// Hardcoded data - @todo: add database
	books = append(books, Book{ID: "1", Isbn: "438227", Title: "Book One", Author: &Author{Firstname: "John", Lastname: "Doe"}})
	books = append(books, Book{ID: "2", Isbn: "454555", Title: "Book Two", Author: &Author{Firstname: "Steve", Lastname: "Smith"}})
	// Route handles & endpoints
	txn := app.StartTransaction("transaction_name")
	defer txn.End()
	segment := txn.StartSegment("mySegmentName")
	// ... code you want to time here ...
	r.HandleFunc(newrelic.WrapHandleFunc(app,"/books", createBook)).Methods("POST")
	r.HandleFunc(newrelic.WrapHandleFunc(app,"/books", getBooks)).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc(newrelic.WrapHandleFunc(app,"/books/{id}", updateBook)).Methods("PUT")
	r.HandleFunc(newrelic.WrapHandleFunc(app,"/books/{id}", deleteBook)).Methods("DELETE")
	segment.End()
	// Start server
	log.Fatal(http.ListenAndServe(":8000", r))
}
// Request sample
// {
//     "isbn":"4545454",
//     "title":"Book Three",
//     "author":{"firstname":"Harry","lastname":"White"}
// }