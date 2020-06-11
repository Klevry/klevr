package main

import (
	"fmt"
	"net/http"
	"log"
	"bytes"
	_ "os/exec"
	_"os"
	_ "io"
	_"io/ioutil"
	"strings"
	"strconv"
	"time"
	"flag"
	_"regexp"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ralfyang/klevr/communicator"
	_"../communicator"
)


/// Default API URL set
var API_url="http://localhost"

/// Klevr Download URL
var Agent_download = "https://bit.ly/go_inst"

/// For server set
var Service_port = "8080" 
var API_key_string = "TEMPORARY"

/// Global variable 
var API_provision_script string
var Hostlist string
var Host_purge_result string
var Master_info string
var Http_body_buffer string


/// For custom scripts when the agent download & install
func Get_provision_script() string{
	Http_body_buffer = communicator.Get_http(API_url+"/v1/kv/klevr/form?raw=1", API_key_string )
	if len(string(Http_body_buffer)) == 0 {
		/// Set Script for instruction
		uri := "/v1/kv/klevr/form"
		data := "bash -s 'echo \"hello world\"'" /// Temporary use
		communicator.Put_http(API_url+uri, data, API_key_string)
		/// Read again
		API_provision_script = communicator.Get_http(API_url+"/v1/kv/klevr/form?raw=1", API_key_string)
	}else {
		API_provision_script = communicator.Get_http(API_url+"/v1/kv/klevr/form?raw=1", API_key_string)
	}
	return API_provision_script
}


/// Get API server information
func API_Server_info(w http.ResponseWriter, r *http.Request){
	// w.Write([]byte("<a href='https://bit.ly/startdocker' target='blank'>Download Klever agent</a>"))
	Get_provision_script()
	w.Write([]byte(API_url))
}


/// Get API key to agent
func API_key(w http.ResponseWriter, r *http.Request){
	// w.Write([]byte("<a href='https://bit.ly/startdocker' target='blank'>Download Klever agent</a>"))
	Get_provision_script()
	w.Write([]byte(API_key_string))
}


/// Default Landing page for http
func LandingPage(w http.ResponseWriter, r *http.Request){
	// w.Write([]byte("<a href='https://bit.ly/startdocker' target='blank'>Download Klever agent</a>"))
	Get_provision_script()
	w.Write([]byte("curl -sL "+Agent_download+" | "+API_provision_script))
}

/// Get Config variable when the webconsole start
func Set_param() string{
	//Parsing by Flag
	port := flag.String("port",Service_port,"Set port number for Service")
	api_server := flag.String("apiserver",API_url,"Set API Server URI for comunication")
	flag.Parse()
	Service_port = *port
	API_url = *api_server
	return Service_port
}

/// Get Master server infomation for slave agent control
func Get_master(user string) string{
	Master_info = communicator.Get_http(API_url+"/v1/kv/klevr/"+user+"/masters?raw=1", API_key_string)
		if len(Master_info) == 0{
			Master_info = "Not yet"
		}
	return Master_info
}

/// Get Hostlist
func Get_host(user string) string{
	dataJson := communicator.Get_http(API_url+"/v1/kv/klevr/"+user+"/hosts/?keys", API_key_string)
	var arr []string
	_ = json.Unmarshal([]byte(dataJson), &arr)
	Get_master(user)
		if Master_info == "Not yet"{
			Http_body_buffer = communicator.Get_http(API_url+"/v1/kv/"+arr[0]+"?raw=1", API_key_string)
			strr1 := strings.Split(Http_body_buffer, "&")
			strr2 := strings.Split(strr1[1], "=")
			Master_info = "master="+strr2[1]

			uri := "/v1/kv/klevr/"+user+"/masters?raw=1"
			communicator.Put_http(API_url+uri, Master_info, API_key_string)
		}

	var quee = Master_info+"\n"
	for i := 0; i < len(arr); i++ {
		get_data := arr[i]
		Http_body_buffer = communicator.Get_http(API_url+"/v1/kv/"+get_data+"?raw=1", API_key_string)
		quee = quee+Http_body_buffer+"\n"
	}
	Hostlist = quee
	return Hostlist

}

/// Old hostlist purge
func Hostpool_mgt(user string) string{
	/// Define variables
	var arr  []string
	var queue, target_key string
	Host_purge_result = "\n"

	/// Get Hostlist with Keys
	dataJson := communicator.Get_http(API_url+"/v1/kv/klevr/"+user+"/hosts/?keys", API_key_string)
	_ = json.Unmarshal([]byte(dataJson), &arr)
		for i := 0; i < len(arr); i++ {
			var target_txt, time_arry []string
			var time_string string

			get_data := arr[i]
			queue = communicator.Get_http(API_url+"/v1/kv/"+get_data+"?raw=1", API_key_string)   // klevr/ralf/hosts/0e25c6b9269944a543be0c82fb2fc8ce67e5b2c6/health

			/// Get value of each hosts
			target_key = API_url+"/v1/kv/"+get_data

			/// Parsing the Key/value of host_info
			target_txt = strings.Split(string(queue), "&")
			time_arry = strings.Split(target_txt[0], "=")

			/// Parsing the Key/value for Unix Time
			time_string = string(time_arry[1])
			time_parsing, err := strconv.ParseInt(time_string, 10, 64)
				if err != nil {
					panic(err)
				}
			/// Duration check 
			tm := time.Unix(time_parsing, 0)
				if time.Since(tm).Hours() > 24 {
					/// Delete old host via API server
					Host_purge_result = Host_purge_result+"Overtime: "+get_data+"\n"
					communicator.Delete_http(API_url+"/v1/kv/"+get_data, API_key_string)
				}else{
					Host_purge_result = Host_purge_result+"It's ok: "+get_data+"\n"
				}
		}
	return Host_purge_result
}


func main() {
	Set_param()
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", LandingPage)
	r.HandleFunc("/apiserver", API_Server_info)
	r.HandleFunc("/apikey", API_key)
        r.HandleFunc("/user/{U}/hostsinfo", func(w http.ResponseWriter, r *http.Request) {
                vars := mux.Vars(r)
                user := vars["U"]
                Get_host(user)
	        /// Export result to web
                fmt.Fprintf(w, "User: %s\n", user)
                fmt.Fprintf(w, "\n\nHost(s) info.: \n%s\n", Hostlist)
        })
        r.HandleFunc("/user/{U}/hostmgt", func(w http.ResponseWriter, r *http.Request) {
                vars := mux.Vars(r)
                user := vars["U"]
		Hostpool_mgt(user)
	        /// Export result to web
                fmt.Fprintf(w, "User: %s\n", user)
                fmt.Fprintf(w, "\nHostresult: \n%s\n", Host_purge_result)
        })

	// Bind to a port and pass our router in
	println("Service port:",Service_port)
	println("Target API Server:",API_url)
	log.Printf("Web-console operation error: ",http.ListenAndServe(":"+Service_port, r))
}

