package agent

import (
	"encoding/json"
	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/NexClipper/logger"
	"os/exec"
	"strings"
	"time"
)

var agentsList = "/tmp/agents"
var executor = common.GetTaskExecutor()

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

	rb.Agent.Nodes = list.Nodes

	// update task status
	tasks, _ := executor.GetUpdatedTasks()

	rb.Task = tasks

	// secondary node 정보 취합
	/**
	agent 정보 받아오기
	task 처리상태 처
	*/

	// body marshal
	b, err := json.Marshal(rb)
	if err != nil {
		logger.Error(err)
	}

	//logger.Debugf("%v", rb)

	// polling API 호출
	result := communicator.Put_Json_http(uri, b, agent.AgentKey, agent.API_key, agent.Zone)

	body := JsonUnmarshal(result)

	provcheck := exec.Command("sh", "-c", "ssh provbee-service busybee beestatus hello > /tmp/con")
	errcheck := provcheck.Run()
	if errcheck != nil {
		logger.Errorf("provbee-service is not running: %v", errcheck)
	}

	hi := readFile("/tmp/con")
	str := strings.TrimRight(string(hi), "\n")

	if strings.Compare(str, "hi") == 0 {
		// change task status
		for i := 0; i < len(body.Task); i++ {
			if body.Task[i].Status == common.WaitPolling || body.Task[i].Status == common.HandOver {
				body.Task[i].Status = common.WaitExec
			}

			logger.Debugf("%v", body.Task[i].ExeAgentChangeable)

			if body.Task[i].ExeAgentChangeable {

			} else {
				logger.Debugf("%v", &body.Task[i])

				executor.RunTask(&body.Task[i])

				resultCom := exec.Command("sh", "-c", "echo $TASK_RESULT > /tmp/result")
				err := resultCom.Run()
				if err != nil{
					logger.Errorf("error to get task result : %v", err)
				}
				taskResult := readFile("/tmp/result")

				logger.Debugf("%v", string(taskResult))

				body.Task[i].Result = string(taskResult)

				deleteFile("/tmp/result")
			}
		}
	}

	writeFile(agentsList, body.Agent)
}