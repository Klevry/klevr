package agent

import (
	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/NexClipper/logger"
	"net"
	"strconv"
	"time"
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

	result := communicator.Get_Json_http(uri, agent.AgentKey, agent.ApiKey, agent.Zone)

	body := JsonUnmarshal(result)

	if body.Agent.Primary.IP == Local_ip_add(){
		agent.Primary = body.Agent.Primary
		agent.startScheduler()
	}

	logger.Debugf("%v", body.Agent.Primary)
}
