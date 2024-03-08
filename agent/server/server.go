package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"slate-rmm-agent/collectors"
)

// Register sends a POST request to the server to register the agent
func Register(data collectors.AgentData, ServerURL string) (int32, error) {
	url := ServerURL + "/api/agents/register"
	// Convert data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	// Send a POST request to the AgentRegister endpoint
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Decode the response
	var result struct {
		HostID int32 `json:"host_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	// Check the response status code
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return result.HostID, nil
}

// Heartbeat sends a POST request to the server to indicate that the agent is still alive
func Heartbeat(hostID int32, ServerURL string) error {
	url := ServerURL + "/api/agents/" + fmt.Sprint(hostID) + "/heartbeat"
	// Send a POST request to the AgentHeartbeat endpoint
	resp, err := http.Post(url, "", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	return nil
}
