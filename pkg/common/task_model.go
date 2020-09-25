package common

import "time"

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
	ID                 uint64     // TASK ID
	TaskType           TaskType   // Task 수행 타입
	Schedule           time.Time  // Task가 수행될 일정
	Cron               string     // Task 타입이 iteration일 경우 반복 실행 cron 주기
	UntilRun           time.Time  // Task 타입이 iteration일 경우 실행 기한
	Timeout            uint       // Task 실행 timeout 시간(seconds)
	AgentKey           string     // Task가 수행될 에이전트 key
	ExeAgentKey        string     // 실제 task가 수행된 에이전트 key
	ExeAgentChangEable bool       // Task를 수행할 에이전트 변동 가능 여부
	TotalStepCount     int        // 전체 task step 수
	CurrentStep        int        // 현재 진행중인 task step 번호 (대기 or 실행중)
	HasRecover         bool       // recover step 존재 여부
	Parameter          string     // Task 실행 파라미터(JSON)
	Steps              []Step     // Task 실행 step
	CallbackUrl        string     // Task 완료 결과를 전달받을 URL(Klevr manager 외의 별도 등록 서버)
	Result             string     // Task 수행 결과물
	Status             TaskStatus // Task 상태
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          time.Time
}

type Step struct {
	ID          uint64      // TASK STEP ID
	CommandType CommandType // 커맨드 타입
	Command     string      // inline script 또는 예약어
}

type Step struct {
	ID uint64
}
