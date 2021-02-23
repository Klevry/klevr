package agent

import (
	"net"
	"strconv"
	"time"

	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/NexClipper/logger"
)

func StatusCheck(agent *KlevrAgent) {
	_, err := net.DialTimeout("tcp", agent.Primary.IP+":"+strconv.Itoa(agent.Primary.Port), 3*time.Second)
	if err != nil {
		logger.Errorf("%v", err)
		PrimaryStatusReport(agent)
	}
}

/*
in: body.me, body.agent.primary
out: body.me, body.agent.primary
*/
func PrimaryStatusReport(agent *KlevrAgent) {
	uri := agent.Manager + "/agents/reports/" + agent.AgentKey

	httpHandler := communicator.Http{
		URL:        uri,
		AgentKey:   agent.AgentKey,
		APIKey:     agent.ApiKey,
		ZoneID:     agent.Zone,
		RetryCount: 1,
		Timeout:    agent.HttpTimeout,
	}
	result, err := httpHandler.GetJson()
	if err != nil {
		logger.Debugf("PrimaryStatusReport url:%s, agent:%s, api:%s, zone:%s", uri, agent.AgentKey, agent.ApiKey, agent.Zone)
		logger.Errorf("Failed PrimaryStatusReport (%v)", err)
		return
	}

	body, err := JsonUnmarshal(result)
	if err != nil {
		logger.Debugf("PrimaryStatusReport url:%s, agent:%s, api:%s, zone:%s", uri, agent.AgentKey, agent.ApiKey, agent.Zone)
		logger.Errorf("The content of payload passed after primarystatusreport is unknown (%v)", err)
		return
	}

	if body.Agent.Primary.IP == LocalIPAddress(agent.NetworkInterfaceName) {
		agent.Primary = body.Agent.Primary
		agent.startScheduler()
	}

	logger.Debugf("%v", body.Agent.Primary)
}
