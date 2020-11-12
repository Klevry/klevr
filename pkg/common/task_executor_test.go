package common_test

import (
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

func TestRunTask(t *testing.T) {
	task := getDefaultTask()

	executor := common.GetTaskExecutor()

	err := executor.RunTask(task)
	assert.NoError(t, err, "RunTask failed.")

	var updatedTask common.KlevrTask

	for {
		updated, cnt := executor.GetUpdatedTasks()

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

func TestIterationTask(t *testing.T) {
	task := getDefaultTask()

	curTime := time.Now()

	task.TaskType = common.Iteration
	task.Cron = "* * * * *"
	task.UntilRun = common.JSONTime{curTime.Add(2 * time.Minute)}

	executor := common.GetTaskExecutor()

	err := executor.RunTask(task)
	assert.NoError(t, err, "RunTask failed.")

	var updatedTask common.KlevrTask

	for {
		updated, cnt := executor.GetUpdatedTasks()

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
