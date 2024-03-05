package main

import (
	"net/http"
	"slate-rmm/handlers"

	"github.com/gorilla/mux"
)

// NewGateway creates a new router and defines the routes for the microservices
func NewGateway() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	// Define routes for each microservice
	agentRoutes(router.PathPrefix("/api/agents").Subrouter())

	// Serve the agent executable
	router.HandleFunc("/download/agent", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Disposition", "attachment; filename=slate-rmm-agent")
		http.ServeFile(w, r, "../agent/slate-rmm-agent")
	})

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
