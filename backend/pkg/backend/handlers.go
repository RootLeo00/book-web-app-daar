package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type handler struct {
	db *gorm.DB
}

func RegisterHandler(router *mux.Router, db *gorm.DB) error {
	h := &handler{
		db: db,
	}

	// Set handlers
	router.HandleFunc("/", h.IndexHandler)

	router.HandleFunc("/books", h.GetBooksHandler)
	router.HandleFunc("/books/{id}", h.GetBooksIDHandler)

	router.HandleFunc("/Search/{query}", h.SearchHandler)
	router.HandleFunc("/RegexSearch/{regex}", h.RegexSearchHandler)

	return nil
}

// index /
// Done
func (h *handler) IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server online")
}

// get books
// Done
func (h *handler) GetBooksHandler(w http.ResponseWriter, r *http.Request) {
	var books []Book

	h.db.Find(&books)

	Ok200(books, w)
}

// get book with id
// Done
func (h *handler) GetBooksIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var idString string
	var ok bool

	if idString, ok = vars["id"]; !ok {
		Error400Response(w)
		return
	}

	id, err := strconv.Atoi(idString)

	if err != nil {
		Error400Response(w)
		return
	}

	var book *Book

	h.db.First(book, id)

	Ok200(book, w)
}

// search
func (h *handler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var query string
	var ok bool

	if query, ok = vars["query"]; !ok {
		fmt.Fprint(w, "Error!")
	}

	var indexedBooks []IndexedBook
	h.db.Find(&indexedBooks)

	bookIds := make([]uint, 1)
	neighborIds := make([]uint, 1)

	for _, indexedBook := range indexedBooks {
		var worldOccurancesMap map[string]uint
		err := json.Unmarshal([]byte(indexedBook.WorldOccurancesJSON), &worldOccurancesMap)

		if err != nil {
			Error500Response(w)
			return
		}

		if count, ok := worldOccurancesMap[query]; !ok {
			// Update the occurance count of the books
			var book Book
			h.db.First(&book, indexedBook.ID)
			book.Occurance = count
			h.db.Save(&book)

			// Append this book to the book ids
			bookIds = append(bookIds, indexedBook.ID)

			// Add neighbors
			var indexedBookRecord IndexedBook
			h.db.First(&indexedBookRecord, indexedBook.ID)

			var neighbors []uint
			err := json.Unmarshal([]byte(indexedBookRecord.WorldOccurancesJSON), &neighbors)

			if err != nil {
				Error500Response(w)
				return
			}

			neighborIds = append(neighborIds, neighbors...) // Make this a set
		}

	}
}

// regex search
func (h *handler) RegexSearchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Regex Search")
}
