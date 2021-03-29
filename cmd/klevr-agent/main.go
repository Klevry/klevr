package main

import (
	"flag"
	"os"
	_ "regexp"

	"github.com/Klevry/klevr/pkg/agent"
	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
)

var AGENT_VERSION = "0.0.1"

func main() {
	// TimeZone UTC로 설정
	os.Setenv("TZ", "")

	loggerEnv := &common.LoggerEnv{
		Level:      "debug",
		LogPath:    "./log/klevr.log",
		MaxSize:    3,
		MaxBackups: 5,
		MaxAge:     10,
		Compress:   false,
	}
	common.InitLogger(loggerEnv)

	// Flag options
	// Sample: -apiKey=\"{apiKey}\" -platform=\"{platform}\" -manager=\"{managerUrl}\" -zoneId=\"{zoneId}\" -iface=\"{networkInterfaceName}\"  -timeout=\"{timeout}\"
	apikey := flag.String("apiKey", "", "API Key from Klevr service")
	platform := flag.String("platform", "", "[baremetal|aws] - Service Platform for Host build up")
	zone := flag.String("zoneId", "", "zone will be a [Dev/Stg/Prod]")
	klevrAddr := flag.String("manager", "", "Klevr webconsole(server) address (URL or IP, Optional: Port) for connect")
	iface := flag.String("iface", "", "The name of the network interface to use.(If the value is empty, the first searched name is used.)")
	requestTimeout := flag.Int("timeout", 0, "Timeout(seconds) for an http request to receive a response")

	flag.Parse() // Important for parsing

	// Check the null data from CLI
	if len(*apikey) == 0 {
		logger.Error("Please insert an API Key")
		os.Exit(0)
	}
	if len(*platform) == 0 {
		logger.Error("Please make sure the platform")
		os.Exit(0)
	}
	if len(*zone) == 0 {
		logger.Error("Please insert zoneId")
		os.Exit(0)
	}
	if len(*klevrAddr) == 0 {
		logger.Error("Please insert manager addr")
		os.Exit(0)
	}

	instance := agent.NewKlevrAgent()

	instance.ApiKey = *apikey
	instance.Platform = *platform
	instance.Zone = *zone
	instance.Manager = *klevrAddr
	instance.NetworkInterfaceName = *iface
	instance.HttpTimeout = *requestTimeout

	logger.Debug("platform: ", instance.Platform)
	logger.Debug("Local_ip_add:", agent.LocalIPAddress(instance.NetworkInterfaceName))
	logger.Debug("Agent UniqID:", instance.AgentKey)

	instance.Run()

}
