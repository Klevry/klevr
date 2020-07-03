package main

import (
	"fmt"
	"net/http"
	"log"
	_"bytes"
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
)


/// Default API URL set
var API_url="http://localhost:8500"

/// Klevr Download URL
var Agent_download = "https://bit.ly/go_inst"

/// For server set
var Service_port = "8080" 
var API_key_string = "TEMPORARY"

/// Global variable 
var API_provision_script string
var Hostlist string
var Buffer_result string
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


func LogRequest(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method: %s | url: %s", r.Method, r.URL.String())
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}


/// Get Hostlist
func Get_host(user string) string{
	var arr []string
	var arr_stop, fail_count, array_count int
	dataJson := communicator.Get_http(API_url+"/v1/kv/klevr/"+user+"/hosts/?keys", API_key_string)
	_ = json.Unmarshal([]byte(dataJson), &arr)
	Get_master(user)
		if Master_info == "Not yet"{
			for i := 0; i <len(arr); i++{
				endpoint := arr[i][strings.LastIndex(arr[i], "/")+1:]
				if endpoint == "health" {
					Http_body_buffer = communicator.Get_http(API_url+"/v1/kv/"+arr[i]+"?raw=1", API_key_string)  /// Endpoing value will be "~/health" part of API
					strr1 := strings.Split(Http_body_buffer, "&")
					strr2 := strings.Split(strr1[1], "=")
					Master_info = "master="+strr2[1]
					arr_stop = i
				}
				uri := "/v1/kv/klevr/"+user+"/masters"
				communicator.Put_http(API_url+uri, Master_info, API_key_string)
			}
		}else{
			array_count = 0
			fail_count = 0
			for i := 0; i <len(arr); i++{
				endpoint := arr[i][strings.LastIndex(arr[i], "/")+1:]
					if endpoint == "health" {
						Http_body_buffer = communicator.Get_http(API_url+"/v1/kv/"+arr[i]+"?raw=1", API_key_string)  /// Endpoing value will be "~/health" part of API
						strr1 := strings.Split(Http_body_buffer, "&")
						strr2 := strings.Split(strr1[1], "=")
						marr1 := strings.Split(Http_body_buffer, "&")
						for mm := 0 ; mm < len(marr1) ; mm++{
							marr2 := marr1[mm][strings.LastIndex(marr1[mm], "=")+1:]
							if marr2 == "failed"{
								fail_count = fail_count + 1
							}
						}
						Master_info = "master="+strr2[1]
	//					log.Println("Error: Target endpoint will be /health, but current address is: "+endpoint+" please check the range of array from API.")
						array_count = array_count + 1
					}
//				uri := "/v1/kv/klevr/"+user+"/masters"
//				communicator.Put_http(API_url+uri, Master_info, API_key_string)
			}
			println("fail_countfail_countfail_countfail_countfail_countfail_countfail_countfail_count:",fail_count)
			println("array_countarray_countarray_countarray_countarray_countarray_countarray_countarray_countarray_count:",array_count)
			if array_count == fail_count+1{
				println("Master is dead!!!!") // test output
			}else if array_count/2 <= fail_count+1 {
				println("Master has something wrong!!!") // test output
			}
		}

	var quee = Master_info+"\n"
	/// for From range 1 to end. Due to the overlap
//	for i := 1; i < len(arr); i++ {
//		get_data := arr[i]
		get_data := arr[arr_stop]
		Http_body_buffer = communicator.Get_http(API_url+"/v1/kv/"+get_data+"?raw=1", API_key_string)
//		println("HHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH",Http_body_buffer)
		quee = quee+Http_body_buffer+"\n"
//	}
	Hostlist = quee
	return Hostlist

}

func Get_info_master(user string){
	/// initial master info
	Get_host(user)
	Get_master(user)
}



func Put_master_ack(user, ack string){
	uri := "/v1/kv/klevr/"+user+"/master_ack"
	communicator.Put_http(API_url+uri, ack, API_key_string)
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
			endpoint := arr[i][strings.LastIndex(arr[i], "/")+1:]
			if endpoint == "health" {
				queue = communicator.Get_http(API_url+"/v1/kv/"+arr[i]+"?raw=1", API_key_string)  /// Endpoing value will be "~/health" part of API
				get_data := arr[i]

				/// Get value of each hosts
				target_key = API_url+"/v1/kv/"+get_data
				println("target_key=", target_key) ///////////  Test output
				/// Parsing the Key/value of host_info
				target_txt = strings.Split(string(queue), "&")
				time_arry = strings.Split(target_txt[0], "=")

				/// Parsing the Key/value for Unix Time
				time_string = string(time_arry[1])
				time_parsing, err := strconv.ParseInt(time_string, 10, 64)
				if err != nil {
					log.Println(err)
				}
				/// Duration check 
				tm := time.Unix(time_parsing, 0)
				if time.Since(tm).Hours() > 1 {
					/// Delete old host via API server
					Host_purge_result = Host_purge_result+"Overtime: "+get_data+"\n"
					communicator.Delete_http(API_url+"/v1/kv/"+get_data, API_key_string)
				}else{
					Host_purge_result = Host_purge_result+"It's ok: "+get_data+"\n"
				}
			}
		}
	return Host_purge_result
}


