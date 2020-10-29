package manager

import (
	"time"

	"github.com/Klevry/klevr/pkg/common"
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
	IsActive           byte
	LastAliveCheckTime time.Time
	LastAccessTime     time.Time
	Ip                 string
	Port               int
	HmacKey            string
	EncKey             string
	Cpu                int
	Memory             int
	Disk               int
	Version            string
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

type RetriveTask struct {
	Tasks      `xorm:"extends"`
	TaskDetail `xorm:"extends"`
	TaskLogs   `xorm:"extends"`
}

func (RetriveTask) TableName() string {
	return "TASKS"
}

type Tasks struct {
	Id          uint64 `xorm:"PK"`
	ZoneId      uint64
	Name        string
	TaskType    common.TaskType
	Schedule    time.Time
	AgentKey    string
	ExeAgentKey string
	Status      common.TaskStatus
	TaskDetail  *TaskDetail  `xorm:"-"`
	TaskSteps   *[]TaskSteps `xorm:"-"`
	Logs        *TaskLogs    `xorm:"-"`
	CreatedAt   time.Time    `xorm:"created"`
	UpdatedAt   time.Time    `xorm:"updated"`
	DeletedAt   time.Time
}

type TaskLogs struct {
	TaskId uint64 `xorm:"PK"`
	Logs   string
}

type TaskSteps struct {
	Id              uint64 `xorm:"PK"`
	Seq             int
	TaskId          uint64
	CommandName     string
	CommandType     common.CommandType
	ReservedCommand string
	InlineScript    string
	IsRecover       bool
}

type TaskDetail struct {
	TaskId             uint64 `xorm:"PK"`
	Cron               string
	UntilRun           time.Time
	Timeout            uint
	ExeAgentChangeable bool
	TotalStepCount     uint
	CurrentStep        uint
	HasRecover         bool
	Parameter          string
	CallbackUrl        string
	Result             string
	FailedStep         uint
	IsFailedRecover    bool
	ShowLog            bool
}

func (tl *TaskLock) expired() bool {
	// if tl.LockDate.Unix() > time.Now().UTC().Add()
	return true
}
