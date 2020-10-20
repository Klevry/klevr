package agent

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/NexClipper/logger"
)

// get task command from git
func getCommand(agent *KlevrAgent) {
	// SSH_provbee := "ssh provbee-service "

	uri := agent.Manager + "/agents/commands/init"

	var loop = true

	for loop {
		time.Sleep(1 * time.Second)

		provcheck := exec.Command("sh", "-c", "ssh provbee-service busybee beestatus hello > /tmp/con")
		errcheck := provcheck.Run()
		if errcheck != nil {
			logger.Errorf("provbee-service is not running: %v", errcheck)
		}

		by := readFile("/tmp/con")
		str := strings.TrimRight(string(by), "\n")

		if strings.Compare(str, "hi") == 0 {
			loop = false

			result := communicator.Get_Json_http(uri, agent.AgentKey, agent.API_key, agent.Zone)

			var Body common.Body
			err := json.Unmarshal(result, &Body)

			if err != nil {
				logger.Error(err)
			}

			// coms := Body.Task[0].Command
			// com := strings.Split(coms, "\n")

			// filenum := len(com)

			// for i := 0; i < filenum-1; i++ {
			// num := strconv.Itoa(i)

			// var read string

			// err := json.Unmarshal(readFile(Commands+num), &read)
			// if err != nil {
			// 	logger.Error(err)
			// }

			// if com[i] == read {
			// 	logger.Debugf("same command")
			// } else {
			// 	logger.Debugf("%d-----%s", i, com[i])
			// 	writeFile(Commands+num, com[i])

			// execute := SSH_provbee + com[i]
			//execute := com[i]

			// exe := exec.Command("sh", "-c", execute)
			// errExe := exe.Run()
			// if errExe != nil {
			// 	logger.Error(errExe)
			// }

			// }
			// }

			if _, err := os.Stat("/tmp/grafana"); !os.IsNotExist(err) {
				data, err := ioutil.ReadFile("/tmp/grafana")
				if err != nil {
					logger.Errorf("/tmp/grafana is not exist: %v", err)
				}

				//logger.Debugf("%v", string(data))

				if string(data) != "" {

					da := strings.Split(string(data), "\n")

					logger.Debugf("%v", da[0])
					// primaryInit(Body, coms, "done", da[0], agent)

					agent.initialized = true
				}

			}
		}
	}
}

func primaryInit(bod common.Body, command string, status string, param string, agent *KlevrAgent) []byte {
	uri := agent.Manager + "/agents/zones/init"

	rb := &common.Body{}

	SendMe(rb)

	par := make(map[string]interface{})
	par["grafana"] = param

	// rb.Task = make([]common.Task, 1)
	// rb.Task[0].ID = bod.Task[0].ID
	// rb.Task[0].AgentKey = bod.Task[0].AgentKey
	// rb.Task[0].Command = command
	// rb.Task[0].Status = status
	// rb.Task[0].Params = par
	// rb.Task[0].Result = bod.Task[0].Result
	// rb.Task[0].Type = bod.Task[0].Type

	b, err := json.Marshal(rb)
	if err != nil {
		logger.Error(err)
	}

	logger.Debugf("request body : [%s]", b)

	result := communicator.Post_Json_http(uri, b, agent.AgentKey, agent.API_key, agent.Zone)
	return result
}
