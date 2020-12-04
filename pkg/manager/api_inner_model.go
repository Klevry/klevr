package manager

import (
	"encoding/json"

	"github.com/Klevry/klevr/pkg/common"
)

// EventType Klevr event type
type EventType string

// Klevr Event constraints
const (
	// 에이전트 접속 이벤트
	AgentConnect EventType = "AGENT_CONNECT"
	// 에이전트 접속 해제 이벤트
	AgentDisconnect EventType = "AGENT_DISCONNECT"
	// primary 에이전트 선출 이벤트
	PrimaryElected EventType = "PRIMARY_ELECTED"
	// primary 에이전트 retire 이벤트
	PrimaryRetire EventType = "PRIMARY_RETIRE"
	// task 수행 결과 전달 이벤트
	TaskCallback EventType = "TASK_CALLBACK"
)

type KlevrVariable struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Length      string `json:"length"`
	Description string `json:"description"`
	Example     string `json:"example"`
	Value       string `json:"value"`
}

// KlevrEvent klevr event struct
type KlevrEvent struct {
	EventType EventType        `json:"eventType"`
	AgentKey  string           `json:"agentKey"`
	GroupID   uint64           `json:"groupId"`
	EventTime *common.JSONTime `json:"eventTime"`
	Result    string           `json:"result"`
	Log       string           `json:"log"`
}

type KlevrEventTaskInfo struct {
	ID              uint64            `json:"id"`
	ZoneID          uint64            `json:"zoneId"`
	Name            string            `json:"name"`
	TaskType        common.TaskType   `json:"taskType"`
	AgentKey        string            `json:"agentKey"`
	ExeAgentKey     string            `json:"exeAgentKey"`
	Status          common.TaskStatus `json:"status"`
	TotalStepCount  uint              `json:"totalStepCount"`
	CurrentStep     uint              `json:"currentStep"`
	FailedStep      uint              `json:"failedStep"`
	IsFailedRecover bool              `json:"isFailedRecover"`
	UpdatedAt       *common.JSONTime  `json:"updatedAt"`
}

type KlevrEventTaskResult struct {
	Task             KlevrEventTaskInfo `json:"taskInfo"`
	Complete         bool               `json:"complete"`
	Success          bool               `json:"success"`
	IsCommandError   bool               `json:"isCommandError"`
	Result           string             `json:"result"`
	Log              string             `json:"log"`
	ExceptionMessage string             `json:"exceptionMessage"`
	ExceptionTrace   string             `json:"exceptionTrace"`
}

func NewKlevrEventTaskResultString(task *Tasks, complete bool, success bool, isCommandError bool, result string, log string, exceptionMessage string, exceptionTrace string) string {
	b, err := json.Marshal(KlevrEventTaskResult{
		Task: KlevrEventTaskInfo{
			ID:              task.Id,
			ZoneID:          task.ZoneId,
			Name:            task.Name,
			TaskType:        task.TaskType,
			AgentKey:        task.AgentKey,
			ExeAgentKey:     task.ExeAgentKey,
			Status:          task.Status,
			TotalStepCount:  task.TaskDetail.TotalStepCount,
			CurrentStep:     task.TaskDetail.CurrentStep,
			FailedStep:      task.TaskDetail.FailedStep,
			IsFailedRecover: task.TaskDetail.IsFailedRecover,
			UpdatedAt:       &common.JSONTime{Time: task.UpdatedAt},
		},
		Complete:         complete,
		Success:          success,
		IsCommandError:   isCommandError,
		ExceptionMessage: exceptionMessage,
		ExceptionTrace:   exceptionTrace,
		Result:           result,
		Log:              log,
	})

	if err != nil {
		panic(err)
	}

	return string(b)
}

type ReservedCommand struct {
	Description    string      `json:"description"`
	ParameterModel interface{} `json:"parameterModel"`
	ResultModel    interface{} `json:"resultModel"`
	HasRecover     bool        `json:"hasRecover"`
}

type SimpleReservedCommand struct {
	Parameter string `json:"parameter"`
	Command   string `json:"command"`
}
