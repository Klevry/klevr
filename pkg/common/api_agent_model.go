package common

import (
	"net/http"

	"github.com/gorilla/context"
)

// CustomHeaderName custom header name
const CustomHeaderName = "CTX-CUSTOM-HEADER"

// CustomHeader header for klevr
type CustomHeader struct {
	APIKey         string `header:"X-API-KEY"`
	AgentKey       string `header:"X-AGENT-KEY"`
	HashCode       string `header:"X-HASH-CODE"`
	ZoneID         uint64 `header:"X-ZONE-ID"`
	SupportVersion string `header:"X-SUPPORT-AGENT-VERSION"`
	Timestamp      int64  `header:"X-TIMESTAMP"`
}

// Body body for message
type Body struct {
	Me    Me        `json:"me"`
	Agent BodyAgent `json:"agent"`
	Task  []Task    `json:"task"`
}

// Me requester
type Me struct {
	IP        string `json:"ip"`
	Port      int    `json:"port"`
	HmacKey   string `json:"hmacKey"`
	EncKey    string `json:"encKey"`
	Version   string `json:"version"`
	CallCycle int    `json:"callCycle"`
	LogLevel  string `json:"logLevel"`
	Resource
}

// BodyAgent agents
type BodyAgent struct {
	Primary Primary `json:"primary"`
	Nodes   []Agent `json:"nodes"`
}

// Primary primary info
type Primary struct {
	AgentKey       string `json:"agentKey"`
	IP             string `json:"ip"`
	Port           int    `json:"port"`
	IsActive       bool   `json:"isActive"`
	LastAccessTime int64  `json:"lastAccessTime"`
}

// Agent agent info
type Agent struct {
	AgentKey           string `json:"agentKey"`
	IsActive           bool   `json:"isActive"`
	LastAliveCheckTime int64  `json:"lastAliveCheckTime"`
	IP                 string `json:"ip"`
	Port               int    `json:"port"`
	Version            string `json:"version"`
	*Resource
}

// Resource agent resource
type Resource struct {
	Core   int `json:"core"`
	Memory int `json:"memory"`
	Disk   int `json:"disk"`
}

// Task tasks
type Task struct {
	ID       uint64                 `json:"id"`
	Type     TaskType               `json:"taskType"`
	Command  string                 `json:"command"`
	AgentKey string                 `json:"agentKey"`
	Status   string                 `json:"status"`
	Params   map[string]interface{} `json:"params"`
	Result   Result                 `json:"result"`
}

// Result result for task struct
type Result struct {
	Success bool                   `json:"success"`
	Params  map[string]interface{} `json:"params"`
	Errors  string                 `json:"errors"`
}

// TaskType for Task struct
type TaskType string

// TaskStatus for Task struct
type TaskStatus string

// Define TaskTypes
const (
	RESERVED = TaskType("reserved") // 지정된 예약어(커맨드)
	INLINE   = TaskType("inline")   // CLI inline 커맨드
)

// Define TaskStatuses
const (
	NEW       = TaskStatus("new")       // 신규 생성
	WAITING   = TaskStatus("waiting")   // agent에 전달 대기중(master에서 holding)
	DELIVERED = TaskStatus("delivered") // agent에 전달됨(수행중)
	CHECKING  = TaskStatus("checking")  // task 상태 확인중(병렬로 수행된 상태확인 task에 의한 상태값)
	RUNNING   = TaskStatus("running")   // 실행중(병렬로 수행된 상태확인 task에 의한 상태값)
	DONE      = TaskStatus("done")      // 정상 완료
	FAILED    = TaskStatus("failed")    // 실패
)

// NewTask constructor for task struct
func NewTask(id uint64, taskType TaskType, command string, agentKey string, status string, params map[string]interface{}) *Task {
	if taskType == RESERVED || taskType == INLINE {
		return &Task{
			ID:       id,
			Type:     taskType,
			Command:  command,
			AgentKey: agentKey,
			Status:   status,
			Params:   params,
		}
	}

	return nil
}

func GetCustomHeader(r *http.Request) *CustomHeader {
	return context.Get(r, CustomHeaderName).(*CustomHeader)
}
