package manager_test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/Klevry/klevr/pkg/manager"
	"github.com/NexClipper/logger"
)

func TestToTasks(t *testing.T) {
	rts := make([]manager.RetriveTask, 0)

	rts = append(rts, manager.RetriveTask{
		Tasks:      &manager.Tasks{Id: 1},
		TaskDetail: &manager.TaskDetail{TaskId: 1},
	})

	rts = append(rts, manager.RetriveTask{
		Tasks:      &manager.Tasks{Id: 2},
		TaskDetail: &manager.TaskDetail{TaskId: 2},
	})

	var nrts *[]manager.RetriveTask
	nrts = &rts

	var tasks = make([]manager.Tasks, 0)
	var tasks2 *[]manager.Tasks

	for _, rt := range *nrts {
		// logger.Debugf("retreive task : [%+v], %v", rt.Tasks, unsafe.Pointer(&rt.Tasks))
		// logger.Debugf("retreive taskDetail : [%+v], %v", rt.TaskDetail, unsafe.Pointer(&rt.TaskDetail))
		// logger.Debugf("retreive taskDetail2 : [%+v], %v", rt.Tasks.TaskDetail, unsafe.Pointer(&rt.Tasks.TaskDetail))

		// fmt.Println()

		rt.Tasks.TaskDetail = rt.TaskDetail

		// logger.Debugf("retreive task : [%+v], %v", rt.Tasks, unsafe.Pointer(&rt.Tasks))
		// logger.Debugf("retreive taskDetail : [%+v], %v", rt.TaskDetail, unsafe.Pointer(&rt.TaskDetail))
		// logger.Debugf("retreive taskDetail2 : [%+v], %v", rt.Tasks.TaskDetail, unsafe.Pointer(&rt.Tasks.TaskDetail))

		// fmt.Println()

		tasks = append(tasks, *rt.Tasks)
	}

	tasks2 = &tasks

	for _, t := range *tasks2 {
		logger.Debugf("[%+v], %v", t, unsafe.Pointer(&t))
		logger.Debugf("[%+v], %v", t.TaskDetail, unsafe.Pointer(&t.TaskDetail))
		fmt.Println()
	}
}
