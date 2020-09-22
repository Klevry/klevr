package common

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

type Task struct {
	ID             uint64
	AgentKey       string
	ExeAgentKey    string
	TotalStepCount int
	CurrentStep    int
	HasRecover     bool
}

type Command struct {
}

type Step struct {
	ID uint64
	
}
