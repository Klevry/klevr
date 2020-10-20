package common

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

func init() {
	// 샘플 커맨드
	// InitCommand(testCommand())

	InitCommand(stopTask())               // 실행중인 Task 중지 처리(작업 취소)
	InitCommand(forceShutdownAgent())     // 에이전트 즉시 중지 처리
	InitCommand(gracefuleShutdownAgent()) // 현재 실행중인 작업을 마치고 에이전트 중지 처리
}

// 샘플 커맨드 작성
func testCommand() Command {
	return Command{
		// 커맨드 명칭
		Name: "SampleCommand",
		// 커맨드 실행 로직 구현
		Run: func(jsonOriginalParam string, jsonPreResult string) (string, error) {
			return "", nil
		},
		// 커맨드 중지 시 복구 로직 구현
		Recover: func(jsonOriginalParam string, jsonPreResult string) (string, error) {
			return "", nil
		},
	}
}

// 실행중인 Task 중지 처리(작업 취소)
func stopTask() Command {
	var stopTaskResultValues = []string{
		"NOT_EXIST", "SUCCESS",
	}

	type ParamModel struct {
		TargetTaskID interface{} `json:"targetTaskID"`
	}

	type ResultModel struct {
		StopTaskResult interface{} `json:"stopTaskResult"`
	}

	return Command{
		// 커맨드 명칭
		Name: "StopTask",
		Description: "Stop other running task using targetTaskID as parameter.\n" +
			"Even if the request to stop a running task is successful, the task may be completed before stopping.",
		ParameterModel: ParamModel{
			TargetTaskID: CommandDescriptor{
				Type:        "uint64",
				Description: "task ID to be stopped",
				Values:      "",
			},
		},
		ResultModel: ResultModel{
			StopTaskResult: CommandDescriptor{
				Type:        "string",
				Description: "task execution result",
				Values:      stopTaskResultValues,
			},
		},
		// 커맨드 실행 로직 구현
		Run: func(jsonOriginalParam string, jsonPreResult string) (string, error) {
			p := ParamModel{}

			err := json.Unmarshal([]byte(jsonOriginalParam), &p)
			if err != nil {
				return "", err
			}

			if p.TargetTaskID == nil {
				return "", errors.New("taskId does not exist in parameter for execution command")
			}

			taskID := uint64(p.TargetTaskID.(float64))

			executor := GetTaskExecutor()

			tw, err := executor.getTaskWrapper(taskID)
			if err != nil {
				return "", err
			}

			if tw == nil {
				b, err := json.Marshal(ResultModel{StopTaskResult: stopTaskResultValues[0]})
				if err != nil {
					return "", err
				}

				return string(b), nil
			} else {
				err := tw.future.Cancel()
				if err != nil {
					return "", err
				}

				b, err := json.Marshal(ResultModel{StopTaskResult: stopTaskResultValues[1]})
				if err != nil {
					return "", err
				}

				return string(b), nil
			}
		},
	}
}

func forceShutdownAgent() Command {
	return Command{
		Name:        "ForceShutdownAgent",
		Description: "Ignore the ongoing task and stop the agent immediately.",
		Run: func(jsonOriginalParam string, jsonPreResult string) (string, error) {
			os.Exit(0)

			return "", nil
		},
	}
}

func gracefuleShutdownAgent() Command {
	return Command{
		Name: "GracefulShutdownAgent",
		Description: "When all ongoing AtOnce tasks are finished, the agent is terminated.\n" +
			"Tasks other than AtOnce are forcibly terminated.",
		Run: func(jsonOriginalParam string, jsonPreResult string) (string, error) {
			executor := GetTaskExecutor()

			executor.closed = true

			var complete = false

			for !complete {
				complete = true

				for _, e := range executor.runningTasks.ToSlice() {
					tw := e.Value().(*TaskWrapper)

					if AtOnce == tw.TaskType {
						complete = false
						break
					}
				}

				time.Sleep(1 * time.Second)
			}

			os.Exit(0)

			return "", nil
		},
	}
}
