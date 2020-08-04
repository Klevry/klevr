package manager

import (
	"time"
)

// AgentGroups model for AGENT_GROUPS
type AgentGroups struct {
	Id        uint64    `xorm:"pk autoincr"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time
	GroupName string
	UserId    uint64
	Platform  string
}

// Agents model for AGENTS
type Agents struct {
	Id        uint64    `xorm:"pk autoincr"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	AgentKey  string
	GroupId   uint64
	// Group              AgentGroups `gorm:"foreignkey:GroupID"`
	IsActive           bool
	LastAliveCheckTime time.Time
	LastAccessTime     time.Time
	Ip                 string
	Port               int
	HmacKey            string
	EncKey             string
	Cpu                int
	Memory             int
	Disk               int
}

// PrimaryAgents model for PIMARY_AGENTS
type PrimaryAgents struct {
	GroupId uint64 `xorm:"pk"`
	// Group          AgentGroups `gorm:"foreignKey:GroupId"`
	AgentId uint64 `xorm:"pk"`
	// Agents         Agents
}

// ApiAuthentications model for API_AUTHENTICATIONS
type ApiAuthentications struct {
	ApiKey    string `xorm:"pk"`
	GroupId   uint64
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

// TaskLock model for TASK_LOCK
type TaskLock struct {
	Task       string `xorm:"PK"`
	InstanceId string
	LockDate   time.Time
}

func (tl *TaskLock) expired() bool {
	// if tl.LockDate.Unix() > time.Now().UTC().Add()
	return true
}
