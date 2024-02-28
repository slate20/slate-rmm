package main

import (
	"fmt"
	"log"
	"net/http"
	"slate-rmm/database"
	"slate-rmm/handlers"

	"github.com/gorilla/mux"
)

func main() {
	//Initialize the database connection
	dsn := "host=localhost user=postgres password=slatermm dbname=agents_db sslmode=disable"
	database.InitDB(dsn)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	router.HandleFunc("/api/agents/register", handlers.AgentRegistration).Methods("POST")

	// Start the server
	fmt.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))

	log.Println("Server is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
