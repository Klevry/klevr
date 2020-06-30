package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"io/ioutil"
	"os"
)

var Status string
var Master_status = "./info/master_status.info"

func Check_master() string{
	body, err := ioutil.ReadFile(Master_status)
	if err != nil{
		Status = "Master uptime check file open error."
	}
	line_break := strings.Split(string(body), "\n")
	//master_time, _ := strconv.Atoi(line_break[0])
	master_time := string(line_break[0])
//	time_parsing, _ := strconv.Atoi(line_break[0])
//fmt.Printf("&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&&:%d", master_time)
	time_parsing, err := strconv.ParseInt(master_time, 10, 64)
	tm := time.Unix(time_parsing, 0)
	if time.Since(tm).Minutes() > 60{
		println("Master agent is not working...")
		os.Exit(1)
	}else{
		Status = "OK"
	}
	return Status
}


func main() {
	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		Check_master()
		w.Write([]byte(Status))
	})
	http.ListenAndServe(":18800", nil)
}
