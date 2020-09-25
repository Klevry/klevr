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

// TaskType for Task struct
type TaskType string

// TaskStatus for Task struct
type TaskStatus string

type CommandType string

const (
	AtOnce    = TaskType("atOnce")    // 한번만 실행
	Iteration = TaskType("iteration") // 반복 수행(with condition)
	LongTerm  = TaskType("longTerm")  // 장기간 수행
)

const (
	Complete = TaskStatus("complete") // Task 수행 완료
)

// Define TaskTypes
const (
	RESERVED = CommandType("reserved") // 지정된 예약어(커맨드)
	INLINE   = CommandType("inline")   // CLI inline 커맨드
)

type Tasks struct {
	Id          uint64 `xorm:"PK"`
	ZoneId      uint64
	Name        string
	TaskType    TaskType
	Schedule    time.Time
	AgentKey    string
	ExeAgentKey string
	Status      TaskStatus
	TaskDetail  *TaskDetail  `xorm:"foreignKey:Id"`
	TaskSteps   *[]TaskSteps `xorm:"foreignKey:Id"`
	Logs        *TaskLogs    `xorm:"foreignKey:Id"`
	Result      string
	CreatedAt   time.Time `xorm:"created"`
	UpdatedAt   time.Time `xorm:"updated"`
	DeletedAt   time.Time
}

type TaskLogs struct {
	TaskId uint64 `xorm:"PK"`
	Logs   string
}

type TaskSteps struct {
	TaskId          uint64 `xorm:"PK"`
	CommandName     string
	CommandType     CommandType
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
}

func (tl *TaskLock) expired() bool {
	// if tl.LockDate.Unix() > time.Now().UTC().Add()
	return true
}
