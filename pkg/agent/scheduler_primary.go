package agent

import (
	"encoding/json"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/NexClipper/logger"
)

var receivedTasks []common.KlevrTask = make([]common.KlevrTask, 0)
var notSentTasks map[uint64]common.KlevrTask = make(map[uint64]common.KlevrTask)

func polling(agent *KlevrAgent) {
	if agent.taskPollingPause == true {
		logger.Debug("Polling aborted because authentication failed.")
		return
	}
	executor := common.GetTaskExecutor()
	uri := agent.Manager + "/agents/" + agent.AgentKey

	rb := &common.Body{}
	agent.setBodyMeInfo(rb)

	agent.primaryStatusCheck()

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
	b := jsonMarshal(rb)

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

	for i := 0; i < len(body.Task); i++ {
		if body.Task[i].Status == common.WaitPolling || body.Task[i].Status == common.HandOver {
			body.Task[i].Status = common.WaitExec
		}

		logger.Debugf("%v", body.Task[i].ExeAgentChangeable)

		if body.Task[i].ExeAgentChangeable {
			body.Task[i].ExeAgentKey = agent.AgentKey
			executor.RunTask(agent.AgentKey, &body.Task[i])
		} else {
			logger.Debugf("%v", &body.Task[i])

			if len(body.Task[i].AgentKey) > 0 {
				sendCompleted := false
				for _, v := range agent.Agents {
					if v.AgentKey == body.Task[i].AgentKey {
						body.Task[i].ExeAgentKey = v.AgentKey

						if v.AgentKey == body.Agent.Primary.AgentKey {
							executor.RunTask(agent.AgentKey, &body.Task[i])
						} else {
							ip := v.IP
							t := jsonMarshal(&body.Task[i])
							logger.Debugf("%v", body.Task[i])
							agent.primaryTaskSend(ip, t)
						}

						sendCompleted = true
						break
					}
				}
				if sendCompleted == false {
					body.Task[i].ExeAgentKey = agent.AgentKey
					executor.RunTask(agent.AgentKey, &body.Task[i])
				}
			} else {
				body.Task[i].ExeAgentKey = agent.AgentKey
				executor.RunTask(agent.AgentKey, &body.Task[i])
			}

		}
	}

	receivedTasks = body.Task
}
