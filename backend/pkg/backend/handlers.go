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

	router.GET("/reset", func(context *gin.Context) {
		books := []Book{
			{
				Model:     Model{ID: 1},
				Title:     "Peter Pan",
				Author:    "Barrie, J. M. (James Matthew)",
				Language:  "en",
				Text:      "Peter PanFairies -- Fiction Fantasy literature Never-Never Land (Imaginary place) -- Fiction Peter Pan (Fictitious character) -- Fiction Pirates -- Fiction",
				ImageURL:  "https://www.gutenberg.org/cache/epub/16/pg16.cover.medium.jpg",
				CRank:     0,
				Occurance: 0,
			},
			{
				Model:     Model{ID: 2},
				Title:     "The Book of Mormon",
				Author:    "Church of Jesus Christ of Latter-day Saints",
				Language:  "en",
				Text:      "The Book of Mormon: An Account Written by the Hand of Mormon, Upon Plates Taken from the Plates of NephiChurch of Jesus Christ of Latter-day Saints -- Sacred books Mormon Church -- Sacred books",
				ImageURL:  "https://www.gutenberg.org/cache/epub/17/pg17.cover.medium.jpg",
				CRank:     0,
				Occurance: 0,
			},
			{
				Model:     Model{ID: 3},
				Title:     "The Federalist Papers",
				Author:    "Hamilton, Alexander",
				Language:  "en",
				Text:      "The Federalist PapersConstitutional history -- United States -- Sources Constitutional law -- United States",
				ImageURL:  "https://www.gutenberg.org/cache/epub/18/pg18.cover.medium.jpg",
				CRank:     0,
				Occurance: 0,
			},
		}

		indexedBooks := []IndexedBook{
			{
				Model:               Model{ID: 1},
				Title:               "Peter Pan",
				WorldOccurancesJSON: "{\"test\": 1}",
			},
			{
				Model:               Model{ID: 2},
				Title:               "The Book of Mormon",
				WorldOccurancesJSON: "{\"test\": 2}",
			},
			{
				Model:               Model{ID: 3},
				Title:               "The Federalist Papers",
				WorldOccurancesJSON: "{\"papers\": 1}",
			},
		}

		jaccardNeighbors := []JaccardNeighbors{
			{
				Model:         Model{ID: 1},
				NeighborsJSON: "[2]",
			},
			{
				Model:         Model{ID: 2},
				NeighborsJSON: "[1]",
			},
			{
				Model:         Model{ID: 3},
				NeighborsJSON: "[]",
			},
		}

		db.Exec("DELETE FROM books")
		db.Exec("DELETE FROM indexed_books")
		db.Exec("DELETE FROM jaccard_neighbors")

		h.db.Create(books)
		h.db.Create(indexedBooks)
		h.db.Create(jaccardNeighbors)

		context.JSON(http.StatusOK, "Success")
	})

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
		// Get the World occurance json
		var worldOccurancesMap map[string]uint
		err := json.Unmarshal([]byte(indexedBook.WorldOccurancesJSON), &worldOccurancesMap)

		if err != nil {
			context.JSON(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}

		if count, ok := worldOccurancesMap[query]; ok {
			// Update the occurance count of the books, bad performance!
			var book Book
			h.db.First(&book, indexedBook.ID)
			book.Occurance = count
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

	// Get the occurance names
	var occuranceBooks []Book
	h.db.Where("id in ?", onlyNeigbors.Slice()).Find(&occuranceBooks)

	// Set all the occurances to zero
	for _, element := range occuranceBooks {
		element.Occurance = 0
	}

	// Update the occurances according to zero
	h.db.Updates(occuranceBooks)

	// Now get the neigboring books
	context.JSON(http.StatusOK, map[string]any{
		"books":     returnBooks,
		"neighbors": occuranceBooks,
	})
}

// regex search
func (h *handler) RegexSearchHandler(context *gin.Context) {
	// Get the parameter query
	regex := context.Param("regex")
	fmt.Printf("ZOZZOZOZOOZOZOZOZOZOOZORT: %q", regex)

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
		// Get the World occurance json
		var worldOccurancesMap map[string]uint
		err := json.Unmarshal([]byte(indexedBook.WorldOccurancesJSON), &worldOccurancesMap)

		if err != nil {
			context.JSON(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}

		// Update the occurance count of the books, bad performance!
		var book Book
		h.db.First(&book, indexedBook.ID)
		count := uint(len(expression.FindAllStringIndex(book.Text, -1)))

		if count != 0 {
			book.Occurance = count
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

	// Get the occurance names
	var occuranceBooks []Book
	h.db.Where("id in ?", onlyNeigbors.Slice()).Find(&occuranceBooks)

	// Set all the occurances to zero
	for _, element := range occuranceBooks {
		element.Occurance = 0
	}

	// Update the occurances according to zero
	h.db.Updates(occuranceBooks)

	// Now get the neigboring books
	context.JSON(http.StatusOK, map[string]any{
		"books":     returnBooks,
		"neighbors": occuranceBooks,
	})
}
