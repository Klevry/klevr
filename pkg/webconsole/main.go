package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	_ "io"
	_ "io/ioutil"
	"log"
	"net/http"
	_ "os"
	_ "os/exec"
	_ "regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Klevry/klevr/pkg/communicator"
	"github.com/gorilla/mux"
)

/// Default API URL set
var Consul_url = "http://localhost:8500"

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
var Primary_info string
var Http_body_buffer string

/// For custom scripts when the agent download & install
func Get_provision_script() string {
	Http_body_buffer = communicator.Get_http(Consul_url+"/v1/kv/klevr/form?raw=1", API_key_string)
	if len(string(Http_body_buffer)) == 0 {
		/// Set Script for instruction
		uri := "/v1/kv/klevr/form"
		data := "bash -s 'echo \"hello world\"'" /// Temporary use
		communicator.Put_http(Consul_url+uri, data, API_key_string)
		/// Read again
		API_provision_script = communicator.Get_http(Consul_url+"/v1/kv/klevr/form?raw=1", API_key_string)
	} else {
		API_provision_script = communicator.Get_http(Consul_url+"/v1/kv/klevr/form?raw=1", API_key_string)
	}
	return API_provision_script
}

/// Default Landing page for http
func LandingPage(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("<a href='https://bit.ly/startdocker' target='blank'>Download Klever agent</a>"))
	Get_provision_script()
	w.Write([]byte("curl -sL " + Agent_download + " | " + API_provision_script))
}

/// Get Config variable when the webconsole start
func Set_param() string {
	//Parsing by Flag
	port := flag.String("port", Service_port, "Set port number for Service")
	api_server := flag.String("apiserver", Consul_url, "Set API Server URI for comunication")
	flag.Parse()
	Service_port = *port
	Consul_url = *api_server
	return Service_port
}

//company user zone platform
//%s/+group+"/users/"+user+"/zones/"+zone+/+group+"\/groups/"+group+"/users/"+user+"/zones/"+zone+"\/zones/"+zone+/g
/// Get Primary server information for secondary agent control
func Get_primary(group, user, zone, platform string) string {
	Primary_info = communicator.Get_http(Consul_url+"/v1/kv/klevr/groups/"+group+"/users/"+user+"/zones/"+zone+"/platforms/"+platform+"/primarys?raw=1", API_key_string)
	if len(Primary_info) == 0 {
		Primary_info = "Not yet"
	}
	return Primary_info
}

