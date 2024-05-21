package models

type Hardware struct {
	CPU     string `json:"cpu"`
	Memory  string `json:"memory"`
	Storage string `json:"storage"`
}

// Agent represents a system that the RMM tool will monitor.
type Agent struct {
	ID            int32    `json:"host_id"`
	Hostname      string   `json:"hostname"`
	IPAddress     string   `json:"ip_address"`
	OS            string   `json:"os"`
	OSVersion     string   `json:"os_version"`
	HardwareSpecs Hardware `json:"hardware_specs"`
	AgentVersion  string   `json:"agent_version"`
	LastSeen      string   `json:"last_seen"`
	LastUser      string   `json:"last_user"`
	Token         string   `json:"token"`
}
