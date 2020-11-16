package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/NexClipper/logger"
	"github.com/gorhill/cronexpr"
	"github.com/pkg/errors"

	concurrent "github.com/fanliao/go-concurrentMap"
	"github.com/fanliao/go-promise"
)

var once sync.Once
var tExecutor taskExecutor

// taskExecutor task executor for KlevrTask. Use a constructor GetTaskExecutor() for creation.
type taskExecutor struct {
	sync.RWMutex
	runningTasks *concurrent.ConcurrentMap // 실행중인 TASK map
	updatedTasks *concurrent.ConcurrentMap // 업데이트된 TASK map
	closed       bool
}

// TaskWrapper for running task management
type TaskWrapper struct {
	*KlevrTask
	future  *promise.Future
	recover *KlevrTaskStep
}

// GetTaskExecutor constructor for taskExecutor.
func GetTaskExecutor() *taskExecutor {
	once.Do(func() {
		tExecutor = taskExecutor{
			runningTasks: concurrent.NewConcurrentMap(), // *TaskWrapper
			updatedTasks: concurrent.NewConcurrentMap(), // KlevrTask
		}
	})

	return &tExecutor
}

func (executor *taskExecutor) getTaskWrapper(ID uint64) (*TaskWrapper, error) {
	tw, err := executor.runningTasks.Get(ID)

	return tw.(*TaskWrapper), err
}

// GetRunningTaskCount 현재 진행중인 TASK의 개수를 반환
func (executor *taskExecutor) GetRunningTaskCount() int {
	return int(executor.runningTasks.Size())
}

// GetUpdatedTasks 진행 상태가 변경된 task 조회
func (executor *taskExecutor) GetUpdatedTasks() (updated []KlevrTask, count int) {
	executor.Lock()
	defer executor.Unlock()

	m := executor.updatedTasks
	size := int(m.Size())

	tasks := make([]KlevrTask, 0, size)

	if size > 0 {
		for _, e := range m.ToSlice() {
			v, _ := m.Remove(e.Key())

			tasks = append(tasks, v.(KlevrTask))
		}
	}

	return tasks, len(tasks)
}

// RunTask Run the task.
func (executor *taskExecutor) RunTask(task *KlevrTask) error {
	if executor.closed {
		return errors.New("Task executor was closed")
	}

	tw := &TaskWrapper{KlevrTask: task}
	tw.Status = Running

	_, err := executor.runningTasks.Put(task.ID, tw)
	if err != nil {
		panic(err)
	}

	go executor.execute(tw)

	return nil
}

func (executor *taskExecutor) execute(tw *TaskWrapper) {
	// execute() 종료 시 runningTask에서 삭제 및 상태 업데이트
	defer func() {
		executor.runningTasks.Remove(tw.ID)
		executor.updatedTasks.Put(tw.ID, *tw.KlevrTask)
	}()

	// Promise function definition
	f := func(canceller promise.Canceller) (interface{}, error) {
		steps := tw.Steps
		size := len(steps)

		var result string
		var err error

		if size > 0 {
			// Task step 순서 정렬
			sort.Slice(steps, func(i, j int) bool {
				return steps[i].Seq < steps[j].Seq
			})

			// Recover 스텝을 제외한 정규 step 개수
			regularCnt := size
			if tw.HasRecover {
				regularCnt--
			}

			// Iteration task 반복 수행 시작 지점
		ITERATION:

			tw.CurrentStep = 1

			// Task 실행 시작 상태 업데이트
			executor.updatedTasks.Put(tw.ID, *tw.KlevrTask)

			// Task step 순차 실행
			for i, step := range steps {
				// task cancel 처리
				if canceller.IsCancelled() {
					return tw, nil
				}

				// recover step 처리
				if step.IsRecover {
					// recover step이 2개 이상이면 오류
					if tw.recover != nil {
						return tw, errors.New(fmt.Sprintf("%d Task has two or more recovers.", tw.ID))
					}

					tw.recover = step
					continue
				}

				// task step 실행
				if RESERVED == step.CommandType {
					result, err = runReservedCommand(result, tw.KlevrTask, step)
				} else if INLINE == step.CommandType {
					result, err = runInlineCommand(result, tw.KlevrTask, step)
				} else {
					return tw, errors.New(fmt.Sprintf("%d Task invalid command type - %s", tw.ID, step.CommandType))
				}

				// task result 갱신 - task result는 최종 step의 result로 갱신된다.
				tw.Result = result

				// step의 처리 결과가 error인 경우 task 실행 중지 및 error return -> OnFailure 처리
				if err != nil {
					return tw, err
				}

				// 마지막 step이 아니면 task 진행상황 업데이트
				if i < regularCnt {
					tw.CurrentStep++

					executor.updatedTasks.Put(tw.ID, *tw.KlevrTask)
				}
			}
			fmt.Println("요기1", Iteration == tw.TaskType)
			// Iteration task 반복 수행
			if Iteration == tw.TaskType {
				expr, err := cronexpr.Parse(tw.Cron)

				if err != nil {
					tw.Log += "Invalid cron expression - " + tw.Cron + "\n"
					return tw, err
				}

				curTime := time.Now()
				nextTime := expr.Next(curTime)

				fmt.Println("요기2", tw.UntilRun.After(nextTime))
				if tw.UntilRun.After(nextTime) {
					tw.Status = WaitInterationSchedule
					executor.updatedTasks.Put(tw.ID, *tw.KlevrTask)

					time.Sleep(nextTime.Sub(curTime))

					tw.Status = Running

					fmt.Println("요기3")

					goto ITERATION
				}
			}
		}

		// 최종 task step까지 처리 완료된 경우 정상 return -> OnSuccess 처리
		return tw, nil
	}

	// Promise function 실행 및 handler 정의
	future := promise.Start(f).OnSuccess(func(v interface{}) {
		tw.Status = Complete
	}).OnFailure(func(v interface{}) {
		tw := v.(*TaskWrapper)

		tw.FailedStep = tw.CurrentStep

		if tw.HasRecover {
			var result string
			var err error

			tw.CurrentStep = uint(tw.recover.Seq)
			tw.Status = Recovering
			executor.updatedTasks.Put(tw.ID, *tw.KlevrTask)

			defer func(err error) {
				v := recover()

				if v != nil {
					tw.Status = FailedRecover
					tw.IsFailedRecover = true

					tw.Log += fmt.Sprintf("%+v\n\n", v)
				} else if err != nil {
					tw.Status = FailedRecover
					tw.IsFailedRecover = true

					tw.Log += fmt.Sprintf("%+v\n\n", err)
				}
			}(err)

			if RESERVED == tw.recover.CommandType {
				result, err = runReservedCommand(tw.Result, tw.KlevrTask, tw.recover)
			} else if INLINE == tw.recover.CommandType {
				result, err = runInlineCommand(tw.Result, tw.KlevrTask, tw.recover)
			}

			tw.Result = result
			tw.Status = Failed
		} else {
			tw.Status = Failed
		}
	}).OnCancel(func() {
		tw.Status = Stopped
	})

	tw.future = future

	if tw.Timeout > 0 { // 태스크 실행 with Timeout
		_, err, timeout := future.GetOrTimeout(tw.Timeout * 1000)
		if err != nil {
			logger.Errorf("%+v", errors.WithStack(err))
			tw.Log += fmt.Sprintf("%+v\n\n", errors.WithStack(err))
		}

		if timeout {
			tw.Status = Timeout
		}
	} else { // 태스크 실행 without Timeout
		_, err := future.Get()
		if err != nil {
			logger.Errorf("%+v", errors.WithStack(err))
			tw.Log += fmt.Sprintf("%+v\n\n", errors.WithStack(err))
		}
	}

	// CallbackURL 이 존재하는 경우 비동기 callback 처리
	if tw.CallbackURL != "" {
		d := KlevrTaskCallback{
			ID:     tw.ID,
			Name:   tw.Name,
			Status: tw.Status,
			Result: tw.Result,
		}

		go callback(tw.CallbackURL, d)
	}
}

