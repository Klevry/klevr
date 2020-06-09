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
	"flag"
	_"regexp"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/ralfyang/klevr/communicator"
)


//var api_url string
//var api_url="http://192.168.2.100:8500"
var api_url="http://localhost"
//var api_provision_script = api_url+"/ui/dc1/kv/klevr/"
var agent_download = "https://bit.ly/go_inst"
var api_provision_script string
var Hostlist string
var master_info string
var http_body_buffer string
var service_port = "8080" 

var api_key_string = "TEMPORARY"


// For Post funcition
func Init_script_api(uri, data string){
	req_body := bytes.NewBufferString(data)
	req, err := http.Post(uri,"text/plain", req_body)
	if err != nil {
		// handle error
		panic(err)
	}
	defer req.Body.Close()
	req.Header.Add("nexcloud-auth-token",api_key_string)
	req.Header.Add("cache-control", "no-cache")
}
// Init_script_api(api_url+"/v1/kv/klevr/form", data)



func Put_http(uri, data string) {
//	data, err := os.Open("text.txt")
//	println(uri,":",data)
	url := api_url+uri
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(data)))
	if err != nil {
		log.Printf("HTTP Put Request error: ",err)
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Add("nexcloud-auth-token",api_key_string)
	req.Header.Add("cache-control", "no-cache")
    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()
}


func Get_provision_script() string{
	http_body_buffer = communicator.Get_http(api_url+"/v1/kv/klevr/form?raw=1", api_key_string )
	if len(string(http_body_buffer)) == 0 {
		/// Set Script for instruction
		uri := "/v1/kv/klevr/form"
		data := "bash -s 'echo \"hello world\"'" /// Temporary use
		Put_http(uri, data)
		/// Read again
		api_provision_script = communicator.Get_http(api_url+"/v1/kv/klevr/form?raw=1", api_key_string)
	}else {
		api_provision_script = communicator.Get_http(api_url+"/v1/kv/klevr/form?raw=1", api_key_string)
	}
	return api_provision_script
}


//println("Body :", api_provision_script)


func API_Server_info(w http.ResponseWriter, r *http.Request){
	// w.Write([]byte("<a href='https://bit.ly/startdocker' target='blank'>Download Klever agent</a>"))
	Get_provision_script()
	w.Write([]byte(api_url))
}
func LandingPage(w http.ResponseWriter, r *http.Request){
	// w.Write([]byte("<a href='https://bit.ly/startdocker' target='blank'>Download Klever agent</a>"))
	Get_provision_script()
	w.Write([]byte("curl -sL "+agent_download+" | "+api_provision_script))
}

func Set_param() string{
	//Parsing by Flag
	port := flag.String("port",service_port,"Set port number for Service")
	api_server := flag.String("apiserver",api_url,"Set API Server URI for comunication")
	flag.Parse()
	service_port = *port
	api_url = *api_server
	return service_port
}

func Get_master(user string) string{
	master_info = communicator.Get_http(api_url+"/v1/kv/klevr/"+user+"/masters?raw=1", api_key_string)
		if len(master_info) == 0{
			master_info = "Not yet"
		}
	return master_info
}

func Get_host(user string) string{
	dataJson := communicator.Get_http(api_url+"/v1/kv/klevr/"+user+"/hosts/?keys", api_key_string)
	var arr []string
	_ = json.Unmarshal([]byte(dataJson), &arr)
	Get_master(user)
		if master_info == "Not yet"{
			http_body_buffer = communicator.Get_http(api_url+"/v1/kv/"+arr[0]+"?raw=1", api_key_string)
			strr1 := strings.Split(http_body_buffer, "&")
			strr2 := strings.Split(strr1[1], "=")
			master_info = "master="+strr2[1]

			uri := "/v1/kv/klevr/"+user+"/masters?raw=1"
			Put_http(uri, master_info)
//			curl -sL -H 'nexcloud-auth-token:testfordev' Cache-Control: no-cache --request PUT -d'master=192.168.2.100' klevr_account_api+"/"+account_name+"/masters"
		}

	var quee = master_info+"\n"
	for i := 0; i < len(arr); i++ {
		get_data := arr[i]
//		println("/v1/kv/"+get_data+"?raw=1")  // Test output
		//Get_http("/v1/kv/"+get_data)
		http_body_buffer = communicator.Get_http(api_url+"/v1/kv/"+get_data+"?raw=1", api_key_string)
		quee = quee+http_body_buffer+"\n"
	}
	Hostlist = quee
	return Hostlist

}

func main() {
	Set_param()
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", LandingPage)
	r.HandleFunc("/apiserver", API_Server_info)
        r.HandleFunc("/user/{U}/hostsinfo", func(w http.ResponseWriter, r *http.Request) {
        /// Export result to web
                vars := mux.Vars(r)
                user := vars["U"]
                /// Test out to the Browser
                Get_host(user)
                fmt.Fprintf(w, "User: %s\n", user)
                fmt.Fprintf(w, "\n\nHost(s) info.: \n%s\n", Hostlist)
        })

	// Bind to a port and pass our router in
	println("Service port:",service_port)
	println("Target API Server:",api_url)
	log.Printf("Web-console operation error: ",http.ListenAndServe(":"+service_port, r))
}



