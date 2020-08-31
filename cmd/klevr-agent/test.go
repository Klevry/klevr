package main

import (
	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/jasonlvhit/gocron"
	"io/ioutil"
	"strings"
)

var test string
var s = gocron.NewScheduler()

func read(){
	data, err := ioutil.ReadFile("test.txt")
	if err != nil {
		logger.Error(err)
	}

	test = string(data)

	logger.Debugf(test)
}

func rem(){
	var hi = "hi"
	var str  = strings.TrimRight(test, "\n")

	if strings.Compare(str, hi) == 0 {
		logger.Debugf("same")
		s.Remove(read)
	} else {
		logger.Debugf("%s-------%s", str, hi)
	}
}

func main(){
	common.InitLogger(common.NewLoggerEnv())

	s.Every(5).Seconds().Do(read)
	s.Every(5).Seconds().Do(rem)



	<-s.Start()


}