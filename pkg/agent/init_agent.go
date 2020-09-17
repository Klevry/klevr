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

var agent_id_file = "/tmp/klevr_agent.id"
var agent_id_string string

// generate agent key
func AgentKeyGen() (string, error) {

	now_time := strconv.FormatInt(time.Now().UTC().Unix(), 10)

	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	key := hex.EncodeToString(uuid) + now_time

	return key, nil
	//return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func CheckAgentKey() string {
	//if agent file exist
	if _, err := os.Stat(agent_id_file); !os.IsNotExist(err) {
		data := readFile(agent_id_file)

		if string(data) != "" {
			agent_id_string = string(data)
		} else {
			logger.Error("There is no agent ID")
		}

	} else {
		key, err := AgentKeyGen()
		if err != nil {
			logger.Error(err)
		}

		// writeFile(agent_id_file, key)

		err = ioutil.WriteFile(agent_id_file, []byte(key), os.FileMode(0644))
		if err != nil {
			logger.Error(err)
		}

		agent_id_string = key
	}

	return agent_id_string
}

func Check_primary(prim string) bool {
	if prim == Local_ip_add() {
		logger.Debug("I am Primary")

		return true
	} else {
		logger.Debug("I am Secondary")

		return false
	}
}
