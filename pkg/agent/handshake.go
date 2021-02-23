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
	httpHandler := communicator.Http{
		URL:        uri,
		AgentKey:   agent.AgentKey,
		APIKey:     agent.ApiKey,
		ZoneID:     agent.Zone,
		RetryCount: 3,
		Timeout:    agent.HttpTimeout,
	}
	result, err := httpHandler.PutJson(b)
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
