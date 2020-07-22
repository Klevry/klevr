package manager

type Header struct {
	APIKey         string `header:"X-API-KEY"`
	AgentKey       string `header:"X-AGENT-KEY"`
	HashCode       string `header:"X-HASH-CODE"`
	ZoneID         uint64 `header:"X-ZONE-ID"`
	SupportVersion string `header:"X-SUPPORT-AGENT-VERSION"`
}

type Primary struct {
	IP      string `json:"ip"`
	Running bool   `json:"running"`
}

type Agent struct {
	IP       string `json:"ip"`
	Running  bool   `json:"running"`
	AgentKey string `json:"agentKey"`
}

type Resource struct {
	Core   int `json:"core"`
	Memory int `json:"memory"`
	Disk   int `json:"disk"`
}

type Task struct {
}
