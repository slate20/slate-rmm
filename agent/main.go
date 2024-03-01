package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slate-rmm-agent/collectors"
	"slate-rmm-agent/server"
	"time"
)

// Config represents the configuration for the agent
type Config struct {
	ServerURL string `json:"server_url"`
	HostID    int32  `json:"host_id"`
}

func main() {
	var config Config
	var err error

	//Get the directory of the executable
	exe, err := os.Executable()
	if err != nil {
		log.Fatalf("could not get the directory of the executable: %v", err)
	}
	dir := filepath.Dir(exe)

	// Define the path to the config file
	configPath := filepath.Join(dir, "config.json")

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If the config file does not exist, register the agent

		// Prompt the user for the server URL
		fmt.Print("Enter the server URL: ")
		_, err := fmt.Scan(&config.ServerURL)
		if err != nil {
			log.Fatalf("could not read server URL: %v", err)
		}

		// Collect data
		data, err := collectors.CollectData()
		if err != nil {
			log.Fatalf("could not collect data: %v", err)
		}

		// Register the agent
		config.HostID, err = server.Register(data, config.ServerURL)
		if err != nil {
			log.Fatalf("could not register with the server: %v", err)
		}

		// Save the config in the config file
		bytes, err := json.Marshal(config)
		if err != nil {
			log.Fatalf("could not marshal config: %v", err)
		}

		err = os.WriteFile(configPath, bytes, 0644)
		if err != nil {
			log.Fatalf("could not write config file: %v", err)
		}
	} else {
		// If the config file exists, read the config from the file
		bytes, err := os.ReadFile(configPath)
		if err != nil {
			log.Fatalf("could not read config file: %v", err)
		}

		err = json.Unmarshal(bytes, &config)
		if err != nil {
			log.Fatalf("could not unmarshal config: %v", err)
		}
	}

	// Send a heartbeat every minute
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := server.Heartbeat(config.HostID, config.ServerURL); err != nil {
			log.Printf("could not send heartbeat: %v", err)
		}
	}

}