func LogRequest(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method: %s | url: %s", r.Method, r.URL.String())
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

/// Get Hostlist
func Get_host(group, user, zone, platform, priyes string) string {
	var arr []string
	var quee string
	var arr_stop, fail_count, array_count int
	dataJson := communicator.Get_http(Consul_url+"/v1/kv/klevr/groups/"+group+"/users/"+user+"/zones/"+zone+"/platforms/"+platform+"/hosts/?keys", API_key_string)
	_ = json.Unmarshal([]byte(dataJson), &arr)
	Get_primary(group, user, zone, platform)
	if Primary_info == "Not yet" {
		for i := 0; i < len(arr); i++ {
			endpoint := arr[i][strings.LastIndex(arr[i], "/")+1:]
			if endpoint == "health" {
				Http_body_buffer = communicator.Get_http(Consul_url+"/v1/kv/"+arr[i]+"?raw=1", API_key_string) /// Endpoing value will be "~/health" part of API
				strr1 := strings.Split(Http_body_buffer, "&")
				strr2 := strings.Split(strr1[1], "=")
				Primary_info = "primary=" + strr2[1]
				arr_stop = i
			}
			uri := "/v1/kv/klevr/groups/" + group + "/users/" + user + "/zones/" + zone + "/platforms/" + platform + "/primarys"
			communicator.Put_http(Consul_url+uri, Primary_info, API_key_string)
			if priyes == "yes" {
				quee = quee + Primary_info
			} else {
				quee = quee
			}
			get_data := arr[arr_stop]
			Http_body_buffer = communicator.Get_http(Consul_url+"/v1/kv/"+get_data+"?raw=1", API_key_string)
			println("HHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH", Http_body_buffer)
			quee = quee + Http_body_buffer + "\n"
		}
	} else {
		array_count = 0
		fail_count = 0
		for i := 0; i < len(arr); i++ {
			endpoint := arr[i][strings.LastIndex(arr[i], "/")+1:]
			if endpoint == "health" {
				Http_body_buffer = communicator.Get_http(Consul_url+"/v1/kv/"+arr[i]+"?raw=1", API_key_string) /// Endpoing value will be "~/health" part of API
				strr1 := strings.Split(Http_body_buffer, "&")
				strr2 := strings.Split(strr1[1], "=")
				marr1 := strings.Split(Http_body_buffer, "&")
				// Failed counting with host listing
				for mm := 0; mm < len(marr1); mm++ {
					marr2 := marr1[mm][strings.LastIndex(marr1[mm], "=")+1:]
					if marr2 == "failed" {
						fail_count = fail_count + 1
					}
				}

				Primary_info = "primary=" + strr2[1]
				//					log.Println("Error: Target endpoint will be /health, but current address is: "+endpoint+" please check the range of array from API.")
				if priyes == "yes" {
					quee = quee + Primary_info
				} else {
					quee = quee
				}
				array_count = array_count + 1
				arr_stop = i
				get_data := arr[i]
				Http_body_buffer = communicator.Get_http(Consul_url+"/v1/kv/"+get_data+"?raw=1", API_key_string)
				quee = quee + Http_body_buffer + "\n"
				//						println("aaaaaaaaa--fail_countfail_countfail_countfail_countfail:",quee)
			}
		}
		//			println("fail_countfail_countfail_countfail_countfail_countfail_countfail_countfail_count:",fail_count)
		//			println("array_countarray_countarray_countarray_countarray_countarray_countarray_countarray_countarray_count:",array_count)
		if array_count == fail_count+1 {
			println("Primary is dead!!!!") // test output
		} else if array_count/2 <= fail_count+1 {
			println("Primary has something wrong!!!") // test output
		}
	}

	Hostlist = quee
	return Hostlist

}

func Get_info_primary(group, user, zone, platform string) {
	/// initial primary info
	Get_host(group, user, zone, platform, "")
	Get_primary(group, user, zone, platform)
}

func Put_platform_init(platform, data string) {
	uri := "/v1/kv/klevr/systems/platform_types/" + platform
	communicator.Put_http(Consul_url+uri, data, API_key_string)
}

func Put_primary_ack(group, user, zone, platform, ack string) {
	uri := "/v1/kv/klevr/groups/" + group + "/users/" + user + "/zones/" + zone + "/platforms/" + platform + "/primary_ack"
	communicator.Put_http(Consul_url+uri, ack, API_key_string)
}

/// Old hostlist purge
func Hostpool_mgt(group, user, zone, platform string) string {
	/// Define variables
	var arr []string
	var queue, target_key string
	Host_purge_result = "\n"

	/// Get Hostlist with Keys
	dataJson := communicator.Get_http(Consul_url+"/v1/kv/klevr/groups/"+group+"/users/"+user+"/zones/"+zone+"/platforms/"+platform+"/hosts/?keys", API_key_string)
	_ = json.Unmarshal([]byte(dataJson), &arr)
	for i := 0; i < len(arr); i++ {
		var target_txt, time_arry []string
		var time_string string
		endpoint := arr[i][strings.LastIndex(arr[i], "/")+1:]
		if endpoint == "health" {
			queue = communicator.Get_http(Consul_url+"/v1/kv/"+arr[i]+"?raw=1", API_key_string) /// Endpoing value will be "~/health" part of API
			get_data := arr[i]

			/// Get value of each hosts
			target_key = Consul_url + "/v1/kv/" + get_data
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
				Host_purge_result = Host_purge_result + "Overtime: " + get_data + "\n"
				communicator.Delete_http(Consul_url+"/v1/kv/"+get_data, API_key_string)
			} else {
				Host_purge_result = Host_purge_result + "It's ok: " + get_data + "\n"
			}
		}
	}
	return Host_purge_result
}

