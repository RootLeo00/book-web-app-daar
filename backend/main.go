package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/RootLeo00/book-web-app-daar/pkg/backend"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const MAX_TRIES = 3

func main() {
	var dsn string
	var db *gorm.DB
	var ok bool
	var err error

	if dsn, ok = os.LookupEnv("DATABASE_URL"); !ok {
		panic("environment variable DATABASE_URL is not provided")
	}

	for tries := 1; tries <= MAX_TRIES; tries++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			if tries < MAX_TRIES {
				// Sleep for 1 second to try again
				fmt.Println("cannot connect to the database")
				time.Sleep(time.Second)
			} else {
				panic("cannot connect to the database, all tries failed")
			}
		} else {
			break
		}
	}

	fmt.Println("Connected to the database, migrating...")

	// Automatically migrate the following models
	db.AutoMigrate(
		&backend.Book{},
		&backend.IndexedBook{},
		&backend.JaccardNeighbors{})

	router := gin.Default()

	// Register the handlers to the mux
	err = backend.RegisterHandler(router, db)

	if err != nil {
		panic("Cannot create the handler")
	}

	server := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Starting server at port 8080...")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
