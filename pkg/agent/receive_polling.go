package agent

import (
	"encoding/json"
	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/NexClipper/logger"
)

func Polling(agent *KlevrAgent) {
	uri := agent.Manager + "/agents/" + agent.AgentKey

	rb := &common.Body{}

	SendMe(rb)

	b, err := json.Marshal(rb)
	if err != nil {
		logger.Error(err)
	}

	agent_me := &common.Agent{}
	agent_me.AgentKey = agent.AgentKey

	// rb.Task = make([]common.Task, 1)
	// rb.Task[0].ID = bod.Task[0].ID
	// rb.Task[0].AgentKey = bod.Task[0].AgentKey
	// rb.Task[0].Command = command
	// rb.Task[0].Status = status
	// rb.Task[0].Params = par
	// rb.Task[0].Result = bod.Task[0].Result
	// rb.Task[0].Type = bod.Task[0].Type

	// put in & get out
	result := communicator.Put_Json_http(uri, b, agent.AgentKey, agent.API_key, agent.Zone)

	logger.Debugf("%v", string(result))
}
