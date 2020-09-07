package agent

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/NexClipper/logger"
	"io"
	"os"
	"strconv"
	"time"
)

var agent_id_file = "/tmp/klevr_agent.id"
var agent_id_string string

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

		writeFile(agent_id_file, key)

		agent_id_string = key
	}

	return agent_id_string
}

func Check_primary(prim string) string {
	var ami string

	if prim == Local_ip_add() {
		ami = "true"
	} else if prim != Local_ip_add() {
		ami = "false"
	}
	return ami
}

func printprimary(prim string) {
	if Check_primary(prim) == "true" {
		logger.Debugf("-----------I am Primary")
	} else {
		logger.Debugf("-----------I am Secondary")
	}
	logger.Debugf("Primary ip : %s, My ip : %s", prim, Local_ip_add)
}
