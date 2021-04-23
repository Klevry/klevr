package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/NexClipper/logger"
)

func init() {
	// 샘플 커맨드
	// InitCommand(testCommand())

	InitCommand(changeLogLevel())        // 에이전트 로그레벨 변경
	InitCommand(collectAgentLog())       // 에이전트 로그 수집
	InitCommand(stopTask())              // 실행중인 Task 중지 처리(작업 취소)
	InitCommand(forceShutdownAgent())    // 에이전트 즉시 중지 처리
	InitCommand(gracefulShutdownAgent()) // 현재 실행중인 작업을 마치고 에이전트 중지 처리
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

func changeLogLevel() Command {
	var levels = []string{
		"DEBUG", "INFO", "WARN", "ERROR",
	}

	type ParamModel struct {
		Level interface{} `json:"level"`
	}

	type ResultModel struct {
		Before  interface{} `json:"before"`
		Current interface{} `json:"current"`
	}

	return Command{
		Name: "ChangeLogLevel",
		Description: "Change the log level of the agent.\n" +
			"After performing the command, the log level before change and the current log level are returned.",
		ParameterModel: ParamModel{
			Level: CommandDescriptor{
				Type:        "string",
				Description: "",
				Values:      levels,
			},
		},
		ResultModel: ResultModel{
			Before: CommandDescriptor{
				Type:        "string",
				Description: "Log level before change",
			},
			Current: CommandDescriptor{
				Type:        "string",
				Description: "Log level after change",
			},
		},
		Run: func(jsonOriginalParam string, jsonPreResult string) (string, error) {
			p := ParamModel{}

			err := json.Unmarshal([]byte(jsonOriginalParam), &p)
			if err != nil {
				return "", err
			}

			before := logger.GetLevel()

			var level logger.Level

			switch strings.ToLower(p.Level.(string)) {
			case "debug":
				level = 0
			case "info":
				level = 1
			case "warn", "warning":
				level = 2
			case "error":
				level = 3
			case "fatal":
				level = 4
			}

			logger.SetLevel(level)
			current := logger.GetLevel()

			result := ResultModel{
				Before:  before,
				Current: current,
			}

			jsonBytes, err := json.Marshal(result)
			if err != nil {
				return "", err
			}

			return string(jsonBytes), nil
		},
	}
}

// 에이전트의 log 수집
func collectAgentLog() Command {
	var sortParamValues = []string{
		"OLDEST", "LATEST",
	}

	type ParamModel struct {
		Sort interface{} `json:"sort"`
	}

	type ResultModel struct {
		Logs interface{} `json:"logs"`
	}

	return Command{
		Name: "CollectAgentLog",
		Description: "Collect agent logs.\n" +
			"If the log size exceeds 64KB, only 64KB is collected according to the 'sort' parameter.",
		ParameterModel: ParamModel{
			Sort: CommandDescriptor{
				Type: "string",
				Description: "If the log size exceeds 64KB, collect between the oldest and the latest.\n" +
					"If not specified, latest is the default.",
				Values: sortParamValues,
			},
		},
		ResultModel: ResultModel{
			Logs: CommandDescriptor{
				Type:        "string",
				Description: "Agent's log strings.",
			},
		},
		Run: func(jsonOriginalParam string, jsonPreResult string) (string, error) {
			data, err := ioutil.ReadFile(LoggerEnvironment.LogPath)

			logger.Debugf("log path : [%s]", LoggerEnvironment.LogPath)
			logger.Debugf("read original : [%d]", len(string(data)))

			if err != nil {
				return "", err
			}

			baseSize := 30 * 1024
			length := len(data)

			if length > 60*1024 {
				p := ParamModel{}

				err := json.Unmarshal([]byte(jsonOriginalParam), &p)
				if err != nil {
					return "", err
				}

				sort := p.Sort.(string)

				sIndex := 0
				eIndex := 0

				if sort == "OLDEST" {
					eIndex = baseSize - 1
				} else {
					eIndex = length - 1
					sIndex = eIndex - baseSize
				}

				cutData := make([]byte, baseSize)

				copy(cutData[:], data[sIndex:eIndex])

				data = cutData
			}

			logger.Debugf("CollectAgentLog : [%d]", len(string(data)))

			result := ResultModel{
				Logs: string(data),
			}

			jsonBytes, err := json.Marshal(result)
			if err != nil {
				return "", err
			}

			return string(jsonBytes), nil
		},
	}
}

// 실행중인 Task 중지 처리(작업 취소)
func stopTask() Command {
	var stopTaskResultValues = []string{
		"NOT_EXIST", "SUCCESS",
	}

	type ParamModel struct {
		TargetTaskID interface{} `json:"targetTaskId"`
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

			tw, exist := executor.getTaskWrapper(taskID)

			if !exist {
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

func gracefulShutdownAgent() Command {
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

				for _, e := range executor.runningTasks.Items() {
					tw := e.(*TaskWrapper)

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
