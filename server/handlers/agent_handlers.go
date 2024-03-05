package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"slate-rmm/database"
	"slate-rmm/models"

	"github.com/gorilla/mux"
)

// AgentRegistration handles the registration of a new agent
func AgentRegistration(w http.ResponseWriter, r *http.Request) {
	var newAgent models.Agent
	// Decode the incoming JSON to the newAgent struct
	if err := json.NewDecoder(r.Body).Decode(&newAgent); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := database.RegisterNewAgent(&newAgent); err != nil {
		http.Error(w, "error registering agent", http.StatusInternalServerError)
		return
	}

	// Respond with the registered agent
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAgent)
}

// GetAllAgents returns all the agents in the database
func GetAllAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := database.GetAllAgents()
	if err != nil {
		log.Printf("error getting agents: %v", err)
		http.Error(w, "error getting agents", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(agents)
}

// GetAgent returns a single agent from the database
func GetAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	agent, err := database.GetAgent(id)
	if err != nil {
		log.Printf("error getting agent: %v", err)
		http.Error(w, "error getting agent", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(agent)
}

// UpdateAgent updates an agent in the database
func UpdateAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var updatedAgent models.Agent
	err := json.NewDecoder(r.Body).Decode(&updatedAgent)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = database.UpdateAgent(id, &updatedAgent)
	if err != nil {
		http.Error(w, "error updating agent", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteAgent deletes an agent from the database
func DeleteAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := database.DeleteAgent(id)
	if err != nil {
		http.Error(w, "error deleting agent", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// AgentHeartbeat updates the last seen time of an agent
func AgentHeartbeat(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := database.AgentHeartbeat(id)
	if err != nil {
		http.Error(w, "error updating agent", http.StatusInternalServerError)
		return
	}
}
