package main

import (
	"fmt"
	"net/http"

	"github.com/RootLeo00/book-web-app-daar/pkg/backend"
)

func main() {

	// Set up the directories
	http.HandleFunc("/", backend.Index) // Set the handler for the root path

	fmt.Println("Starting server at port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
