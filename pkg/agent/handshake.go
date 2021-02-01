package agent

import (
	"fmt"
	"time"

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

	retryCnt := 0
RETRY:
	// put in & get out
	result, err := communicator.Put_Json_http(uri, b, agent.AgentKey, agent.ApiKey, agent.Zone)
	if err != nil {
		logger.Debug(fmt.Sprintf("Failed Handshake %v", err))
		if retryCnt < 3 {
			time.Sleep(time.Second * 1)
			retryCnt++
			goto RETRY
		}
		return nil
	}

	logger.Debug("%s", string(result))

	body := JsonUnmarshal(result)

	logger.Debugf("%v", body)
	agent.schedulerInterval = body.Me.CallCycle

	if len(body.Agent.Nodes) > 0 {
		for _, v := range body.Agent.Nodes {
			agent.Agents = append(agent.Agents, v)
		}
	}

	return &body.Agent.Primary
}
