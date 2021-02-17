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
func HandShake(agent *KlevrAgent) *common.Primary {
	uri := agent.Manager + "/agents/handshake"

	rb := &common.Body{}
	agent.SendMe(rb)
	logger.Debugf("%v", rb)
	b := JsonMarshal(rb)
	// put in & get out
	result, err := communicator.Put_Json_http(uri, b, agent.AgentKey, agent.ApiKey, agent.Zone)
	if err != nil {
		logger.Debugf("Handshake url:%s, agent:%s, api:%s, zone:%s",
			uri, agent.AgentKey, agent.ApiKey, agent.Zone)
		logger.Errorf("Failed Handshake (%v)", err)
		return nil
	}

	logger.Debugf("%s", string(result))

	body, unmarshalError := JsonUnmarshal(result)
	if unmarshalError != nil {
		logger.Debugf("Handshake url:%s, agent:%s, api:%s, zone:%s",
			uri, agent.AgentKey, agent.ApiKey, agent.Zone)
		logger.Errorf("The content of payload passed after handshake is unknown (%v)", err)
		return nil
	}

	logger.Debugf("%v", body)
	agent.schedulerInterval = body.Me.CallCycle

	if len(body.Agent.Nodes) > 0 {
		for _, v := range body.Agent.Nodes {
			agent.Agents = append(agent.Agents, v)
		}
	}

	return &body.Agent.Primary
}
