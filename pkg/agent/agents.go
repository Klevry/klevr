package agent

import (
	"github.com/Klevry/klevr/pkg/common"
	"github.com/jasonlvhit/gocron"
	"net"
	"net/http"
)

const defaultSchedulerInterval int = 5

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

type KlevrAgent struct {
	ApiKey            string
	Platform          string
	Zone              string
	Manager           string
	AgentKey          string
	Version           string
	schedulerInterval int
	connect           net.Listener
	scheduler         *gocron.Scheduler
	Primary           common.Primary
	Agents            []common.Agent
}

func NewKlevrAgent() *KlevrAgent {
	agentKey := CheckAgentKey()

	instance := &KlevrAgent{
		AgentKey: agentKey,
	}

	return instance
}

func (agent *KlevrAgent) Run() {
	agent.Primary = HandShake(agent)
	agent.startScheduler()

	http.ListenAndServe(":18800", nil)
}

func (agent *KlevrAgent) startScheduler() {
	//var scheduleFunc interface{}

	s := gocron.NewScheduler()

	if Check_primary(agent.Primary.IP) {
		var interval int
		if interval = agent.schedulerInterval; interval <= 0 {
			interval = defaultSchedulerInterval
		}

		s.Every(5).Seconds().Do(Polling, agent)

	} else {
		go agent.SecondaryServer()
		s.Every(5).Seconds().Do(StatusCheck, agent)
	}

	go func() {
		<-s.Start()
	}()
}
