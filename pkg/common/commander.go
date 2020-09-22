package common

import (
	"fmt"
	"sync"

	"github.com/NexClipper/logger"
	"github.com/fanliao/go-promise"
)

// TaskStatusType task status type
type TaskStatusType string

// TaskStatusType const define
const (
	TaskStatusStart    = TaskStatusType("START")
	TaskStatusRunning  = TaskStatusType("RUNNING")
	TaskStatusSuccess  = TaskStatusType("SUCCESS")
	TaskStatusFailed   = TaskStatusType("FAILED")
	TaskStatusCanceled = TaskStatusType("CANCELED")
)

// 커맨드 구현체를 담는 map
var commands = make(map[string]Command)

// 실행중인 task를 담는 map
var taskMap = sync.Map{}

// CommandWrapper 커맨드를 래핑한 task struct
type CommandWrapper struct {
	Command
	id     uint64
	result TaskResult
	param  *map[string]interface{}
	future *promise.Future
}

// Command 종류별로 구현해야 하는 커맨드 struct.
// 사용자는 필요한 Command를 추가로 생성해야 한다.
type Command struct {
	Name    string
	Run     func(*map[string]interface{}) (interface{}, error)
	Recover func(*map[string]interface{}) (interface{}, error)
}

// TaskResult task 수행 결과 struct
// Status는
type TaskResult struct {
	Status TaskStatusType
	Result interface{}
}

// InitCommand 커맨드 추가
func InitCommand(c Command) {
	commands[c.Name] = c
}

// RunCommand 커맨드를 실행
func RunCommand(id uint64, commandName string, param *map[string]interface{}) error {
	c, ok := commands[commandName]
	if !ok {
		return NewStandardError(fmt.Sprintf("%s command does not exist.", commandName))
	}

	// 커맨드 인스턴스 생성
	cw := &CommandWrapper{
		Command: c,
		id:      id,
		result:  TaskResult{},
		param:   param,
	}

	cw.result.Status = TaskStatusStart

	logger.Debug(cw)

	// 커맨드(task) 실행
	go cw.execute()

	return nil
}

func (c *CommandWrapper) execute() {
	callbacked := false

	// task 맵에 현재 커맨드(task) 적재 for 상태관리
	taskMap.Store(c.id, c)

	task := func() (interface{}, error) {
		return c.Run(c.param)
	}

	// task 실행 및 상태 처리
	f := promise.Start(task).OnSuccess(func(v interface{}) {
		c.result.Status = TaskStatusSuccess
	}).OnFailure(func(v interface{}) {
		c.result.Status = TaskStatusFailed
	}).OnComplete(func(v interface{}) {
		callbacked = true
	}).OnCancel(func() {
		callbacked = true
		c.result.Status = TaskStatusCanceled
	})

	c.result.Status = TaskStatusRunning

	c.future = f

	// task 수행 결과 수신
	r, err := f.Get()

	// cancel 처리 시 OnComplete()이 호출 되지 않는 케이스에 대한 보완 처리
	if !callbacked {
		if f.IsCancelled() {
			c.result.Status = TaskStatusCanceled
		} else if err != nil {
			c.result.Status = TaskStatusFailed
			c.result.Result = r
		}
	}

	// task 수행 결과물 적재
	c.result.Result = r

	// TODO: 매니저로 실행 결과 전송 구현

	// task 맵에서 현재 커맨드(task) 삭제
	taskMap.Delete(c.id)
}

// GetTaskResult get task result
func GetTaskResult(ID uint64) (TaskResult, bool) {
	r, ok := taskMap.Load(ID)

	if ok {
		return r.(TaskResult), ok
	}

	return TaskResult{}, ok
}
