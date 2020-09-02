package manager

import (
	"time"
)

const layout = "2006-01-02T15:04:05.000000Z"

type JSONTime struct {
	time.Time
}

func (t *JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Time.Format(layout) + `"`), nil
}

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
	EventTime *JSONTime `json:"eventTime"`
	Result    string    `json:"result"`
	Log       string    `json:"log"`
}
