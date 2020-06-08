package main

import (
	"net/http"
	"log"
	"bytes"
	_ "os/exec"
	_"os"
	_ "io"
	"io/ioutil"
	"strings"
	"flag"
	"github.com/gorilla/mux"
)


//var api_url string
//var api_url="http://192.168.2.100:8500"
var api_url="http://localhost"
//var api_provision_script = api_url+"/ui/dc1/kv/klevr/"
var agent_download = "https://bit.ly/startdocker"
var api_provision_script string
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



func Put_request(uri string) {
//	data, err := os.Open("text.txt")
//	println(uri,":",data)
	data := "bash -c 'hello world'"
	req, err := http.NewRequest("PUT", uri, strings.NewReader(string(data)))
	if err != nil {
		log.Fatal(err)
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
	uri := api_url+"/v1/kv/klevr/form?raw=1"
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Add("nexcloud-auth-token",api_key_string)
	req.Header.Add("cache-control", "no-cache")
	res, err := http.DefaultClient.Do(req)
	if err == nil {
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		api_provision_script = string(body)
		if len(string(body)) == 0 {
//			data := "Ralf Test!"
			uri := api_url+"/v1/kv/klevr/form"
			Put_request(uri)
		}

	}
	return api_provision_script
}


//println("Body :", api_provision_script)



func LandingPage(w http.ResponseWriter, r *http.Request){
	// w.Write([]byte("<a href='https://bit.ly/startdocker' target='blank'>Download Klever agent</a>"))
	Get_provision_script()
	w.Write([]byte("curl -sL "+agent_download+" | "+api_provision_script))
}


func set_param() string{
	//Parsing by Flag
	port := flag.String("port",service_port,"Set port number for Service")
	api_server := flag.String("apiserver",api_url,"Set API Server URI for comunication")
	flag.Parse()
	service_port = *port
	api_url = *api_server
	return service_port
}



func main() {
	set_param()
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", LandingPage)
	// Bind to a port and pass our router in
	println("Service port:",service_port)
	println("Target API Server:",api_url)
	log.Fatal(http.ListenAndServe(":"+service_port, r))
}






