package backend

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	db *gorm.DB
}

func RegisterHandler(router *gin.Engine, db *gorm.DB) error {
	h := &handler{
		db: db,
	}

	// Set handlers
	router.GET("/", h.IndexHandler)

	router.GET("/books", h.GetBooksHandler)
	router.GET("/books/:id", h.GetBooksIDHandler)

	router.GET("/Search/:query", h.SearchHandler)
	router.GET("/RegexSearch/:regex", h.RegexSearchHandler)

	return nil
}

// The index page
func (h *handler) IndexHandler(context *gin.Context) {
	context.JSON(http.StatusOK, "Success")
}

// Gets all of the books in the database
func (h *handler) GetBooksHandler(context *gin.Context) {
	// Get all of the books
	var books []Book
	h.db.Find(&books)

	// Response with the books list
	context.JSON(http.StatusOK, books)
}

// Gets the book object by id
func (h *handler) GetBooksIDHandler(context *gin.Context) {
	// Get the parameter id
	id := context.Param("id")

	// Get the book with the id from database
	var book *Book
	h.db.First(book, id)

	// Response with the book object
	context.JSON(http.StatusOK, book)
}

// search
func (h *handler) SearchHandler(context *gin.Context) {
	// Get the parameter query
	query := context.Param("query")

	var indexedBooks []IndexedBook
	h.db.Find(&indexedBooks)

	bookIds := make([]uint, 1)
	neighborIds := make([]uint, 1)

	for _, indexedBook := range indexedBooks {
		var worldOccurancesMap map[string]uint
		err := json.Unmarshal([]byte(indexedBook.WorldOccurancesJSON), &worldOccurancesMap)

		if err != nil {
			context.JSON(http.StatusInternalServerError, []Book{})
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
				context.JSON(http.StatusInternalServerError, []Book{})
				return
			}

			neighborIds = append(neighborIds, neighbors...) // Make this a set
		}
	}
	context.JSON(http.StatusOK, []Book{})
}

// regex search
func (h *handler) RegexSearchHandler(context *gin.Context) {
	context.JSON(http.StatusOK, []Book{})
}
