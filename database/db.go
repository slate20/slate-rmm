package database

import (
	"database/sql"
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
	// Prepare for SQL Statement
	stmt, err := db.Prepare(`
		INSERT INTO agents (hostname, ip_address, os, os_version, agent_version, last_seen)
		VALUES ($1, $2, $3, $4, $5, $6)
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
		agent.AgentVersion,
		time.Now(),
	).Scan(&agent.ID)

	if err != nil {
		return err
	}

	return nil
}
