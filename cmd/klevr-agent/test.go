package main

import (
	"github.com/NexClipper/logger"
	"github.com/jasonlvhit/gocron"
	"sync"
)

var path = "./test.txt"
var once sync.Once


func fileExist(){
	once.Do(func() {
		logger.Debugf("once test")
	})
	logger.Debugf("JES")
}
func main(){
	s := gocron.NewScheduler()
	s.Every(5).Seconds().Do(fileExist)

	<- s.Start()
}