package agent

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/NexClipper/logger"
)

var agentIdFile = "/tmp/klevr_agent.id"
var agentIdString string

// generate agent key
func AgentKeyGen() (string, error) {

	nowTime := strconv.FormatInt(time.Now().UTC().Unix(), 10)

	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	key := hex.EncodeToString(uuid) + nowTime

	return key, nil
}

func CheckAgentKey() string {
	//if agent file exist
	if _, err := os.Stat(agentIdFile); !os.IsNotExist(err) {
		data := ReadFile(agentIdFile)

		if string(data) != "" {
			agentIdString = string(data)
		} else {
			logger.Error("There is no agent ID")
		}

	} else {
		key, err := AgentKeyGen()
		if err != nil {
			logger.Error(err)
		}

		err = ioutil.WriteFile(agentIdFile, []byte(key), os.FileMode(0644))
		if err != nil {
			logger.Error(err)
		}

		agentIdString = key
	}

	return agentIdString
}

func (agent *KlevrAgent) checkPrimary(prim string) bool {
	if prim == LocalIPAddress(agent.NetworkInterfaceName) {
		logger.Debug("I am Primary")

		return true
	} else {
		logger.Debug("I am Secondary")

		return false
	}
}
