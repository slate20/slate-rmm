package database

import (
	"database/sql"
	"encoding/json"
	"log"
	"slate-rmm/models"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
}

// RegisterNewAgent stores a new agent in the database
func RegisterNewAgent(agent *models.Agent) error {

	// Marshal the hardware specs
	hardwareSpecsJSON, err := json.Marshal(agent.HardwareSpecs)
	if err != nil {
		return err
	}

	// Prepare for SQL Statement
	stmt, err := db.Prepare(`
		INSERT INTO agents (hostname, ip_address, os, os_version, hardware_specs, agent_version, last_seen)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING host_id
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the SQL Statement. For 'last_seen', we use the current time
	err = stmt.QueryRow(
		agent.Hostname,
		agent.IPAddress,
		agent.OS,
		agent.OSVersion,
		hardwareSpecsJSON,
		agent.AgentVersion,
		time.Now(),
	).Scan(&agent.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetAllAgents returns all the agents in the database
func GetAllAgents() ([]models.Agent, error) {
	rows, err := db.Query("SELECT * FROM agents")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and add the agents to the slice
	var agents []models.Agent
	for rows.Next() {
		var agent models.Agent
		var hardwareSpecsRaw sql.RawBytes
		if err := rows.Scan(&agent.ID, &agent.Hostname, &agent.IPAddress, &agent.OS, &agent.OSVersion, &hardwareSpecsRaw, &agent.AgentVersion, &agent.LastSeen); err != nil {
			return nil, err
		}

		// Unmarshal the hardware specs
		if len(hardwareSpecsRaw) > 0 {
			if err := json.Unmarshal([]byte(hardwareSpecsRaw), &agent.HardwareSpecs); err != nil {
				return nil, err
			}
		}

		agents = append(agents, agent)
	}

	return agents, nil
}

// GetAgent returns a single agent from the database
func GetAgent(id string) (*models.Agent, error) {
	row := db.QueryRow("SELECT * FROM agents WHERE host_id = $1", id)

	// Scan the row into an Agent struct
	var agent models.Agent
	var hardwareSpecsRaw []byte
	if err := row.Scan(&agent.ID, &agent.Hostname, &agent.IPAddress, &agent.OS, &agent.OSVersion, &hardwareSpecsRaw, &agent.AgentVersion, &agent.LastSeen); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Unmarshal the hardware specs
	if len(hardwareSpecsRaw) > 0 {
		if err := json.Unmarshal([]byte(hardwareSpecsRaw), &agent.HardwareSpecs); err != nil {
			return nil, err
		}
	}

	return &agent, nil
}

// UpdateAgent updates an agent in the database
func UpdateAgent(id string, agent *models.Agent) error {
	_, err := db.Exec("UPDATE agents SET hostname = $1, ip_address = $2, os = $3, os_version = $4, agent_version = $5, last_seen = $6 WHERE host_id = $7",
		agent.Hostname, agent.IPAddress, agent.OS, agent.OSVersion, agent.AgentVersion, time.Now(), id)
	return err
}

// DeleteAgent deletes an agent from the database
func DeleteAgent(id string) error {
	_, err := db.Exec("DELETE FROM agents WHERE host_id = $1", id)
	return err
}
