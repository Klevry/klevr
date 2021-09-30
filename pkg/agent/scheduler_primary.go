package agent

import (
	"encoding/json"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/NexClipper/logger"
)

var receivedTasks []common.KlevrTask = make([]common.KlevrTask, 0)
var notSentTasks map[uint64]common.KlevrTask = make(map[uint64]common.KlevrTask)

// agentkey가 지정되었지만 실행하지 못한 Task는 실패로 처리
func (agent *KlevrAgent) assignmentTask(primaryAgentKey string, task []common.KlevrTask) {
	executor := common.GetTaskExecutor()

	for i := 0; i < len(task); i++ {
		beforeStatus := task[i].Status
		if task[i].Status == common.WaitPolling || task[i].Status == common.HandOver {
			task[i].Status = common.WaitExec
		}

		logger.Debugf("%v", task[i].ExeAgentChangeable)

		if task[i].ExeAgentChangeable {
			task[i].ExeAgentKey = agent.AgentKey
			executor.RunTaskInLocal(&task[i])
		} else {
			logger.Debugf("%v", &task[i])

			if len(task[i].AgentKey) > 0 {
				for _, v := range agent.Agents {
					if v.AgentKey == task[i].AgentKey {
						task[i].ExeAgentKey = v.AgentKey
						if v.AgentKey == primaryAgentKey {
							executor.RunTaskInLocal(&task[i])
						} else {
							ip := v.IP
							logger.Debugf("%v", task[i])
							if err := executor.RunTaskInRemote(ip, agent.grpcPort, &task[i]); err != nil {
								task[i].ExeAgentKey = ""
								task[i].Status = common.TaskStatus(beforeStatus)
							}
						}

						break
					}
				}
			} else {
				task[i].ExeAgentKey = agent.AgentKey
				executor.RunTaskInLocal(&task[i])
			}

		}
	}
}

func (agent *KlevrAgent) polling() {
	if agent.taskPollingPause == true {
		logger.Debug("Polling aborted because authentication failed.")
		return
	}

	executor := common.GetTaskExecutor()
	uri := agent.Manager + "/agents/" + agent.AgentKey

	rb := &common.Body{}
	agent.setBodyMeInfo(rb)

	agent.checkZoneStatus()

	var updateMap = make(map[uint64]common.KlevrTask)

	for _, t := range receivedTasks {
		updateMap[t.ID] = t
	}

	// 전송하지 못 했던 task
	for k, v := range notSentTasks {
		updateMap[k] = v
		delete(notSentTasks, k)
	}

	rb.Agent.Nodes = agent.Agents

	// update task status
	tasks, _ := executor.GetUpdatedTasks()

	for _, t := range tasks {
		updateMap[t.ID] = t
	}

	updateTasks := []common.KlevrTask{}
	for _, value := range updateMap {
		logger.Debugf("polling updated task [%+v]", value)
		updateTasks = append(updateTasks, value)
	}

	rb.Task = updateTasks

	// secondary node 정보 취합
	/**
	agent 정보 받아오기
	task 처리상태 처
	*/

	// body marshal
	b := common.JsonMarshal(rb)

	// polling API 호출
	// polling은 5초마다 시도되는 작업으로 요청이 실패하면 다음 작업을 기다린다.(retryCount가 0인 이유)
	httpHandler := communicator.Http{
		URL:        uri,
		AgentKey:   agent.AgentKey,
		APIKey:     agent.ApiKey,
		ZoneID:     agent.Zone,
		RetryCount: 0,
		Timeout:    agent.HttpTimeout,
	}
	result, err := httpHandler.PutJson(b)
	if err != nil {
		if serr, ok := err.(*common.HTTPError); ok {
			if serr.StatusCode() == 401 {
				agent.taskPollingPause = true
			}
		}
		for k, v := range updateMap {
			notSentTasks[k] = v
		}
		logger.Debugf("Polling url:%s, agent:%s, api:%s, zone:%s", uri, agent.AgentKey, agent.ApiKey, agent.Zone)
		logger.Error(err)
		return
	}

	var body common.Body

	err = json.Unmarshal(result, &body)
	if err != nil {
		for k, v := range updateMap {
			notSentTasks[k] = v
		}
		logger.Debugf("%v", string(result))
		logger.Error(err)
		return
	}

	defer func() {
		agent.Agents = body.Agent.Nodes
	}()

	agent.schedulerInterval = body.Me.CallCycle

	// change task status
	logger.Debugf("%+v", body.Task)
	agent.assignmentTask(body.Agent.Primary.AgentKey, body.Task)

	receivedTasks = body.Task
}
