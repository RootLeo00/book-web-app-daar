package backend

import (
	"encoding/json"
	"fmt"
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

	router.GET("/reset", func(context *gin.Context) {
		books := []Book{
			{
				Model:     gorm.Model{ID: 1},
				Title:     "Peter Pan",
				Author:    "Barrie, J. M. (James Matthew)",
				Language:  "en",
				Text:      "Peter PanFairies -- Fiction Fantasy literature Never-Never Land (Imaginary place) -- Fiction Peter Pan (Fictitious character) -- Fiction Pirates -- Fiction",
				ImageURL:  "https://www.gutenberg.org/cache/epub/16/pg16.cover.medium.jpg",
				CRank:     0,
				Occurance: 0,
			},
			{
				Model:     gorm.Model{ID: 2},
				Title:     "The Book of Mormon: An Account Written by the Hand of Mormon, Upon Plates Taken from the Plates of Nephi",
				Author:    "Church of Jesus Christ of Latter-day Saints",
				Language:  "en",
				Text:      "The Book of Mormon: An Account Written by the Hand of Mormon, Upon Plates Taken from the Plates of NephiChurch of Jesus Christ of Latter-day Saints -- Sacred books Mormon Church -- Sacred books",
				ImageURL:  "https://www.gutenberg.org/cache/epub/17/pg17.cover.medium.jpg",
				CRank:     0,
				Occurance: 0,
			},
			{
				Model:     gorm.Model{ID: 3},
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
				Model:               gorm.Model{ID: 1},
				Title:               "Peter Pan",
				WorldOccurancesJSON: "{\"peter\": 2, \"panfairies\": 1, \"fiction\": 4, \"fantasy\": 1, \"literature\": 1, \"never\": 2, \"land\": 1, \"imaginary\": 1, \"place\": 1, \"pan\": 1, \"fictitious\": 1, \"character\": 1, \"pirates\": 1}",
			},
			{
				Model:               gorm.Model{ID: 2},
				Title:               "The Book of Mormon: An Account Written by the Hand of Mormon, Upon Plates Taken from the Plates of Nephi",
				WorldOccurancesJSON: "{\"book\": 1, \"mormon\": 3, \"account\": 1, \"written\": 1, \"hand\": 1, \"upon\": 1, \"plates\": 2, \"taken\": 1, \"nephichurch\": 1, \"jesus\": 1, \"christ\": 1, \"latter\": 1, \"day\": 1, \"saints\": 1, \"sacred\": 2, \"books\": 2, \"church\": 1}",
			},
			{
				Model:               gorm.Model{ID: 3},
				Title:               "The Federalist Papers",
				WorldOccurancesJSON: "{\"federalist\": 1, \"papersconstitutional\": 1, \"history\": 1, \"united\": 2, \"states\": 2, \"sources\": 1, \"constitutional\": 1, \"law\": 1}",
			},
		}

		jaccardNeighbors := []JaccardNeighbors{
			{
				Model:         gorm.Model{ID: 1},
				NeighborsJSON: "[2]",
			},
			{
				Model:         gorm.Model{ID: 2},
				NeighborsJSON: "[1, 3]",
			},
			{
				Model:         gorm.Model{ID: 3},
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

	// For each book calculate
	for _, indexedBook := range indexedBooks {
		// Get the World occurance json
		var worldOccurancesMap map[string]uint
		err := json.Unmarshal([]byte(indexedBook.WorldOccurancesJSON), &worldOccurancesMap)

		if err != nil {
			context.JSON(http.StatusInternalServerError, fmt.Sprint(err))
			return
		}

		if count, ok := worldOccurancesMap[query]; !ok {
			// Update the occurance count of the books, bad performance!
			var book Book
			h.db.First(&book, indexedBook.ID)
			book.Occurance = count
			h.db.Save(&book)

			// Append this book to the book ids
			bookIds = append(bookIds, indexedBook.ID)

			// Add neighbors
			var jaccardNeighbors JaccardNeighbors
			h.db.First(&jaccardNeighbors, indexedBook.ID)

			var neighbors []uint
			err := json.Unmarshal([]byte(jaccardNeighbors.NeighborsJSON), &neighbors)

			if err != nil {
				context.JSON(http.StatusInternalServerError, fmt.Sprint(err))
				return
			}

			neighborIds = append(neighborIds, neighbors...) // Make this a set
		}
	}

	// Now get the neigboring books
	context.JSON(http.StatusOK, map[string]any{
		"books":     bookIds,
		"neighbors": neighborIds,
	})
}

// regex search
func (h *handler) RegexSearchHandler(context *gin.Context) {
	context.JSON(http.StatusOK, []Book{})
}
