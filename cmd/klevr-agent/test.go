package main

import (
	"encoding/json"
	"github.com/NexClipper/logger"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

var Path = "./command"

func test(){
	var com []string

	com[0] = "mkdir -p teat"
	com[1] = "echo \"test\" >> test"

	filenum := len(com)

	for i:=0; i<filenum; i++{
		num := strconv.Itoa(i)

		if(com[i] == string(readFile1(Path + num))){
			logger.Debugf("same command")
		} else {
			logger.Debugf("%d-------%s", i, com[i])
			writeFile1(Path+num, com[i])
		}
	}

	for i:=0; i<filenum-1; i++{
		num := strconv.Itoa(i)
		command := readFile1(Path+num)

		execute := string(command)[1:len(string(command))-1]
		logger.Debugf(execute)
		exe := exec.Command("sh", "-c", execute)
		errExe := exe.Run()
		if errExe != nil{
			logger.Error(errExe)
		} else {
			exe.Wait()
			deleteFile1(Path+num)
		}
	}
}

func writeFile1(path string, data string) {
	d, _ := json.MarshalIndent(data, "", "  ")
	err := ioutil.WriteFile(path, d, os.FileMode(0644))
	if err != nil{
		logger.Error(err)
	}
}

func readFile1(path string) []byte{
	data, err := ioutil.ReadFile(path)
	if err != nil{
		logger.Error(err)
	}

	return data
}


func deleteFile1(path string){
	err := os.Remove(path)
	if err != nil {
		logger.Error(err)
	}
}


func main(){
	//s := gocron.NewScheduler()
	//s.Every(1).Seconds().Do(test)
	//
	//go func() {
	//	<-s.Start()
	//}()
	test()
}