package agent

import (
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strings"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/jasonlvhit/gocron"
	"github.com/mackerelio/go-osstat/memory"
)

const defaultSchedulerInterval int = 5

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
	grpcPort             string
}

func New(apiKey, platform, zone, manager, networkInterfaceName string, httpTimeout int) *KlevrAgent {
	agentKey := checkAgentKey()

	instance := &KlevrAgent{
		AgentKey:             agentKey,
		ApiKey:               apiKey,
		Platform:             platform,
		Zone:                 zone,
		Manager:              manager,
		NetworkInterfaceName: networkInterfaceName,
		HttpTimeout:          httpTimeout,
		grpcPort:             "9350",
	}

	logger.Debug("platform: ", instance.Platform)
	logger.Debug("Local_ip_add:", localIPAddress(instance.NetworkInterfaceName))
	logger.Debug("Agent UniqID:", instance.AgentKey)

	return instance
}

func (agent *KlevrAgent) Run() {
	primary := agent.handShake()
	if primary == nil || primary.IP == "" {
		logger.Error("Failed Handshake: Invalid Primary")
		return
	}
	agent.Primary = *primary
	agent.startScheduler()

	http.HandleFunc("/loglevel", agent.logLevelHandler)
	if agent.NetworkInterfaceName == "" {
		err := http.ListenAndServe(":18800", nil)
		if err != nil {
			panic(err)
		}
	} else {
		address := localIPAddress(agent.NetworkInterfaceName)
		err := http.ListenAndServe(fmt.Sprintf("%s:18800", address), nil)
		if err != nil {
			panic(err)
		}
	}
}

func (agent *KlevrAgent) startScheduler() {
	s := gocron.NewScheduler()

	if agent.checkPrimary(agent.Primary.IP) {
		var interval int
		if interval = agent.schedulerInterval; interval <= 0 {
			interval = defaultSchedulerInterval
		}

		s.Every(5).Seconds().Do(polling, agent)

	} else {
		go agent.secondaryServer()
		s.Every(5).Seconds().Do(statusCheck, agent)
	}

	go func() {
		<-s.Start()
	}()
}

func (agent *KlevrAgent) logLevelHandler(w http.ResponseWriter, r *http.Request) {
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

// send agent info to manager
func (agent *KlevrAgent) setBodyMeInfo(body *common.Body) {
	body.Me.IP = localIPAddress(agent.NetworkInterfaceName)
	body.Me.Port = 18800
	body.Me.Version = "0.1.0"

	disk := diskUsage("/")

	memory, _ := memory.Get()

	body.Me.Resource.Core = runtime.NumCPU()
	body.Me.Resource.Memory = int(memory.Total / MB)
	body.Me.Resource.Disk = int(disk.All / MB)
}
