package agent

import (
	"encoding/json"
	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/NexClipper/logger"
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

	var body common.Body
	logger.Debugf("%v", string(result))
	err2 := json.Unmarshal(result, &body)
	if err2 != nil {
		logger.Errorf("%v", err2)
	}

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

			//for _, v := range list.Nodes {
			//	if v.AgentKey == body.Task[i].AgentKey {
			//		ip := v.IP
			//
			//		logger.Debugf("%v", body.Task[i])
			//		agent.taskExecute(ip, &body.Task[i])
			//	}
			//}
		}
	}

	//exec.RunTask(body.Task)

	//logger.Debugf("%v", string(result))

	writeFile(agentsList, body.Agent)

	//logger.Debugf("%v", string(result))
}
