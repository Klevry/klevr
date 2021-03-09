package common

import (
	"fmt"

	"github.com/NexClipper/logger"
	"github.com/pkg/errors"
)

// 커맨드 구현체를 담는 map
var commands = make(map[string]Command)

// Command 종류별로 구현해야 하는 커맨드 struct.
// 사용자는 필요한 Command를 추가로 생성해야 한다.
type Command struct {
	Name           string
	Description    string
	ParameterModel interface{}
	ResultModel    interface{}
	Run            func(jsonOriginalParam string, jsonPreResult string) (jsonResult string, err error)
	Recover        func(jsonOriginalParam string, jsonPreResult string) (jsonResult string, err error)
}

type CommandDescriptor struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Values      interface{} `json:"values"`
}

// InitCommand 커맨드 추가
func InitCommand(c Command) {
	commands[c.Name] = c
}

// GetCommands 등록된 command map을 반환
func GetCommands() map[string]Command {
	t := make(map[string]Command)

	for k, v := range commands {
		t[k] = v
	}

	return t
}

// RunCommand 커맨드를 실행
func RunCommand(jsonPreResult string, task *KlevrTask, command *KlevrTaskStep) (result string, err error) {
	// 커맨드 구현체 function 획득
	c, ok := commands[command.Command]
	if !ok {
		return "", NewStandardError(fmt.Sprintf("%s command does not exist.", command.Command))
	}

	var r string
	var e error

	// Recover function 실행 (Run()이 panic 또는 error 일때만)
	defer func() {
		v := recover()
		var ex error

		if v != nil {
			e = v.(error)
			task.Log += fmt.Sprintf("%+v\n\n", errors.WithStack(e))

			if c.Recover != nil {
				r, ex = c.Recover(task.Parameter, task.Result)
			}
		} else if e != nil {
			task.Log += fmt.Sprintf("%+v\n\n", errors.WithStack(e))
			if c.Recover != nil {
				r, ex = c.Recover(task.Parameter, task.Result)
			}
		}

		if ex != nil {
			task.Log += fmt.Sprintf("%+v\n\n", errors.WithStack(ex))
		}
	}()

	// Run function 실행
	r, e = c.Run(task.Parameter, task.Result)

	logger.Debugf("Reserved command result : [%d]", len(r))

	return r, e
}
