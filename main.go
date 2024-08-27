package main

import (
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/Nikhils-179/connection-pool/handler"
)

func main() {
	router := http.NewServeMux()

	// Set up the routes
	router.HandleFunc("/list-following", handler.Handler)
	router.HandleFunc("/list-following-fix", handler.HandlerFix)

	// Start the server
	fmt.Println("Starting server at :4000")
	if err := http.ListenAndServe(":4000", router); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}