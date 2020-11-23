package common_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Klevry/klevr/pkg/common"
)

func getDefaultTask() *common.KlevrTask {
	return &common.KlevrTask{
		ID:                 1,
		ZoneID:             1,
		Name:               "UNIT_TEST",
		TaskType:           common.AtOnce,
		Schedule:           common.JSONTime{},
		AgentKey:           "",
		ExeAgentKey:        "",
		Status:             common.WaitExec,
		Cron:               "",
		UntilRun:           common.JSONTime{},
		Timeout:            0,
		ExeAgentChangeable: false,
		TotalStepCount:     0,
		CurrentStep:        0,
		HasRecover:         false,
		Parameter:          "{}",
		CallbackURL:        "",
		Result:             "",
		FailedStep:         0,
		IsFailedRecover:    false,
		Steps: []*common.KlevrTaskStep{
			&common.KlevrTaskStep{
				ID:          1,
				Seq:         1,
				CommandName: "COMMAND1",
				CommandType: common.INLINE,
				Command:     "echo hello",
				IsRecover:   false,
			},
		},
		ShowLog:   true,
		Log:       "echo ${TASK_ORIGIN_PARAM}\necho ${TASK_RESULT}\nTASK_RESULT='{\"step\"=\"command1\", \"success\"=true}'",
		CreatedAt: common.JSONTime{},
		UpdatedAt: common.JSONTime{},
	}
}

func getIterationTask() *common.KlevrTask {
	return &common.KlevrTask{
		ID:                 2,
		ZoneID:             1,
		Name:               "MONITOR CLUSTER STATUS",
		TaskType:           common.Iteration,
		Schedule:           common.JSONTime{},
		AgentKey:           "",
		ExeAgentKey:        "",
		Status:             common.WaitExec,
		Cron:               "",
		UntilRun:           common.JSONTime{},
		Timeout:            0,
		ExeAgentChangeable: false,
		TotalStepCount:     0,
		CurrentStep:        0,
		HasRecover:         false,
		Parameter:          "{}",
		CallbackURL:        "",
		Result:             "",
		FailedStep:         0,
		IsFailedRecover:    false,
		Steps: []*common.KlevrTaskStep{
			&common.KlevrTaskStep{
				ID:          1,
				Seq:         1,
				CommandName: "COMMAND1",
				CommandType: common.INLINE,
				Command:     "echo hello",
				IsRecover:   false,
			},
		},
		ShowLog:   true,
		Log:       "echo ${TASK_ORIGIN_PARAM}\necho ${TASK_RESULT}\nTASK_RESULT='{\"step\"=\"command1\", \"success\"=true}'",
		CreatedAt: common.JSONTime{},
		UpdatedAt: common.JSONTime{},
	}
}

func TestSingleRunTask(t *testing.T) {
	task := getDefaultTask()

	executor := common.GetTaskExecutor()

	err := executor.RunTask(task)
	assert.NoError(t, err, "RunTask failed.")

	var updatedTask common.KlevrTask

	for {
		updated, cnt := executor.GetUpdatedTasks()

		for i, u := range updated {
			fmt.Println(i, " = ", u)
		}

		assert.LessOrEqual(t, cnt, 1, "Updated task count not matched.")

		if cnt == 1 {
			status := updated[0].Status
			expected := []common.TaskStatus{common.Running, common.Complete}

			assert.Contains(t, expected, status, "Invalid task status")

			updatedTask = updated[0]
		}

		if count := executor.GetRunningTaskCount(); count < 1 {
			break
		}

		time.Sleep(50 * time.Millisecond)
	}

	assert.Equal(t, common.Complete, updatedTask.Status, "")
}

// func TestIterationTask(t *testing.T) {
// 	iterStr := "{\"id\":37,\"zoneId\":17,\"name\":\"MONITOR CLUSTER STATUS\",\"taskType\":\"iteration\",\"schedule\":null,\"agentKey\":\"\",\"exeAgentKey\":\"\",\"status\":\"wait-polling\",\"cron\":\"* * * * *\",\"untilRun\":null,\"timeout\":0,\"exeAgentChangeable\":true,\"totalStepCount\":1,\"currentStep\":0,\"hasRecover\":false,\"parameter\":\"{}\",\"callbackUrl\":\"\",\"result\":\"\",\"failedStep\":0,\"isFailedRecover\":false,\"steps\":[{\"id\":38,\"seq\":1,\"commandName\":\"GET P8S STATUS\",\"commandType\":\"inline\",\"command\":\"JSON_TASK_PARAMS=${TASK_ORIGIN_PARAM}\\nPROVNS=$(echo ${JSON_TASK_PARAMS} | jq -r '.p8s_namespace')\\np8s_status=`ssh provbee-service busybee kps wow $PROVNS` \\nTASK_RESULT=$(echo ${p8s_status})\\n\",\"isRecover\":false}],\"showLog\":false,\"log\":\"\",\"createdAt\":\"2020-11-16T05:17:51.000000Z\",\"updatedAt\":\"2020-11-16T05:17:51.000000Z\"}"

