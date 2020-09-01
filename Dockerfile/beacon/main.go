package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var Status string
var Master_status = "/info/master_status.info"

func Check_master() string {
	body, err := ioutil.ReadFile(Master_status)
	if err != nil {
		Status = "Master uptime check file open error."
	}
	line_break := strings.Split(string(body), "\n")
	master_time := string(line_break[0])
	time_parsing, err := strconv.ParseInt(master_time, 10, 64)
	tm := time.Unix(time_parsing, 0)
	if time.Since(tm).Minutes() > 10 {
		println("Error: Master agent is not working...")
		os.Exit(1)
	} else {
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
