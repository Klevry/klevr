package manager

import (
	"time"

	"github.com/jinzhu/gorm"
)

type AgentGroups struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	GroupName string
	UserId    uint
	Platform  string
}

type Agents struct {
	gorm.Model
	AgentKey string
	GroupID  uint
	// Group              AgentGroups `gorm:"foreignkey:GroupID"`
	IsActive           bool
	LastAliveCheckTime time.Time
	LastAccessTime     time.Time
	Ip                 string
	HmacKey            string
	EncKey             string
}

type PrimaryAgents struct {
	GroupId uint `gorm:"primary_key"`
	// Group          AgentGroups `gorm:"foreignKey:GroupId"`
	AgentId uint `gorm:"primary_key"`
	// Agents         Agents
	LastAccessTime time.Time
}
