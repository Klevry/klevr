package agent

import (
	"net/http"

	"github.com/jasonlvhit/gocron"
)

const defaultSchedulerInterval int = 5

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

type KlevrAgent struct {
	API_key           string
	Platform          string
	Zone              string
	Manager           string
	AgentKey          string
	Version           string
	schedulerInterval int
	initialized       bool
	scheduler         *gocron.Scheduler
}

func NewKlevrAgent() *KlevrAgent {
	agentKey := CheckAgentKey()

	instance := &KlevrAgent{
		AgentKey: agentKey,
	}

	return instance
}

func (agent *KlevrAgent) Run() {

	primary_IP := HandShake(agent)
	agent.startScheduler(primary_IP)

	http.ListenAndServe(":18800", nil)
}

func (agent *KlevrAgent) startScheduler(prim string) {
	var scheduleFunc interface{}

	if Check_primary(prim) {
		var interval int
		if interval = agent.schedulerInterval; interval <= 0 {
			interval = defaultSchedulerInterval
		}

		go getCommand(agent)

		scheduleFunc = agent.tempHealthCheck
	} else {
		scheduleFunc = PrimaryStatusReport
	}

	s := gocron.NewScheduler()
	s.Every(5).Seconds().Do(scheduleFunc)

	agent.scheduler = s

	go func() {
		<-s.Start()
	}()
}
