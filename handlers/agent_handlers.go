package handlers

import (
	"encoding/json"
	"net/http"
	"slate-rmm/database"
	"slate-rmm/models"
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
