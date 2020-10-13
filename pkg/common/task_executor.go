package common

import (
	"sync"

	concurrent "github.com/fanliao/go-concurrentMap"
)

type taskExecutor struct {
	sync.RWMutex
	runningTasks *concurrent.ConcurrentMap // 실행중인 TASK map
	updatedTasks *concurrent.ConcurrentMap // 업데이트된 TASK map
	taskChannel  chan *KlevrTask
	alive        bool
}

func NewTaskExecutor() *taskExecutor {
	executor := &taskExecutor{
		runningTasks: concurrent.NewConcurrentMap(),
		updatedTasks: concurrent.NewConcurrentMap(),
		taskChannel:  make(chan *KlevrTask),
		alive:        true,
	}

	go executor.handle()

	return executor
}

func (executor *taskExecutor) handle() {
	for executor.alive {
		select {
		case t := <-executor.taskChannel:
			if t != nil {
				if t.
			}
		}
	}
}

func (executor *taskExecutor) RunTask(task *KlevrTask) {
	_, err := executor.runningTasks.Put(task.ID, task)
	if err != nil {
		panic(err)
	}
}

func execute() {

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

	tasks := make([]KlevrTask, size)

	if size > 0 {
		for i, e := range m.ToSlice() {
			v, _ := m.Remove(e.Key())

			tasks[i] = v.(KlevrTask)
		}
	}

	return tasks, size
}

// Close executor를 종료한다.
func (executor *taskExecutor) Close() {
	executor.alive = false
	executor.taskChannel <- nil
	close(executor.taskChannel)
}
