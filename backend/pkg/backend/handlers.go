package backend

import (
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

	fmt.Fprintf(w, "Search q %q", query)
}

// regex search
func (h *handler) RegexSearchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Regex Search")
}
