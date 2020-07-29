package manager

import (
	"time"
)

type AgentGroups struct {
	Id        uint      `xorm:"pk autoincr"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt *time.Time
	GroupName string
	UserId    uint
	Platform  string
}

type Agents struct {
	Id        uint      `xorm:"pk autoincr"`
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	AgentKey  string
	GroupId   uint
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

type PrimaryAgents struct {
	GroupId uint `xorm:"pk"`
	// Group          AgentGroups `gorm:"foreignKey:GroupId"`
	AgentId uint `xorm:"pk"`
	// Agents         Agents
}

type ApiAuthentications struct {
	ApiKey    string `xorm:"pk"`
	GroupId   uint
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

type TaskLock struct {
	Task       string `xorm:"PK"`
	InstanceId string
	LockDate   time.Time
}

func (tl *TaskLock) expired() bool {
	// if tl.LockDate.Unix() > time.Now().UTC().Add()
	return true
}