func Client_receiver(user, hostname, host_ip, host_type, host_alive, master_alive string)string{
	uri := "/v1/kv/klevr/"+user+"/hosts/"+hostname+"/health"
	data := "last_check="+host_alive+"&ip="+host_ip+"&clientType="+host_type+"&masterConnection="+master_alive
	communicator.Put_http(API_url+uri, data, API_key_string)
	Buffer_result = data
	return Buffer_result
}

func Put_hostinfo(user, hostname, body string)string{
	uri := "/v1/kv/klevr/"+user+"/hosts/"+hostname+"/hostinfo"
	data := body
	communicator.Put_http(API_url+uri, data, API_key_string)
	Buffer_result = data
	return Buffer_result
}


func main() {
	Set_param()
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", LandingPage)

	/// Master ack receiver
        r.HandleFunc("/user/{U}/ackmaster", func(w http.ResponseWriter, r *http.Request) {
                vars := mux.Vars(r)
                user := vars["U"]
		ack_time := fmt.Sprint(time.Now().Unix())
		Put_master_ack(user, ack_time)
        })

	/// Hostinfo receiver
        r.HandleFunc("/user/{U}/hostsinfo", func(w http.ResponseWriter, r *http.Request) {
                vars := mux.Vars(r)
                user := vars["U"]
                Get_host(user)
	        /// Export result to web
                fmt.Fprintf(w, "User: %s\n", user)
                fmt.Fprintf(w, "\n\nHost(s) info.: \n%s\n", Hostlist)
        })

	/// Master status receiver
        r.HandleFunc("/user/{U}/masterinfo", func(w http.ResponseWriter, r *http.Request) {
                vars := mux.Vars(r)
                user := vars["U"]
                Get_info_master(user)
	        /// Export result to web
                fmt.Fprintf(w, "%s", Master_info)
        })

	/// Check hostpool & purge
        r.HandleFunc("/user/{U}/hostsmgt", func(w http.ResponseWriter, r *http.Request) {
                vars := mux.Vars(r)
                user := vars["U"]
		Hostpool_mgt(user)
	        /// Export result to web
                fmt.Fprintf(w, "User: %s\n", user)
                fmt.Fprintf(w, "\nHostresult: \n%s\n", Host_purge_result)
        })

	/// Callback receiver
        r.HandleFunc("/user/{U}/job/{JOB}/ticket/{TICKET}/{MSG}", func(w http.ResponseWriter, r *http.Request) {
                vars := mux.Vars(r)
                user := vars["U"]
                job := vars["JOB"]
                ticket := vars["TICKET"]
                callback_msg := vars["MSG"]
	        /// Export result to web
                fmt.Fprintf(w, "User: %s\n", user)
                fmt.Fprintf(w, "Job: %s\n", job)
                fmt.Fprintf(w, "Ticket number: %s\n", ticket)
                fmt.Fprintf(w, "Callback message: %s\n", callback_msg)

        })

	/// Host alive time & info receiver
        r.HandleFunc("/user/{U}/hostname/{HH}/{II}/type/{TP}/{TTL}/{MLO}", func(w http.ResponseWriter, r *http.Request) {
		// ralf, c3349a6b4c40908ec07fa4667b661362b76fba7d, 192.168.2.100, baremetal, 1592385021, ok
                vars := mux.Vars(r)
                user := vars["U"]
                hostname := vars["HH"]
                host_ip := vars["II"]
                host_type := vars["TP"]
                host_alive := string(vars["TTL"])
                master_alive := vars["MLO"]
		Client_receiver(user, hostname, host_ip, host_type, host_alive, master_alive)
	        /// Export result to web
                fmt.Fprintf(w, "User: %s\n", user)
                fmt.Fprintf(w, "\nResult: \n%s\n", Buffer_result)
        })

	/// receive json data to KV store
	r.StrictSlash(true)
	r.Use(LogRequest)
	r.HandleFunc("/user/{U}/hostname/{HH}/hostinfo", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["U"]
		hostname := vars["HH"]
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		jsondata := fmt.Sprintln(body)
		Put_hostinfo(user, hostname, jsondata)
		fmt.Fprintf(w, "Push result: %s \n", body)
//		fmt.Fprintf(w, "body: %s\n", body)
	})

	// Bind to a port and pass our router in
	println("Service port:",Service_port)
	println("Target API Server:",API_url)
	log.Printf("Web-console operation error: ",http.ListenAndServe(":"+Service_port, r))
}

