package agent

import (
	"encoding/json"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/NexClipper/logger"
)

func (agent *KlevrAgent) tempHealthCheck() {
	uri := agent.Manager + "/agents/" + agent.AgentKey + "/tempHeartBeat"

	logger.Debugf(agent.AgentKey)

	rb := &common.Body{}

	SendMe(rb)

	logger.Debugf("%v", rb)

	b, err := json.Marshal(rb)
	if err != nil {
		logger.Error(err)
	}

	// put in & get out
	result := communicator.Put_Json_http(uri, b, agent.AgentKey, agent.API_key, agent.Zone)

	if len(result) > 0 {
		var Body common.Body
		err2 := json.Unmarshal(result, &Body)
		if err2 != nil {
			logger.Error(err2)
		}

		if Body.Me.Deleted {
			agent.scheduler.Remove(agent.tempHealthCheck)
		}
	}
}
