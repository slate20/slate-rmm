package main

import (
	"slate-rmm/handlers"

	"github.com/gorilla/mux"
)

// NewGateway creates a new router and defines the routes for the microservices
func NewGateway() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Define routes for each microservice
	agentRoutes(router.PathPrefix("/api/agents").Subrouter())

	return router

}

// agentRoutes defines the routes for the agent database microservice
func agentRoutes(router *mux.Router) {
	router.HandleFunc("/register", handlers.AgentRegistration).Methods("POST")
	router.HandleFunc("", handlers.GetAllAgents).Methods("GET")
	router.HandleFunc("/{id}", handlers.GetAgent).Methods("GET")
	router.HandleFunc("/{id}", handlers.UpdateAgent).Methods("PUT")
	router.HandleFunc("/{id}", handlers.DeleteAgent).Methods("DELETE")
	router.HandleFunc("/{id}/heartbeat", handlers.AgentHeartbeat).Methods("POST")
}
