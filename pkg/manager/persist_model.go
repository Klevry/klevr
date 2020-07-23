package manager

import (
	"time"
)

type AgentGroups struct {
	Id        uint      `xorm:"pk autoincr"`
	CreatedAt time.Time `xorm:created`
	UpdatedAt time.Time `xorm:updated`
	DeletedAt *time.Time
	GroupName string
	UserId    uint
	Platform  string
}

type Agents struct {
	Id        uint      `xorm:"pk autoincr"`
	CreatedAt time.Time `xorm:created`
	UpdatedAt time.Time `xorm:updated`
	AgentKey  string
	GroupId   uint
	// Group              AgentGroups `gorm:"foreignkey:GroupID"`
	IsActive           bool
	LastAliveCheckTime time.Time
	LastAccessTime     time.Time
	Ip                 string
	HmacKey            string
	EncKey             string
}

type PrimaryAgents struct {
	GroupId uint `xorm:"pk"`
	// Group          AgentGroups `gorm:"foreignKey:GroupId"`
	AgentId uint `xorm:"pk"`
	// Agents         Agents
	LastAccessTime time.Time
}
