package common

// CustomHeaderName custom header name
const CustomHeaderName = "CTX-CUSTOM-HEADER"

// CustomHeader header for klevr
type CustomHeader struct {
	APIKey           string `header:"X-API-KEY"`
	AgentKey         string `header:"X-AGENT-KEY"`
	HashCode         string `header:"X-HASH-CODE"`
	ZoneID           uint64 `header:"X-ZONE-ID"`
	SupportVersion   string `header:"X-SUPPORT-AGENT-VERSION"`
	Timestamp        int64  `header:"X-TIMESTAMP"`
	PayloadEncrypted bool   `header:"X-PAYLOAD-ENC"`
}

// Body body for message
type Body struct {
	Me    Me          `json:"me"`
	Agent BodyAgent   `json:"agent"`
	Task  []KlevrTask `json:"task"`
}

// Me requester
type Me struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	HmacKey   string `json:"hmacKey"`
	EncKey    string `json:"encKey"`
	Version   string `json:"version"`
	CallCycle int    `json:"callCycle"`
	LogLevel  string `json:"logLevel"`
	Deleted   bool   `json:"deleted"`
	Resource
}

// BodyAgent agents
type BodyAgent struct {
	Primary Primary `json:"primary"`
	Nodes   []Agent `json:"nodes"`
}

// Primary primary info
type Primary struct {
	AgentKey       string `json:"agentKey"`
	IP             string `json:"ip"`
	Port           int    `json:"port"`
	IsActive       bool   `json:"isActive"`
	LastAccessTime int64  `json:"lastAccessTime"`
}

// Agent agent info
type Agent struct {
	AgentKey           string    `json:"agentKey"`
	IsActive           bool      `json:"isActive"`
	LastAliveCheckTime *JSONTime `json:"lastAliveCheckTime"`
	IP                 string    `json:"ip"`
	Port               int       `json:"port"`
	Version            string    `json:"version"`
	*Resource
}

// Resource agent resource
type Resource struct {
	Core   int `json:"core"`
	Memory int `json:"memory"`
	Disk   int `json:"disk"`
}
