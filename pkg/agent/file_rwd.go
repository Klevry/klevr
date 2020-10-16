package agent

import (
	"encoding/json"
	"github.com/NexClipper/logger"
	"io/ioutil"
	"os"
)

func writeFile(path string, data interface{}) {
	d, _ := json.MarshalIndent(data, "", "  ")
	err := ioutil.WriteFile(path, d, os.FileMode(0644))
	if err != nil {
		logger.Error(err)
	}
}

func readFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error(err)
	}

	return data
}

func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		logger.Error(err)
	}

}
