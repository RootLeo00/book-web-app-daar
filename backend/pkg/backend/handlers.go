package backend

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-set/v2"
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
	var book Book
	h.db.First(&book, id)

	// Response with the book object
	context.JSON(http.StatusOK, book)
}

// search
func (h *handler) SearchHandler(context *gin.Context) {
	// Get the parameter query
	query := context.Param("query")

	// Get all of the indexed books
	var indexedBooks []IndexedBook
	h.db.Find(&indexedBooks)

	// Define the actual books
	returnBooks := set.New[Book](0)
	bookIds := set.New[uint](len(indexedBooks))
	neighborIds := set.New[uint](0)

	// For each book calculate
	for _, indexedBook := range indexedBooks {
		// Get the World occurrence json
		var wordOccurrencesMap map[string]uint
		err := json.Unmarshal([]byte(indexedBook.WordOccurrenceJSON), &wordOccurrencesMap)

		if err != nil {
			context.JSON(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}

		if count, ok := wordOccurrencesMap[query]; ok {
			// Update the occurrence count of the books, bad performance!
			var book Book
			h.db.First(&book, indexedBook.ID)
			book.Occurrence = count
			h.db.Save(&book)

			// Append this book to the book ids
			bookIds.Insert(indexedBook.ID)
			returnBooks.Insert(book)

			// Add neighbors
			var jaccardNeighbors JaccardNeighbors
			h.db.First(&jaccardNeighbors, indexedBook.ID)

			var neighbors []uint
			err := json.Unmarshal([]byte(jaccardNeighbors.NeighborsJSON), &neighbors)

			if err != nil {
				context.JSON(http.StatusInternalServerError, fmt.Sprint(err))
				return
			}

			neighborIds.InsertSlice(neighbors)
		}
	}

	// Get only the negigbor ones
	onlyNeigbors := neighborIds.Difference(bookIds)

	// Get the occurrence names
	var occurrenceBooks []Book
	h.db.Where("id in ?", onlyNeigbors.Slice()).Find(&occurrenceBooks)

	// Set all the occurrences to zero
	for _, element := range occurrenceBooks {
		element.Occurrence = 0
	}

	// Update the occurrences according to zero
	h.db.Updates(occurrenceBooks)

	// Now get the neigboring books
	context.JSON(http.StatusOK, map[string]any{
		"books":     returnBooks,
		"neighbors": occurrenceBooks,
	})
}

// regex search
func (h *handler) RegexSearchHandler(context *gin.Context) {
	// Get the parameter query
	regex := context.Param("regex")

	expression, err := regexp.Compile(regex)

	if err != nil {
		context.JSON(http.StatusBadRequest, fmt.Sprintf("cannot compile regex %q", regex))
		return
	}

	// Get all of the indexed books
	var indexedBooks []IndexedBook
	h.db.Find(&indexedBooks)

	// Define the actual books
	returnBooks := set.New[Book](0)
	bookIds := set.New[uint](len(indexedBooks))
	neighborIds := set.New[uint](0)

	// For each book calculate
	for _, indexedBook := range indexedBooks {
		// Get the World occurrence json
		var wordOccurrencesMap map[string]uint
		err := json.Unmarshal([]byte(indexedBook.WordOccurrenceJSON), &wordOccurrencesMap)

		if err != nil {
			context.JSON(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}

		// Update the occurrence count of the books, bad performance!
		var book Book
		h.db.First(&book, indexedBook.ID)
		count := uint(len(expression.FindAllStringIndex(book.Text, -1)))

		if count != 0 {
			book.Occurrence = count
			h.db.Save(&book)

			// Append this book to the book ids
			bookIds.Insert(indexedBook.ID)
			returnBooks.Insert(book)

			// Add neighbors
			var jaccardNeighbors JaccardNeighbors
			h.db.First(&jaccardNeighbors, indexedBook.ID)

			var neighbors []uint
			err = json.Unmarshal([]byte(jaccardNeighbors.NeighborsJSON), &neighbors)

			if err != nil {
				context.JSON(http.StatusInternalServerError, fmt.Sprint(err))
				return
			}

			neighborIds.InsertSlice(neighbors)
		}

	}

	// Get only the negigbor ones
	onlyNeigbors := neighborIds.Difference(bookIds)

	// Get the occurrence names
	var occurrenceBooks []Book
	h.db.Where("id in ?", onlyNeigbors.Slice()).Find(&occurrenceBooks)

	// Set all the occurrences to zero
	for _, element := range occurrenceBooks {
		element.Occurrence = 0
	}

	// Update the occurrences according to zero
	h.db.Updates(occurrenceBooks)

	// Now get the neigboring books
	context.JSON(http.StatusOK, map[string]any{
		"books":     returnBooks,
		"neighbors": occurrenceBooks,
	})
}
