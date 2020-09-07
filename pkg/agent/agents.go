package agent

import (
	"github.com/jasonlvhit/gocron"
	"net/http"
)

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

type KlevrAgent struct {
	API_key  string
	Platform string
	Zone     string
	Manager  string
	AgentKey string
	Version  string
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

	if Check_primary(prim) == "true" {
		primScheduler := gocron.NewScheduler()
		primScheduler.Every(5).Seconds().Do(printprimary, prim)
		primScheduler.Every(5).Seconds().Do(getCommand, agent, *primScheduler)

		go func() {
			<-primScheduler.Start()
		}()
	} else {
		s := gocron.NewScheduler()
		s.Every(5).Seconds().Do(printprimary)

		go func() {
			<-s.Start()
		}()
	}
}
