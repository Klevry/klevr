package model

import (
	"github.com/Klevry/klevr/pkg/serialize"
)

// TaskType for KlevrTask struct
type TaskType string

// TaskStatus for KlevrTask struct
type TaskStatus string

// CommandType for KlevrTask struct
type CommandType string

// EventHookSendingType for KlevrTask struct
type EventHookSendingType string

// INLINE script reserved variable name constants
const (
	InlineCommandOriginalParamVarName = "TASK_ORIGIN_PARAM"
	InlineCommandTaskResultVarName    = "TASK_RESULT"
)

// TaskType Task 종류 정의
const (
	AtOnce    = TaskType("atOnce")    // 한번만 실행
	Iteration = TaskType("iteration") // 반복 수행(with condition)
	LongTerm  = TaskType("longTerm")  // 장기간 수행
)

// TaskStatus Task 상태 정의
const (
	Scheduled              = TaskStatus("scheduled")               // 실행 스케쥴링됨
	WaitPolling            = TaskStatus("wait-polling")            // 에이전트 polling 대기중
	HandOver               = TaskStatus("hand-over")               // Primary 에이전트로 task 전달
	WaitExec               = TaskStatus("wait-exec")               // Secondary 에이전트로 task 전달 대기중
	Started                = TaskStatus("started")                 // Task가 시작됨
	Running                = TaskStatus("running")                 // Task 실행중
	Recovering             = TaskStatus("recovering")              // 복구중
	WaitInterationSchedule = TaskStatus("wait-iteration-schedule") // 반복 스케쥴 대기중
	Complete               = TaskStatus("complete")                // Task 수행 완료
	FailedRecover          = TaskStatus("failed-recover")          // 복구 실패
	Failed                 = TaskStatus("failed")                  // Task 수행 실패
	Canceled               = TaskStatus("canceled")                // hand-over 전인 Task 실행 취소
	Stopped                = TaskStatus("stopped")                 // 실행중인 Task 취소 (recovery 하지 않음)
	Timeout                = TaskStatus("timeout")                 // Task timeout (recovery 하지 않음)
)

// CommandType Command 종류 정의
const (
	RESERVED = CommandType("reserved") // 지정된 예약어(커맨드)
	INLINE   = CommandType("inline")   // CLI inline 커맨드
)

// EventHookSendingType 이벤트훅 전송 조건 정의
const (
	EventHookWithAll           = EventHookSendingType("all-statues")         // 모든 상태 변화에 전송
	EventHookWithConclusion    = EventHookSendingType("conclusion-statuses") // 완료 상태에 전송
	EventHookWithBothEnds      = EventHookSendingType("both-ends")           // 시작과 완료 상태에 전송
	EventHookWithSuccess       = EventHookSendingType("success")             // 성공일때 전송
	EventHookWithFailed        = EventHookSendingType("failed")              // 실패일때 전송
	EventHookWithChangedResult = EventHookSendingType("changed-result")      // result가 변경 되었을때
	EventHookWithEachSteps     = EventHookSendingType("each-steps")          // 각 task step 단계마다
)

// KlevrTask define Task model
type KlevrTask struct {
	ID                   uint64               `json:"id"`
	ZoneID               uint64               `json:"zoneId"`
	Name                 string               `json:"name"`
	TaskType             TaskType             `json:"taskType"`
	Schedule             serialize.JSONTime   `json:"schedule"`
	AgentKey             string               `json:"agentKey"`
	ExeAgentKey          string               `json:"exeAgentKey"`
	Status               TaskStatus           `json:"status"`
	Cron                 string               `json:"cron"`
	UntilRun             serialize.JSONTime   `json:"untilRun"`
	Timeout              uint                 `json:"timeout"`
	ExeAgentChangeable   bool                 `json:"exeAgentChangeable"`
	TotalStepCount       uint                 `json:"totalStepCount"`
	CurrentStep          uint                 `json:"currentStep"`
	HasRecover           bool                 `json:"hasRecover"`
	Parameter            string               `json:"parameter"`
	CallbackURL          string               `json:"callbackUrl"`
	Result               string               `json:"result"`
	FailedStep           uint                 `json:"failedStep"`
	IsFailedRecover      bool                 `json:"isFailedRecover"`
	Steps                []*KlevrTaskStep     `json:"steps"`
	ShowLog              bool                 `json:"showLog"`
	Log                  string               `json:"log"`
	EventHookSendingType EventHookSendingType `json:"eventHookSendingType"`
	IsChangedResult      bool                 `json:"isChangedResult"`
	CreatedAt            serialize.JSONTime   `json:"createdAt"`
	UpdatedAt            serialize.JSONTime   `json:"updatedAt"`
}

type KlevrTaskStep struct {
	ID          uint64      `json:"id"`
	Seq         int         `json:"seq"`
	CommandName string      `json:"commandName"`
	CommandType CommandType `json:"commandType"`
	Command     string      `json:"command"`
	IsRecover   bool        `json:"isRecover"`
}

type KlevrTaskCallback struct {
	ID     uint64     `json:"id"`
	Name   string     `json:"name"`
	Status TaskStatus `json:"status"`
	Result string     `json:"result"`
}

type KlevrTaskLog struct {
	ID          uint64             `json:"id"`
	ZoneID      uint64             `json:"zoneId"`
	Name        string             `json:"name"`
	ExeAgentKey string             `json:"exeAgentKey"`
	Status      TaskStatus         `json:"status"`
	Log         string             `json:"log"`
	CreatedAt   serialize.JSONTime `json:"createdAt"`
	UpdatedAt   serialize.JSONTime `json:"updatedAt"`
}
