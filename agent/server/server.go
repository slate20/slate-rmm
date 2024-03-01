package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slate-rmm-agent/collectors"
)

const ServerURL = "http://192.168.1.10:8080" // Replace with actual server URL

// Register sends a POST request to the server to register the agent
func Register(data collectors.AgentData) error {
	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Send a POST request to the AgentRegister endpoint
	resp, err := http.Post(ServerURL+"/api/agents/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
