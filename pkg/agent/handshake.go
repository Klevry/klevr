package agent

import (
	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/NexClipper/logger"
)

/*
send handshake to manager

in: body.me
out: body.me, body.agent.primary
*/
func HandShake(agent *KlevrAgent) common.Primary {
	uri := agent.Manager + "/agents/handshake"

	rb := &common.Body{}

	SendMe(rb)

	logger.Debugf("%v", rb)

	b := JsonMarshal(rb)

	// put in & get out
	result := communicator.Put_Json_http(uri, b, agent.AgentKey, agent.ApiKey, agent.Zone)

	body := JsonUnmarshal(result)

	logger.Debugf("%v", body)
	agent.schedulerInterval = body.Me.CallCycle

	if len(body.Agent.Nodes) > 0 {
		for _, v := range body.Agent.Nodes {
			agent.Agents = append(agent.Agents, v)
		}
	}

	return body.Agent.Primary
}
