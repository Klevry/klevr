package common

import (
	"time"

	"github.com/NexClipper/logger"
)

// TaskType for KlevrTask struct
type TaskType string

// TaskStatus for KlevrTask struct
type TaskStatus string

// CommandType for KlevrTask struct
type CommandType string

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
	Scheduled     = TaskStatus("scheduled")      // 실행 스케쥴링됨
	WaitPolling   = TaskStatus("wait-polling")   // 에이전트 polling 대기중
	HandOver      = TaskStatus("hand-over")      // Primary 에이전트로 task 전달
	WaitExec      = TaskStatus("wait-exec")      // Secondary 에이전트로 task 전달 대기중
	Running       = TaskStatus("running")        // Task 실행중
	Recovering    = TaskStatus("recovering")     // 복구중
	Complete      = TaskStatus("complete")       // Task 수행 완료
	FailedRecover = TaskStatus("failed-recover") // 복구 실패
	Failed        = TaskStatus("failed")         // Task 수행 실패
	Canceled      = TaskStatus("canceled")       // Task 취소 (recovery 하지 않음)
	Timeout       = TaskStatus("timeout")        // Task timeout (recovery 하지 않음)
)

// CommandType Command 종류 정의
const (
	RESERVED = CommandType("reserved") // 지정된 예약어(커맨드)
	INLINE   = CommandType("inline")   // CLI inline 커맨드
)

// KlevrTask define Task model
type KlevrTask struct {
	ID                 uint64           `json:"id"`
	ZoneID             uint64           `json:"zoneId"`
	Name               string           `json:"name"`
	TaskType           TaskType         `json:"taskType"`
	Schedule           JSONTime         `json:"schedule"`
	AgentKey           string           `json:"agentKey"`
	ExeAgentKey        string           `json:"exeAgentKey"`
	Status             TaskStatus       `json:"status"`
	Cron               string           `json:"cron"`
	UntilRun           JSONTime         `json:"untilRun"`
	Timeout            uint             `json:"timeout"`
	ExeAgentChangeable bool             `json:"exeAgentChangeable"`
	TotalStepCount     uint             `json:"totalStepCount"`
	CurrentStep        uint             `json:"currentStep"`
	HasRecover         bool             `json:"hasRecover"`
	Parameter          string           `json:"parameter"`
	CallbackURL        string           `json:"callbackUrl"`
	Result             string           `json:"result"`
	FailedStep         uint             `json:"failedStep"`
	IsFailedRecover    bool             `json:"isFailedRecover"`
	Steps              []*KlevrTaskStep `json:"steps"`
	ShowLog            bool             `json:"showLog"`
	Log                string           `json:"log"`
	CreatedAt          JSONTime         `json:"createdAt"`
	UpdatedAt          JSONTime         `json:"updatedAt"`
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

func TaskStatusAdd(task *KlevrTask) *KlevrTask {
	if task.Schedule.IsZero() {
		task.Status = WaitPolling
	} else {
		compare := time.Now().UTC()

		if task.Schedule.After(compare) {
			task.Status = Scheduled
		} else {
			task.Status = WaitPolling
		}

		logger.Debugf("schedule : [%+v], current : [%+v], task status : [%+v]", task.Schedule, compare, task.Status)
	}

	return task
}
