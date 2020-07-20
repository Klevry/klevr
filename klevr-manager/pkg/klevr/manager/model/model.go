package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type AgentGroups struct {
	gorm.Model
	GroupName string
	UserId    uint64
	Zone      string
	Platform  string
}

type Agents struct {
	gorm.Model
	AgentKey           string
	GroupId            uint64
	Group              AgentGroups
	IsActive           bool
	LastAliveCheckTime time.Time
	LastAccessTime     time.Time
	Ip                 string
	HmacKey            string
	EncKey             string
}

type PrimaryAgents struct {
	GroupId        uint64 `gorm:"primary_key"`
	Group          AgentGroups
	AgentId        uint64 `gorm:"primary_key"`
	Agents         Agents
	LastAccessTime time.Time
}
