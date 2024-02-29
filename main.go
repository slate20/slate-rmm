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
	router.HandleFunc("/api/agents", handlers.GetAllAgents).Methods("GET")
	router.HandleFunc("/api/agents/{id}", handlers.GetAgent).Methods("GET")
	router.HandleFunc("/api/agents/{id}", handlers.UpdateAgent).Methods("PUT")
	router.HandleFunc("/api/agents/{id}", handlers.DeleteAgent).Methods("DELETE")

	// Start the server
	fmt.Println("Starting server on the port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))

}
