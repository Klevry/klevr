package manager

import (
	"time"
)

// EventType Klevr event type
type EventType string

// Klevr Event constraints
const (
	AgentConnect    EventType = "AGENT_CONNECT"
	AgentDisconnect EventType = "AGENT_DISCONNECT"
	PrimaryInit     EventType = "PRIMARY_INIT"
	TaskCallback    EventType = "TASK_CALLBACK"
)

// KlevrEvent klevr event struct
type KlevrEvent struct {
	EventType EventType `json:"eventType"`
	AgentId   uint64    `json:"agentId"`
	GroupId   uint64    `json:"groupId"`
	EventTime time.Time `json:"eventTime"`
	Result    string    `json:"result"`
	Log       string    `json:"log"`
}
