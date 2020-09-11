package agent

import (
	"encoding/json"
	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/NexClipper/logger"
)

/*
send handshake to manager

in: body.me
out: body.me, body.agent.primary
*/
func HandShake(agent *KlevrAgent) string {
	uri := agent.Manager + "/agents/handshake"

	rb := &common.Body{}

	SendMe(rb)

	logger.Debugf("%v", rb)

	b, err := json.Marshal(rb)
	if err != nil {
		logger.Error(err)
	}

	// put in & get out
	result := communicator.Put_Json_http(uri, b, agent.AgentKey, agent.API_key, agent.Zone)

	var Body common.Body
	err2 := json.Unmarshal(result, &Body)
	if err2 != nil {
		logger.Error(err2)
	}

	logger.Debugf("%v", Body.Agent.Primary)
	primary := Body.Agent.Primary.IP

	return primary
}
