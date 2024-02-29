package main

import (
	"log"
	"slate-rmm-agent/collectors"
	"slate-rmm-agent/server"
)

func main() {
	// Collect data
	data, err := collectors.CollectData()
	if err != nil {
		log.Fatalf("could not collect data: %v", err)
	}

	// Register the agent
	err = server.Register(data)
	if err != nil {
		log.Fatalf("could not register with the server: %v", err)
	}
}
