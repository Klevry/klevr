package main

import (
	"flag"
	"github.com/Klevry/klevr/pkg/agent"
	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"os"
	_ "regexp"
)

var AGENT_VERSION = "0.0.1"

func main() {
	common.InitLogger(common.NewLoggerEnv())

	// Flag options
	// Sample: -apiKey=\"{apiKey}\" -platform={platform} -manager=\"{managerUrl}\" -zoneId={zoneId}
	apikey := flag.String("apiKey", "", "API Key from Klevr service")
	platform := flag.String("platform", "", "[baremetal|aws] - Service Platform for Host build up")
	zone := flag.String("zoneId", "", "zone will be a [Dev/Stg/Prod]")
	klevr_addr := flag.String("manager", "", "Klevr webconsole(server) address (URL or IP, Optional: Port) for connect")

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
	if len(*klevr_addr) == 0 {
		logger.Error("Please insert manager addr")
		os.Exit(0)
	}

	instance := agent.NewKlevrAgent()

	instance.API_key = *apikey
	instance.Platform = *platform
	instance.Zone = *zone
	instance.Manager = *klevr_addr

	logger.Debug("platform: ", instance.Platform)
	logger.Debug("Local_ip_add:", agent.Local_ip_add())
	logger.Debug("Agent UniqID:", instance.AgentKey)

	instance.Run()

}