// 예약 커맨드(golang function) 실행
func runReservedCommand(preResult string, task *KlevrTask, command *KlevrTaskStep) (result string, err error) {
	return RunCommand(preResult, task, command)
}

// inline shell script 커맨드 실행
func runInlineCommand(preResult string, task *KlevrTask, command *KlevrTaskStep) (result string, err error) {
	var wrapper string
	var path = "/tmp/" + strconv.FormatUint(task.ID, 10)

	// inline command 스크립트 파일 생성을 위한 디렉토리 체크(/tmp/taskID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0700)
	}

	// inline command 스크립트 파일 생성
	scriptFile := path + "/" + strconv.Itoa(command.Seq)
	resultFile := scriptFile + ".result"
	ioutil.WriteFile(scriptFile, []byte(command.Command), 0700)
	defer func() {
		err := os.Remove(scriptFile)
		if err != nil {
			logger.Warningf("%s file remove failed - %+v", scriptFile, err)
		}

		os.Remove(resultFile)
	}()

	wrapper += InlineCommandOriginalParamVarName + "=\"" + strings.ReplaceAll(task.Parameter, "\"", "\\\"") + "\"\n"
	wrapper += InlineCommandTaskResultVarName + "=\"" + strings.ReplaceAll(preResult, "\"", "\\\"") + "\"\n\n"

	wrapper += ". " + scriptFile + "\n\n"

	wrapper += "\necho \"${" + InlineCommandTaskResultVarName + "}\" > " + resultFile

	// command wrapper 스크립트 파일 생성
	wrapperFile := path + "/wrapper.sh"
	ioutil.WriteFile(wrapperFile, []byte(wrapper), 0700)
	defer func() {
		err := os.Remove(wrapperFile)
		if err != nil {
			logger.Warningf("%s file remove failed - %+v", wrapperFile, err)
		}
	}()

	// 실행
	cmd := exec.Command("sh", "-c", wrapperFile)

	var stdOut bytes.Buffer
	var errOut bytes.Buffer

	cmd.Stdout = &stdOut
	cmd.Stderr = &errOut

	runErr := cmd.Run()

	if task.ShowLog {
		task.Log += stdOut.String() + "\n\n"
	}

	task.Log += errOut.String() + "\n\n"

	if runErr != nil {
		task.Log += fmt.Sprintf("%+v\n\n", errors.WithStack(runErr))
		return "", runErr
	}

	// 결과 조회
	b, err := ioutil.ReadFile(resultFile)
	if err != nil {
		logger.Errorf("%s file read failed - %+v", resultFile, err)
		task.Log += fmt.Sprintf("%+v\n\n", errors.WithStack(runErr))

		return "", err
	}

	return string(b), nil
}

func callback(url string, data KlevrTaskCallback) {
	b, err := json.Marshal(data)
	if err != nil {
		logger.Errorf("%+v", errors.WithStack(err))
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		logger.Debugf("%+v", errors.WithStack(err))
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err == nil {
		res.Body.Close()
	} else {
		logger.Debugf("%+v", errors.WithStack(err))
		return
	}
}