// 	var iterTask common.KlevrTask
// 	err := json.Unmarshal([]byte(iterStr), &iterTask)

// 	fmt.Println(err)
// 	fmt.Println(iterTask)
// 	fmt.Printf("iter task : [%+v]\n\n", iterTask)

// 	task := &iterTask

// 	// curTime := time.Now()

// 	task.TaskType = common.Iteration
// 	task.Cron = "* * * * *"
// 	// task.UntilRun = common.JSONTime{curTime.Add(3 * time.Minute)}

// 	fmt.Printf("before task : [%+v]\n\n", iterTask)

// 	executor := common.GetTaskExecutor()

// 	err = executor.RunTask(task)
// 	assert.NoError(t, err, "RunTask failed.")

// 	var updatedTask common.KlevrTask

// 	for {
// 		updated, cnt := executor.GetUpdatedTasks()

// 		assert.LessOrEqual(t, cnt, 1, "Updated task count not matched.")

// 		if cnt == 1 {
// 			status := updated[0].Status
// 			// fmt.Println(status)
// 			expected := []common.TaskStatus{common.Running, common.Complete, common.WaitInterationSchedule}

// 			assert.Contains(t, expected, status, "Invalid task status - "+status)

// 			updatedTask = updated[0]

// 			fmt.Printf("updated : [%+v]\n", updated[0])
// 		}

// 		if count := executor.GetRunningTaskCount(); count < 1 {
// 			break
// 		}

// 		time.Sleep(50 * time.Millisecond)
// 	}

// 	assert.Equal(t, common.Complete, updatedTask.Status, "")
// }

// func TestMultiTaskRunTask(t *testing.T) {
// 	iterStr := "{\"task\":[{\"id\":37,\"zoneId\":17,\"name\":\"MONITOR CLUSTER STATUS\",\"taskType\":\"iteration\",\"schedule\":null,\"agentKey\":\"\",\"exeAgentKey\":\"\",\"status\":\"wait-polling\",\"cron\":\"* * * * *\",\"untilRun\":null,\"timeout\":0,\"exeAgentChangeable\":true,\"totalStepCount\":1,\"currentStep\":0,\"hasRecover\":false,\"parameter\":\"{}\",\"callbackUrl\":\"\",\"result\":\"\",\"failedStep\":0,\"isFailedRecover\":false,\"steps\":[{\"id\":38,\"seq\":1,\"commandName\":\"GET P8S STATUS\",\"commandType\":\"inline\",\"command\":\"JSON_TASK_PARAMS=${TASK_ORIGIN_PARAM}\nPROVNS=$(echo ${JSON_TASK_PARAMS} | jq -r '.p8s_namespace')\np8s_status=`ssh provbee-service busybee kps wow $PROVNS` \nTASK_RESULT=$(echo ${p8s_status})\n\",\"isRecover\":false}],\"showLog\":false,\"log\":\"\",\"createdAt\":\"2020-11-16T05:17:51.000000Z\",\"updatedAt\":\"2020-11-16T05:17:51.000000Z\"}]}"

// 	var iterTask common.KlevrTask
// 	err := json.Unmarshal([]byte(iterStr), &iterTask)

// 	fmt.Printf("iter task : [%+v]", iterTask)

// 	task1 := &iterTask
// 	task2 := getDefaultTask()
// 	task3 := getDefaultTask()

// 	curTime := time.Now()

// 	task1.ID = 1
// 	task1.Name = "T1"
// 	task1.TaskType = common.Iteration
// 	task1.Cron = "* * * * *"
// 	task1.UntilRun = common.JSONTime{curTime.Add(1 * time.Minute)}

// 	task2.ID = 2
// 	task2.Name = "T2"

// 	task3.ID = 3
// 	task3.Name = "T3"

// 	executor := common.GetTaskExecutor()

// 	err = executor.RunTask(task1)
// 	assert.NoError(t, err, "RunTask failed.")

// 	err = executor.RunTask(task2)
// 	assert.NoError(t, err, "RunTask failed.")

// 	err = executor.RunTask(task3)
// 	assert.NoError(t, err, "RunTask failed.")

// 	fmt.Println("running count : ", executor.GetRunningTaskCount())

// 	var updatedTask common.KlevrTask

// 	for {
// 		updated, cnt := executor.GetUpdatedTasks()

// 		assert.LessOrEqual(t, cnt, 3, "Updated task count not matched.")

// 		for _, ut := range updated {
// 			fmt.Printf("updated : [%+v]", ut)
// 			fmt.Println()

// 			status := ut.Status
// 			var expected []common.TaskStatus

// 			if ut.TaskType == common.Iteration {
// 				expected = []common.TaskStatus{common.Running, common.Complete, common.WaitInterationSchedule}
// 			} else {
// 				expected = []common.TaskStatus{common.Running, common.Complete}
// 			}

// 			assert.Contains(t, expected, status, "Invalid task status - "+status)
// 		}

// 		if count := executor.GetRunningTaskCount(); count < 1 {
// 			break
// 		}

// 		time.Sleep(50 * time.Millisecond)
// 	}

// 	assert.Equal(t, common.Complete, updatedTask.Status, "")
// }
