package models

// Agent represents a system that the RMM tool will monitor.
type Agent struct {
	ID           int64  `json:"host_id"`
	Hostname     string `json:"hostname"`
	IPAddress    string `json:"ip_address"`
	OS           string `json:"os"`
	OSVersion    string `json:"os_version"`
	AgentVersion string `json:"agent_version"`
	LastSeen     string `json:"last_seen"`
}
