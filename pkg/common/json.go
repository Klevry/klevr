package common

import (
	"encoding/json"

	"github.com/NexClipper/logger"
)

func JsonMarshal(a interface{}) []byte {
	b, err := json.Marshal(a)
	if err != nil {
		logger.Debugf("%v", string(b))
		logger.Error(err)
	}

	return b
}

func JsonUnmarshal(a []byte) (*Body, error) {
	var body Body

	err := json.Unmarshal(a, &body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}
