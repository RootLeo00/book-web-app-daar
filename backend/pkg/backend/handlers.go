package backend

import (
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

// index /
// Done
func (h *handler) IndexHandler(context *gin.Context) {
	context.JSON(http.StatusOK, "Server Online")
}

// get books
// Done
func (h *handler) GetBooksHandler(context *gin.Context) {
	var books []Book

	h.db.Find(&books)

	context.JSON(http.StatusOK, books)
}

// get book with id
// Done
func (h *handler) GetBooksIDHandler(context *gin.Context) {
	// vars := mux.Vars(r)
	// var id string
	// var ok bool

	// if id, ok = vars["id"]; !ok {
	// 	context.JSON(http.StatusInternalServerError, "Cannot get the id")
	// 	return
	// }

	// var book *Book

	// h.db.First(book, id)

	context.JSON(http.StatusOK, []Book{})
}

// search
func (h *handler) SearchHandler(context *gin.Context) {
	// if query, ok = vars["query"]; !ok {
	// 	context.JSON(http.StatusInternalServerError, "Cannot get the query")
	// }

	// var indexedBooks []IndexedBook
	// h.db.Find(&indexedBooks)

	// bookIds := make([]uint, 1)
	// neighborIds := make([]uint, 1)

	// for _, indexedBook := range indexedBooks {
	// 	var worldOccurancesMap map[string]uint
	// 	err := json.Unmarshal([]byte(indexedBook.WorldOccurancesJSON), &worldOccurancesMap)

	// 	if err != nil {
	// 		Error500Response(w)
	// 		return
	// 	}

	// 	if count, ok := worldOccurancesMap[query]; !ok {
	// 		// Update the occurance count of the books
	// 		var book Book
	// 		h.db.First(&book, indexedBook.ID)
	// 		book.Occurance = count
	// 		h.db.Save(&book)

	// 		// Append this book to the book ids
	// 		bookIds = append(bookIds, indexedBook.ID)

	// 		// Add neighbors
	// 		var indexedBookRecord IndexedBook
	// 		h.db.First(&indexedBookRecord, indexedBook.ID)

	// 		var neighbors []uint
	// 		err := json.Unmarshal([]byte(indexedBookRecord.WorldOccurancesJSON), &neighbors)

	// 		if err != nil {
	// 			Error500Response(w)
	// 			return
	// 		}

	// 		neighborIds = append(neighborIds, neighbors...) // Make this a set
	// 	}
	// }
	context.JSON(http.StatusOK, []Book{})
}

// regex search
func (h *handler) RegexSearchHandler(context *gin.Context) {
	context.JSON(http.StatusOK, []Book{})
}
