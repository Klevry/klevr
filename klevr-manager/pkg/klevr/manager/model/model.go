package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type AgentGroups struct {
	gorm.Model
	GroupName string
	UserId    int64
	Zone      string
	Platform  string
}

type Agents struct {
	gorm.Model
	AgentKey           string
	GroupId            int64
	IsActive           bool
	LastAliveCheckTime time.Time
	LastAccessTime     time.Time
	Ip                 string
	HmacKey            string
	EncKey             string
}

type PrimaryAgents struct {
	GroupId        int64 `gorm:"primary_key"`
	AgentId        int64 `gorm:"primary_key"`
	LastAccessTime time.Time
}