func Client_receiver(group, user, zone, hostname, host_ip, platform, host_alive, primary_alive string) string {
	uri := "/v1/kv/klevr/groups/" + group + "/users/" + user + "/zones/" + zone + "/platforms/" + platform + "/hosts/" + hostname + "/health"
	data := "last_check=" + host_alive + "&ip=" + host_ip + "&clientType=" + platform + "&primaryConnection=" + primary_alive
	communicator.Put_http(Consul_url+uri, data, API_key_string)
	Buffer_result = data
	return Buffer_result
}

func Put_hostinfo(group, user, zone, platform, hostname, body string) string {
	uri := "/v1/kv/klevr/groups/" + group + "/users/" + user + "/zones/" + zone + "/platforms/" + platform + "/hosts/" + hostname + "/hostinfo"
	data := body
	communicator.Put_http(Consul_url+uri, data, API_key_string)
	Buffer_result = data
	return Buffer_result
}

func main() {
	Set_param()
	r := mux.NewRouter()
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", LandingPage)

	/// Primary ack receiver
	r.HandleFunc("/group/{G}/user/{U}/zone/{Z}/platform/{P}/ackprimary", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["U"]
		platform := vars["P"]
		group := vars["G"]
		zone := vars["Z"]
		ack_time := fmt.Sprint(time.Now().Unix())
		Put_primary_ack(group, user, zone, platform, ack_time)
		Get_host(group, user, zone, platform, "")
		/// Export result to web
		fmt.Fprintf(w, "get_timestamp: %s\n", ack_time)
		fmt.Fprintf(w, "%s\n", Hostlist)
	})

	/// Hostinfo receiver
	r.HandleFunc("/group/{G}/user/{U}/zone/{Z}/platform/{P}/hostsinfo", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["U"]
		platform := vars["P"]
		group := vars["G"]
		zone := vars["Z"]
		Get_host(group, user, zone, platform, "")
		/// Export result to web
		fmt.Fprintf(w, "User: %s\n", user)
		fmt.Fprintf(w, "\n\nHost(s) info.: \n%s\n", Hostlist)
	})

	/// Primary status receiver
	r.HandleFunc("/group/{G}/user/{U}/zone/{Z}/platform/{P}/primaryinfo", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["U"]
		platform := vars["P"]
		group := vars["G"]
		zone := vars["Z"]
		Get_info_primary(group, user, zone, platform)
		/// Export result to web
		fmt.Fprintf(w, "%s", Primary_info)
	})

	/// Check hostpool & purge
	r.HandleFunc("/group/{G}/user/{U}/zone/{Z}/platform/{P}/hostsmgt", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["U"]
		platform := vars["P"]
		group := vars["G"]
		zone := vars["Z"]
		Hostpool_mgt(group, user, zone, platform)
		/// Export result to web
		fmt.Fprintf(w, "User: %s\n", user)
		fmt.Fprintf(w, "\nHostresult: \n%s\n", Host_purge_result)
	})

	/// Callback receiver
	r.HandleFunc("/group/{G}/user/{U}/zone/{Z}/platform/{P}/job/{JOB}/ticket/{TICKET}/{MSG}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["U"]
		job := vars["JOB"]
		ticket := vars["TICKET"]
		callback_msg := vars["MSG"]
		platform := vars["P"]
		group := vars["G"]
		zone := vars["Z"]
		/// Export result to web
		fmt.Fprintf(w, "Group: %s\n", group)
		fmt.Fprintf(w, "User: %s\n", user)
		fmt.Fprintf(w, "Zone: %s\n", zone)
		fmt.Fprintf(w, "Platform: %s\n", platform)
		fmt.Fprintf(w, "Job: %s\n", job)
		fmt.Fprintf(w, "Ticket number: %s\n", ticket)
		fmt.Fprintf(w, "Callback message: %s\n", callback_msg)

	})

	/// Host alive time & info receiver
	r.HandleFunc("/group/{G}/user/{U}/zone/{Z}/platform/{P}/hostname/{HH}/{II}/{TTL}/{MLO}", func(w http.ResponseWriter, r *http.Request) {
		// ralf, c3349a6b4c40908ec07fa4667b661362b76fba7d, 192.168.2.100, baremetal, 1592385021, ok
		vars := mux.Vars(r)
		user := vars["U"]
		hostname := vars["HH"]
		host_ip := vars["II"]
		host_alive := string(vars["TTL"])
		primary_alive := vars["MLO"]
		platform := vars["P"]
		group := vars["G"]
		zone := vars["Z"]
		Client_receiver(group, user, zone, hostname, host_ip, platform, host_alive, primary_alive)
		/// Export result to web
		fmt.Fprintf(w, "User: %s\n", user)
		fmt.Fprintf(w, "\nResult: \n%s\n", Buffer_result)
	})

	// http://10.10.33.2:8000/group/klevr-a-team/user/ralf/zone/dev/platform/baremetal/hostname/2db8c6cf4329e44bda97a72f7f9127125d991ba2/192.168.15.50/1594257516/ok
	/// receive json data to KV store
	r.StrictSlash(true)
	r.Use(LogRequest)
	r.HandleFunc("/group/{G}/user/{U}/zone/{Z}/platform/{P}/hostname/{HH}/hostinfo", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["U"]
		hostname := vars["HH"]
		platform := vars["P"]
		group := vars["G"]
		zone := vars["Z"]
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		jsondata := fmt.Sprintln(body)
		Put_hostinfo(group, user, zone, platform, hostname, jsondata)
		fmt.Fprintf(w, "Push result: %s \n", body)
		//		fmt.Fprintf(w, "body: %s\n", body)
	})

	/// Platform setup init. for preinstaller
	r.HandleFunc("/systems/platform_types/{P}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		platform := vars["P"]
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		jsondata := fmt.Sprintln(body)
		println(platform, jsondata)
		//		data := fmt.Sprintln(body)
		//		Put_platform_init(platform, data)
		/// Export result to web
		fmt.Fprintf(w, "%s", Primary_info)
	})

	/// Agent alive status receiver
	/// http://192.168.2.100:8000/groups/klevr-a-team/users/ralf/zones/dev/platforms/baremetal/alive_hosts
	r.HandleFunc("/groups/{G}/users/{U}/zones/{Z}/platforms/{P}/aliveagent", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		user := vars["U"]
		platform := vars["P"]
		group := vars["G"]
		zone := vars["Z"]

		data_buffer := new(bytes.Buffer)
		data_buffer.ReadFrom(r.Body)
		data_Str := data_buffer.String()

		/// Get string from r.Body via bytes buffer
		//		data := base64.NewDecoder(base64.StdEncoding, r.Body)
		data, _ := base64.StdEncoding.DecodeString(data_Str)
		plan_data := fmt.Sprintf("%s", data)
		println("datadatadatadatadata:", plan_data)
		ure := "/v1/kv/klevr/groups/" + group + "/users/" + user + "/zones/" + zone + "/platforms/" + platform + "/aliveagent"
		communicator.Put_http(Consul_url+ure, plan_data, API_key_string)
		fmt.Fprintf(w, "%s", plan_data)
	})

	// Bind to a port and pass our router in
	println("Service port:", Service_port)
	println("Target API Server:", Consul_url)
	log.Println("Web-console operation error: ", http.ListenAndServe(":"+Service_port, r))
}
