package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
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
		log.Printf("could not get the directory of the executable: %v", err)
	}
	dir := filepath.Dir(exe)

	// Define the path to the config file
	configPath := filepath.Join(dir, "config.json")

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// If the config file does not exist, register the agent

		// Prompt the user for the server URL
		fmt.Print("Enter the server IP or Hostname (If DNS is configured): ")
		var serverURL string
		_, err := fmt.Scan(&serverURL)
		if err != nil {
			log.Printf("could not read server IP/Hostname: %v", err)
		}

		// Append "http://" and ":8080" to the server URL
		config.ServerURL = "http://" + serverURL + ":8080"

		// Collect data
		data, err := collectors.CollectData()
		if err != nil {
			log.Printf("could not collect data: %v", err)
		}

		// Save the config in the config file
		configBytes, err := json.Marshal(config)
		if err != nil {
			log.Printf("could not marshal config: %v", err)
		}

		err = os.WriteFile(configPath, configBytes, 0644)
		if err != nil {
			log.Printf("could not write config file: %v", err)
		}

		// Download the CheckMK agent
		err = downloadFile("http://localhost:5000/main/check_mk/agents/windows/check_mk_agent.msi", "check_mk_agent.msi")
		if err != nil {
			log.Printf("could not download CheckMK agent: %v", err)
		}

		// Install the CheckMK agent
		err = exec.Command("msiexec", "/i", "check_mk_agent.msi", "/qn").Run()
		if err != nil {
			log.Printf("could not install CheckMK agent: %v", err)
		}

		// Download CheckMK Inventory plugin
		err = downloadFile("http://localhost:5000/main/check_mk/agents/windows/plugins/mk_inventory.vbs", "mk_inventory.vbs")
		if err != nil {
			log.Printf("could not download CheckMK Inventory plugin: %v", err)
		}

		// Move the CheckMK Inventory plugin to the plugins directory
		err = os.Rename("mk_inventory.vbs", "C:\\Program Files (x86)\\checkmk\\service\\plugins\\mk_inventory.vbs")
		if err != nil {
			log.Printf("could not move CheckMK Inventory plugin: %v", err)
		}

		// Register the agent
		config.HostID, err = server.Register(data, config.ServerURL)
		if err != nil {
			log.Printf("could not register with the server: %v", err)
		}

		// Register the CheckMK agent with the CheckMK server
		cmd := exec.Command("C:\\Program Files (x86)\\checkmk\\service\\cmk-agent-ctl.exe", "register", "--hostname", data.Hostname, "--server", serverURL+":8000", "--site", "main", "--user", "cmkadmin", "--password", "slatermmdev")

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			log.Printf("could not register CheckMK agent: %v", err)
			log.Printf("stdout: %s", stdout.String())
			log.Printf("stderr: %s", stderr.String())
		}

	} else {
		// If the config file exists, read the config from the file
		configBytes, err := os.ReadFile(configPath)
		if err != nil {
			log.Printf("could not read config file: %v", err)
		}

		err = json.Unmarshal(configBytes, &config)
		if err != nil {
			log.Printf("could not unmarshal config: %v", err)
		}
	}

	// Pause the program
	fmt.Println("Press enter to close...")
	fmt.Scanln()

	// Send a heartbeat every minute
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := server.Heartbeat(config.HostID, config.ServerURL); err != nil {
			log.Printf("could not send heartbeat: %v", err)
		}
	}

}

// downloadFile downloads a file from the given URL and saves it to the given path
func downloadFile(url string, path string) error {
	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the data to the file
	_, err = io.Copy(out, resp.Body)
	return err
}
