package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/khengsaurus/go-tutorials/bookstore-api/pkg/models"
	"github.com/khengsaurus/go-tutorials/bookstore-api/pkg/utils"
)

// HELP in utils this fn can't be found?
func json200(payload any, w http.ResponseWriter) {
	res, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	newBooks := models.GetAllBooks()
	json200(newBooks, w)
}

func GetBookById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_id := vars["bookId"]
	if _id == "" {
		http.Error(w, "400 missing parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(_id, 0, 0)
	if err != nil {
		fmt.Println("Error while parsing")
		http.Error(w, "500 internal server error", http.StatusInternalServerError)
		return
	}
	book, _ := models.GetBookById(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if book.ID == 0 {
		w.Write([]byte("Book not found"))
	} else {
		json.NewEncoder(w).Encode((book))
	}
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	createBook := &models.Book{}
	utils.ParseBody(r, createBook)
	newBook := createBook.CreateBook()
	json200(newBook, w)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	var updateBook = &models.Book{}
	utils.ParseBody(r, updateBook)
	vars := mux.Vars(r)
	_id := vars["bookId"]
	id, _ := strconv.ParseInt(_id, 0, 0)
	book, db := models.GetBookById(id)
	if updateBook.Name != "" {
		book.Name = updateBook.Name
	}
	if updateBook.Author != "" {
		book.Author = updateBook.Author
	}
	if updateBook.Publication != "" {
		book.Publication = updateBook.Publication
	}
	db.Save(&book)
	json200(book, w)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	// NB this only adds the 'deleted_at' field in MySQL DB. Due to config?
	vars := mux.Vars(r)
	_id := vars["bookId"]
	id, _ := strconv.ParseInt(_id, 0, 0)
	book := models.DeleteBook(id)
	json200(book, w)
}
