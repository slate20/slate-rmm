package main

import (
	"fmt"
	"log"
	"net/http"
	"slate-rmm/database"
)

// main is the entry point for the application
func main() {
	//Initialize the database connection
	dsn := "host=localhost user=postgres password=slatermm dbname=RMM_db sslmode=disable"
	database.InitDB(dsn)

	// Create a new router
	router := NewGateway()

	// Add the CORS middleware
	corsRouter := CORSMiddleware(router)

	// Start the server
	fmt.Println("Starting server on the port 8081...")
	log.Fatal(http.ListenAndServe(":8081", corsRouter))
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

		// If it's just an OPTIONS request, we don't need to go any further
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
