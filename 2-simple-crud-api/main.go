package main

import (
	"encoding/json"
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type Book struct {
	ID     uuid.UUID `json:"id"`
	Isbn   string    `json:"isbn"`
	Title  string    `json:"title"`
	Author *Author   `json:"author"`
}

type Author struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var books []Book

func getBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		responseJSON(w, http.StatusOK, books)
		return
	}
	responseJSON(w, http.StatusMethodNotAllowed, nil)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		id, _ := uuid.FromString(r.URL.Query().Get("id"))

		for _, item := range books {
			if item.ID == id {
				responseJSON(w, http.StatusOK, item)
				return
			}
		}
		responseJSON(w, http.StatusNotFound, nil)
		return
	}
	responseJSON(w, http.StatusMethodNotAllowed, nil)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var book Book
		_ = json.NewDecoder(r.Body).Decode(&book)
		book.ID = uuid.NewV4()
		books = append(books, book)
		responseJSON(w, http.StatusCreated, book)
		return
	}
	responseJSON(w, http.StatusMethodNotAllowed, nil)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(r.URL.Query().Get("id"))

	for index, item := range books {
		if item.ID == id {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = id
			books = append(books, book)
			responseJSON(w, http.StatusOK, book)
			return
		}
	}
	responseJSON(w, http.StatusBadRequest, nil)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	id, _ := uuid.FromString(r.URL.Query().Get("id"))
	for index, item := range books {
		if item.ID == id {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	responseJSON(w, http.StatusNoContent, books)
}

func responseJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func main() {
	http.HandleFunc("/api/book/list", getBooks)
	http.HandleFunc("/api/book/show", getBook)
	http.HandleFunc("/api/book/create", createBook)
	http.HandleFunc("/api/book/update", updateBook)
	http.HandleFunc("/api/book/delete", deleteBook)

	log.Fatal(http.ListenAndServe(":5555", nil))
}
