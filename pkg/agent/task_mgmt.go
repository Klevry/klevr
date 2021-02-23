package agent

import (
	"encoding/json"
	"os/exec"
	"strings"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/NexClipper/logger"
)

var agentsList = "/tmp/agents"
var executor = common.GetTaskExecutor()

var receivedTasks []common.KlevrTask = make([]common.KlevrTask, 0)

func Polling(agent *KlevrAgent) {
	uri := agent.Manager + "/agents/" + agent.AgentKey

	rb := &common.Body{}
	agent.SendMe(rb)

	for i := 0; i < len(agent.Agents); i++ {
		agent.Agents[i].LastAliveCheckTime = &common.JSONTime{Time: time.Now().UTC()}
	}

	var updateMap = make(map[uint64]common.KlevrTask)

	for _, t := range receivedTasks {
		updateMap[t.ID] = t
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
	b := JsonMarshal(rb)

	//logger.Debugf("%v", rb)

	// polling API 호출
	// polling은 5초마다 시도되는 작업으로 요청이 실패하면 다음 작업을 기다린다.(retryCount가 1인 이유)
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
		logger.Debugf("Polling url:%s, agent:%s, api:%s, zone:%s", uri, agent.AgentKey, agent.ApiKey, agent.Zone)
		logger.Error(err)
	}

	var body common.Body

	err = json.Unmarshal(result, &body)
	if err != nil {
		logger.Debugf("%v", string(result))
		logger.Error(err)
	}

	provcheck := exec.Command("sh", "-c", "ssh provbee-service busybee beestatus hello > /tmp/con")
	errcheck := provcheck.Run()
	if errcheck != nil {
		logger.Errorf("provbee-service is not running!!!: %v", errcheck)
	}

	hi := ReadFile("/tmp/con")
	str := strings.TrimRight(string(hi), "\n")

	if strings.Compare(str, "hi") == 0 {
		// change task status
		logger.Debugf("%+v", body.Task)

		for i := 0; i < len(body.Task); i++ {
			if body.Task[i].Status == common.WaitPolling || body.Task[i].Status == common.HandOver {
				body.Task[i].Status = common.WaitExec
			}

			logger.Debugf("%v", body.Task[i].ExeAgentChangeable)

			if body.Task[i].ExeAgentChangeable {
				executor.RunTask(&body.Task[i])
			} else {
				logger.Debugf("%v", &body.Task[i])

				sendCompleted := false
				for _, v := range agent.Agents {
					if v.AgentKey == body.Task[i].AgentKey {
						ip := v.IP

						t := JsonMarshal(&body.Task[i])

						logger.Debugf("%v", body.Task[i])
						agent.PrimaryTaskSend(ip, t)
						sendCompleted = true
						break
					}
				}

				if sendCompleted == false {
					executor.RunTask(&body.Task[i])
				}
			}
		}

		receivedTasks = body.Task
	}

	agent.Agents = body.Agent.Nodes
}
