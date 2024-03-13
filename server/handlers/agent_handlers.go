package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slate-rmm/database"
	"slate-rmm/models"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// AgentRegistration handles the registration of a new agent
func AgentRegistration(w http.ResponseWriter, r *http.Request) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		http.Error(w, "could not load .env file", http.StatusInternalServerError)
		return
	}

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

	// Prepare payload for CheckMK host creation
	payload := map[string]interface{}{
		"folder":    "/",
		"host_name": newAgent.Hostname,
		"attributes": map[string]string{
			"ipaddress": newAgent.IPAddress,
		},
	}
	payloadStr, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "could not marshal payload", http.StatusInternalServerError)
		return
	}

	log.Printf("Payload: %s\n", payloadStr)

	// Create a new request
	req, err := http.NewRequest("POST", "http://localhost:5000/main/check_mk/api/1.0/domain-types/host_config/collections/all", strings.NewReader(string(payloadStr)))
	if err != nil {
		http.Error(w, "error creating request", http.StatusInternalServerError)
		return
	}

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")
	// Set the authorization header
	req.Header.Set("Authorization", "Bearer "+os.Getenv("API_USER")+" "+os.Getenv("AUTOMATION_SECRET"))
	// req.Header.Set("Authorization", "Bearer cmkadmin slatermm")
	req.Header.Set("Accept", "application/json")

	log.Printf("Request: %v\n", req)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "error sending request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Printf("Response: %v\n", resp)

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
