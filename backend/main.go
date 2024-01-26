package main

import (
	"net/http"
	"os"
	"time"

	"github.com/RootLeo00/book-web-app-daar/pkg/backend"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const MAX_TRIES = 3

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		// c.Writer.Header().Set("Origin", "https://localhost:8080")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	var dsn string
	var db *gorm.DB
	var ok bool
	var err error

	if dsn, ok = os.LookupEnv("DATABASE_URL"); !ok {
		logger.Default.Info(nil, "environment variable DATABASE_URL is not provided, looking for the local sql file...")

		db, err = gorm.Open(sqlite.Open("./db.sqlite3"), &gorm.Config{})

		if err != nil {
			panic("cannot open to the sqlite3 database")
		}

	} else {
		// Try to connect to the postgres, if cannot connect sleep for 1 second and try again
		for tries := 1; tries <= MAX_TRIES; tries++ {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

			if err != nil {
				if tries < MAX_TRIES {
					// Sleep for 1 second to try again
					logger.Default.Info(nil, "cannot connect to the database")
					time.Sleep(time.Second)
				} else {
					panic("cannot connect to the database, all tries failed")
				}
			} else {
				break
			}
		}

	}

	logger.Default.Info(nil, "Connected to the database, migrating...")

	// Automatically migrate the following models
	db.AutoMigrate(
		&backend.Book{},
		&backend.IndexedBook{},
		&backend.JaccardNeighbors{})

	router := gin.Default()

	// Add middewares
	router.Use(CORSMiddleware())
	// router.Use(cors.New(cors.Config{
	// 	AllowAllOrigins: true,
	// 	AllowHeaders:    []string{"Origin"},
	// }))

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
		MaxHeaderBytes: 1 << 32,
	}

	logger.Default.Info(nil, "Starting server at port 8080...")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
