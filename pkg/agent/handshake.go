package agent

import (
	"encoding/json"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/NexClipper/logger"
)

//var agents_list = "/tmp/agents"

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

	var body common.Body
	err2 := json.Unmarshal(result, &body)
	if err2 != nil {
		logger.Error(err2)
	}

	logger.Debugf("%v", body)
	primary := body.Agent.Primary.IP
	agent.schedulerInterval = body.Me.CallCycle

	if len(body.Agent.Nodes) > 0 {
		for _, v := range body.Agent.Nodes {
			agent.SecondaryIP = append(agent.SecondaryIP, Secondary{IP: v.IP})
		}
	}

	return primary
}
