package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/exec"
	"slate-rmm/database"
	"slate-rmm/models"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var agentTokens = make(map[string]string)

// AgentRegistration handles the registration of a new agent
func AgentRegistration(w http.ResponseWriter, r *http.Request) {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		http.Error(w, "could not load .env file", http.StatusInternalServerError)
		return
	}

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		log.Fatal("API_URL is not set in .env file")
	}
	siteName := os.Getenv("SITE_NAME")
	if siteName == "" {
		log.Fatal("SITE_NAME is not set in .env file")
	}
	apiUser := os.Getenv("API_USER")
	if apiUser == "" {
		log.Fatal("API_USER is not set in .env file")
	}
	apiPass := os.Getenv("AUTOMATION_SECRET")
	if apiPass == "" {
		log.Fatal("AUTOMATION_SECRET is not set in .env file")
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
	req, err := http.NewRequest("POST", apiURL+"/domain-types/host_config/collections/all", strings.NewReader(string(payloadStr)))
	if err != nil {
		http.Error(w, "error creating request", http.StatusInternalServerError)
		return
	}

	// Set the content type to application/json
	req.Header.Set("Content-Type", "application/json")
	// Set the authorization header
	req.Header.Set("Authorization", "Bearer "+apiUser+" "+apiPass)
	// req.Header.Set
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

	//Generate a one-time token for the agent
	token := uuid.New().String()
	newAgent.Token = token

	// Store the token and the agent ID in the agentTokens map
	agentTokens[newAgent.Hostname] = token

	// Respond with the registered agent
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newAgent)

	// Sleep for 5 seconds to allow host creation to complete
	time.Sleep(5 * time.Second)

	// Run the CheckMK service discovery script
	cmd := exec.Command("./handlers/cmk_svcd.sh", newAgent.Hostname)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
	}
}

// Verify agent token and return $AUTOMATION_SECRET
func VerifyAgentToken(w http.ResponseWriter, r *http.Request) {
	// Decode the incoming JSON to get the token and agent ID
	var data map[string]string
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, ok := data["token"]
	if !ok {
		http.Error(w, "Token not provided", http.StatusBadRequest)
		return
	}

	agentID, ok := data["agent_id"]
	if !ok {
		http.Error(w, "Agent ID not provided", http.StatusBadRequest)
		return
	}

	//Verify the token
	storedToken, ok := agentTokens[agentID]
	if !ok || token != storedToken {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Delete the token from the agentTokens map
	delete(agentTokens, agentID)

	// If the token is valid, return the AUTOMATION_SECRET
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"AUTOMATION_SECRET": os.Getenv("AUTOMATION_SECRET"),
	})
}

// CheckMKServiceDiscovery runs the CheckMK service discovery script
func CMKSvcDiscovery(w http.ResponseWriter, r *http.Request) {
	log.Println("Received service discovery request")
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		http.Error(w, "could not load .env file", http.StatusInternalServerError)
		return
	}

	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		log.Fatal("API_URL is not set in .env file")
	}
	apiUser := os.Getenv("API_USER")
	if apiUser == "" {
		log.Fatal("API_USER is not set in .env file")
	}
	apiPass := os.Getenv("AUTOMATION_SECRET")
	if apiPass == "" {
		log.Fatal("AUTOMATION_SECRET is not set in .env file")
	}

	// Decode the incoming JSON to get the hostname
	var data map[string]string
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	hostname := data["host_name"]

	// Run the CheckMK service discovery script
	log.Println("Running service discovery script")
	cmd := exec.Command("./handlers/cmk_svcd.sh", hostname)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Printf("cmd.Run() failed with %s\n", err)
	}
	log.Println("Service discovery script complete")
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
