package agent

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/jasonlvhit/gocron"
)

const defaultSchedulerInterval int = 5

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

type KlevrAgent struct {
	ApiKey               string
	Platform             string
	Zone                 string
	Manager              string
	NetworkInterfaceName string
	HttpTimeout          int
	AgentKey             string
	Version              string
	schedulerInterval    int
	connect              net.Listener
	scheduler            *gocron.Scheduler
	Primary              common.Primary
	Agents               []common.Agent
}

func NewKlevrAgent() *KlevrAgent {
	agentKey := CheckAgentKey()

	instance := &KlevrAgent{
		AgentKey: agentKey,
	}

	return instance
}

func (agent *KlevrAgent) Run() {
	primary := HandShake(agent)
	if primary == nil || primary.IP == "" {
		logger.Error("Failed Handshake: Invalid Primary")
		return
	}
	agent.Primary = *primary
	agent.startScheduler()

	http.HandleFunc("/loglevel", agent.LogLevelHandler)
	if agent.NetworkInterfaceName == "" {
		err := http.ListenAndServe(":18800", nil)
		if err != nil {
			panic(err)
		}
	} else {
		address := LocalIPAddress(agent.NetworkInterfaceName)
		err := http.ListenAndServe(fmt.Sprintf("%s:18800", address), nil)
		if err != nil {
			panic(err)
		}
	}
}

func (agent *KlevrAgent) startScheduler() {
	//var scheduleFunc interface{}

	s := gocron.NewScheduler()

	if agent.checkPrimary(agent.Primary.IP) {
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

func (agent *KlevrAgent) LogLevelHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			level := logger.GetLevel()

			var levelValue string
			switch int(level) {
			case 0:
				levelValue = "debug"
			case 1:
				levelValue = "info"
			case 2:
				levelValue = "warn"
			case 3:
				levelValue = "error"
			case 4:
				levelValue = "fatal"
			}

			w.WriteHeader(200)
			fmt.Fprintf(w, levelValue)
		}
	case "PUT":
		{
			targetLevel := make([]byte, r.ContentLength)
			r.Body.Read(targetLevel)
			var level logger.Level

			switch strings.ToLower(string(targetLevel)) {
			case "debug":
				level = 0
			case "info":
				level = 1
			case "warn", "warning":
				level = 2
			case "error":
				level = 3
			case "fatal":
				level = 4
			}

			logger.SetLevel(level)

			w.WriteHeader(200)
			fmt.Fprintf(w, "{\"updated\":ok}")
		}
	}
}
