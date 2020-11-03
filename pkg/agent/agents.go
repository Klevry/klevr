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
	PrimaryIP         string
	SecondaryIP       []Secondary
}

type Secondary struct {
	IP string
}

func NewKlevrAgent() *KlevrAgent {
	agentKey := CheckAgentKey()

	instance := &KlevrAgent{
		AgentKey: agentKey,
	}

	return instance
}

func (agent *KlevrAgent) Run() {
	
	logger.Debugf("new agent")

	agent.PrimaryIP = HandShake(agent)
	agent.startScheduler()

	http.ListenAndServe(":18800", nil)
}

func (agent *KlevrAgent) startScheduler() {
	//var scheduleFunc interface{}

	s := gocron.NewScheduler()

	if Check_primary(agent.PrimaryIP) {
		var interval int
		if interval = agent.schedulerInterval; interval <= 0 {
			interval = defaultSchedulerInterval
		}

		//go getCommand(agent)

		s.Every(5).Seconds().Do(Polling, agent)

		//scheduleFunc = agent.tempHealthCheck
		//scheduleFunc = Polling
	} else {
		//scheduleFunc = PrimaryStatusReport
	}

	//s.Every(5).Seconds().Do(scheduleFunc)

	agent.scheduler = s

	go func() {
		<-s.Start()
	}()
}
