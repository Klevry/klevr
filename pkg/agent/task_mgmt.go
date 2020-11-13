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

var receivedTasks *[]common.KlevrTask

func Polling(agent *KlevrAgent) {
	uri := agent.Manager + "/agents/" + agent.AgentKey

	rb := &common.Body{}

	SendMe(rb)

	// add agent nodes
	by := readFile(agentsList)

	list := common.BodyAgent{}

	_ = json.Unmarshal(by, &list)

	for i := 0; i < len(list.Nodes); i++ {
		list.Nodes[i].LastAliveCheckTime = &common.JSONTime{Time: time.Now().UTC()}

		for j := 0; j < len(agent.SecondaryIP); j++ {
			if list.Nodes[i].IP == agent.SecondaryIP[i].IP {
				break
			} else {
				agent.SecondaryIP = append(agent.SecondaryIP, Secondary{list.Nodes[i].IP})
			}
		}
	}

	var updateMap = make(map[uint64]common.KlevrTask)

	for _, t := range *receivedTasks {
		updateMap[t.ID] = t
	}

	rb.Agent.Nodes = list.Nodes

	// update task status
	tasks, _ := executor.GetUpdatedTasks()

	for _, t := range tasks {
		updateMap[t.ID] = t
	}

	updateTasks := []common.KlevrTask{}
	for _, value := range updateMap {
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
	result := communicator.Put_Json_http(uri, b, agent.AgentKey, agent.API_key, agent.Zone)

	var body common.Body

	err := json.Unmarshal(result, &body)
	if err != nil {
		logger.Debugf("%v", string(result))
		logger.Error(err)
	}

	provcheck := exec.Command("sh", "-c", "ssh provbee-service busybee beestatus hello > /tmp/con")
	errcheck := provcheck.Run()
	if errcheck != nil {
		logger.Errorf("provbee-service is not running: %v", errcheck)
	}

	hi := readFile("/tmp/con")
	str := strings.TrimRight(string(hi), "\n")

	if strings.Compare(str, "hi") == 0 {
		// change task status
		logger.Debugf("%v", body.Task)
		for i := 0; i < len(body.Task); i++ {
			if body.Task[i].Status == common.WaitPolling || body.Task[i].Status == common.HandOver {
				body.Task[i].Status = common.WaitExec
			}

			logger.Debugf("%v", body.Task[i].ExeAgentChangeable)

			if body.Task[i].ExeAgentChangeable {

			} else {
				logger.Debugf("%v", &body.Task[i])

				executor.RunTask(&body.Task[i])
			}
		}

		receivedTasks = &body.Task
	}

	writeFile(agentsList, body.Agent)
}
