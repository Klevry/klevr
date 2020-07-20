package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/NexClipper/logger"
	"github.com/gin-gonic/gin"
	"github.com/ralfyang/pkg/klevr/manager/model"
)

type apiDef struct {
	method   string
	uri      string
	function func(*gin.Context)
}

// var arr = [...]apiDef{
// 	apiDef{"any", "group/:G/user/:U/zone/:Z/platform/:P/ackprimary", test}
// }
// var apiMap = map[int]apiDef{
// 	1: &apiDef{method: "any", uri: "group/:G/user/:U/zone/:Z/platform/:P/ackprimary", function: Test},
// }

var apiSlice []apiDef

func (api *API) initAPI() {
	logger.Debug("API InitLegacy - init URI")
	apiSlice = append(apiSlice, apiDef{"any", "/group/:G/user/:U/zone/:Z/platform/:P/ackprimary", api.ackprimary})
	apiSlice = append(apiSlice, apiDef{"any", "/group/:G/user/:U/zone/:Z/platform/:P/hostsinfo", api.hostsinfo})
	apiSlice = append(apiSlice, apiDef{"any", "/group/:G/user/:U/zone/:Z/platform/:P/primaryinfo", api.primaryinfo})
	apiSlice = append(apiSlice, apiDef{"any", "/group/:G/user/:U/zone/:Z/platform/:P/hostsmgt", api.hostsmgt})
	apiSlice = append(apiSlice, apiDef{"any", "/group/:G/user/:U/zone/:Z/platform/:P/job/:JOB/ticket/:TICKET/:MSG", api.callback})
	apiSlice = append(apiSlice, apiDef{"any", "/group/:G/user/:U/zone/:Z/platform/:P/hostname/:H/hostinfo", api.hostinfo})
	// apiSlice = append(apiSlice, apiDef{"any", "/group/:G/user/:U/zone/:Z/platform/:P/hostname/:H/:I/:TTL/:MLO", api.alivetime})
	apiSlice = append(apiSlice, apiDef{"any", "/systems/platform_types/:P", api.initAgent})
	apiSlice = append(apiSlice, apiDef{"any", "/groups/:G/users/:U/zones/:Z/platforms/:P/aliveagent", api.statusReciever})
}

// InitLegacy initialize legacy API
func (api *API) InitLegacy(legacy *gin.RouterGroup) {
	logger.Debug("API InitLegacy")

	api.initAPI()

	for _, def := range apiSlice {
		switch def.method {
		case "any":
			legacy.Any(def.uri, def.function)
		case "get":
			legacy.GET(def.uri, def.function)
		case "post":
			legacy.POST(def.uri, def.function)
		case "put":
			legacy.PUT(def.uri, def.function)
		case "delete":
			legacy.DELETE(def.uri, def.function)
		case "patch":
			legacy.PATCH(def.uri, def.function)
		}
	}

	// legacy.Any("group/:G/user/:U/zone/:Z/platform/:P/ackprimary", func(c *gin.Context) {
	// 	// vars := mux.Vars(r)
	// 	// user := vars["U"]
	// 	// platform := vars["P"]
	// 	// group := vars["G"]
	// 	// zone := vars["Z"]
	// 	// ackTime := fmt.Sprint(time.Now().Unix())
	// 	// Put_primary_ack(group, user, zone, platform, ackTime)
	// 	// Get_host(group, user, zone, platform, "")
	// 	// /// Export result to web
	// 	// fmt.Fprintf(w, "get_timestamp: %s\n", ackTime)
	// 	// // fmt.Fprintf(w, "%s\n", Hostlist)
	// 	// w.Write([]byte("test"))

	// })
}

func (api *API) ackprimary(c *gin.Context) {
	group, _ := strconv.ParseInt(c.Param("G"), 10, 64)
	user, _ := strconv.ParseInt(c.Param("U"), 10, 64)
	zone := c.Param("Z")
	platform := c.Param("P")

	api.PutPrimaryAck(group, user, zone, platform, fmt.Sprint(time.Now().Unix()))
	// GetHost()
}

func (api *API) hostsinfo(c *gin.Context) {
	c.String(200, "test2")
}

func (api *API) primaryinfo(c *gin.Context) {

}

func (api *API) hostsmgt(c *gin.Context) {

}

func (api *API) callback(c *gin.Context) {

}

func (api *API) alivetime(c *gin.Context) {

}

func (api *API) hostinfo(c *gin.Context) {

}

func (api *API) initAgent(c *gin.Context) {

}

func (api *API) statusReciever(c *gin.Context) {

}

// GetProvisionScript For custom scripts when the agent download & install
func GetProvisionScript() string {
	// Http_body_buffer := communicator.Get_http(ConsulURL+"/v1/kv/klevr/form?raw=1", API_key_string)
	// if len(string(Http_body_buffer)) == 0 {
	// 	/// Set Script for instruction
	// 	uri := "/v1/kv/klevr/form"
	// 	data := "bash -s 'echo \"hello world\"'" /// Temporary use
	// 	communicator.Put_http(ConsulURL+uri, data, API_key_string)
	// 	/// Read again
	// 	API_provision_script = communicator.Get_http(ConsulURL+"/v1/kv/klevr/form?raw=1", API_key_string)
	// } else {
	// 	API_provision_script = communicator.Get_http(ConsulURL+"/v1/kv/klevr/form?raw=1", API_key_string)
	// }
	// return API_provision_script
	return ""
}

// LandingPage Default Landing page for http
func LandingPage(w http.ResponseWriter, r *http.Request) {
	// // w.Write([]byte("<a href='https://bit.ly/startdocker' target='blank'>Download Klever agent</a>"))
	// GetProvisionScript()
	// w.Write([]byte("curl -sL " + AgentDownload + " | " + API_provision_script))
}

// SetParam Get Config variable when the webconsole start
func SetParam() string {
	// //Parsing by Flag
	// port := flag.String("port", Service_port, "Set port number for Service")
	// api_server := flag.String("apiserver", ConsulURL, "Set API Server URI for comunication")
	// flag.Parse()
	// Service_port = *port
	// ConsulURL = *api_server
	// return Service_port
	return ""
}

// GetPrimary company user zone platform
//%s/+group+"/users/"+user+"/zones/"+zone+/+group+"\/groups/"+group+"/users/"+user+"/zones/"+zone+"\/zones/"+zone+/g
/// Get Primary server infomation for secondary agent control
func GetPrimary(group, user, zone, platform string) string {
	// Primary_info = communicator.Get_http(ConsulURL+"/v1/kv/klevr/groups/"+group+"/users/"+user+"/zones/"+zone+"/platforms/"+platform+"/primarys?raw=1", API_key_string)
	// if len(Primary_info) == 0 {
	// 	Primary_info = "Not yet"
	// }
	// return Primary_info

	return ""
}

// LogRequest ..
func LogRequest(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("method: %s | url: %s", r.Method, r.URL.String())
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// GetHost Get Hostlist
func GetHost(group, user, zone, platform, priyes string) string {
	// var arr []string
	// var quee string
	// var arr_stop, fail_count, array_count int
	// dataJson := communicator.Get_http(ConsulURL+"/v1/kv/klevr/groups/"+group+"/users/"+user+"/zones/"+zone+"/platforms/"+platform+"/hosts/?keys", API_key_string)
	// _ = json.Unmarshal([]byte(dataJson), &arr)
	// Get_primary(group, user, zone, platform)
	// if Primary_info == "Not yet" {
	// 	for i := 0; i < len(arr); i++ {
	// 		endpoint := arr[i][strings.LastIndex(arr[i], "/")+1:]
	// 		if endpoint == "health" {
	// 			Http_body_buffer = communicator.Get_http(ConsulURL+"/v1/kv/"+arr[i]+"?raw=1", API_key_string) /// Endpoing value will be "~/health" part of API
	// 			strr1 := strings.Split(Http_body_buffer, "&")
	// 			strr2 := strings.Split(strr1[1], "=")
	// 			Primary_info = "primary=" + strr2[1]
	// 			arr_stop = i
	// 		}
	// 		uri := "/v1/kv/klevr/groups/" + group + "/users/" + user + "/zones/" + zone + "/platforms/" + platform + "/primarys"
	// 		communicator.Put_http(ConsulURL+uri, Primary_info, API_key_string)
	// 		if priyes == "yes" {
	// 			quee = quee + Primary_info
	// 		} else {
	// 			quee = quee
	// 		}
	// 		get_data := arr[arr_stop]
	// 		Http_body_buffer = communicator.Get_http(ConsulURL+"/v1/kv/"+get_data+"?raw=1", API_key_string)
	// 		println("HHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH", Http_body_buffer)
	// 		quee = quee + Http_body_buffer + "\n"
	// 	}
	// } else {
	// 	array_count = 0
	// 	fail_count = 0
	// 	for i := 0; i < len(arr); i++ {
	// 		endpoint := arr[i][strings.LastIndex(arr[i], "/")+1:]
	// 		if endpoint == "health" {
	// 			Http_body_buffer = communicator.Get_http(ConsulURL+"/v1/kv/"+arr[i]+"?raw=1", API_key_string) /// Endpoing value will be "~/health" part of API
	// 			strr1 := strings.Split(Http_body_buffer, "&")
	// 			strr2 := strings.Split(strr1[1], "=")
	// 			marr1 := strings.Split(Http_body_buffer, "&")
	// 			// Failed counting with host listing
	// 			for mm := 0; mm < len(marr1); mm++ {
	// 				marr2 := marr1[mm][strings.LastIndex(marr1[mm], "=")+1:]
	// 				if marr2 == "failed" {
	// 					fail_count = fail_count + 1
	// 				}
	// 			}

	// 			Primary_info = "primary=" + strr2[1]
	// 			//					log.Println("Error: Target endpoint will be /health, but current address is: "+endpoint+" please check the range of array from API.")
	// 			if priyes == "yes" {
	// 				quee = quee + Primary_info
	// 			} else {
	// 				quee = quee
	// 			}
	// 			array_count = array_count + 1
	// 			arr_stop = i
	// 			get_data := arr[i]
	// 			Http_body_buffer = communicator.Get_http(ConsulURL+"/v1/kv/"+get_data+"?raw=1", API_key_string)
	// 			quee = quee + Http_body_buffer + "\n"
	// 			//						println("aaaaaaaaa--fail_countfail_countfail_countfail_countfail:",quee)
	// 		}
	// 	}
	// 	//			println("fail_countfail_countfail_countfail_countfail_countfail_countfail_countfail_count:",fail_count)
	// 	//			println("array_countarray_countarray_countarray_countarray_countarray_countarray_countarray_countarray_count:",array_count)
	// 	if array_count == fail_count+1 {
	// 		println("Primary is dead!!!!") // test output
	// 	} else if array_count/2 <= fail_count+1 {
	// 		println("Primary has something wrong!!!") // test output
	// 	}
	// }

	// Hostlist = quee
	// return Hostlist

	return ""
}

// GetInfoPrimary ..
func GetInfoPrimary(group, user, zone, platform string) {
	/// initial primary info
	GetHost(group, user, zone, platform, "")
	GetPrimary(group, user, zone, platform)
}

// PutPlatformInit ..
func PutPlatformInit(platform, data string) {
	// uri := "/v1/kv/klevr/systems/platform_types/" + platform
	// communicator.Put_http(ConsulURL+uri, data, API_key_string)
}

// PutPrimaryAck ..
func (api *API) PutPrimaryAck(group int64, user int64, zone, platform, ack string) {
	logger.Debug(fmt.Sprintf("group : %d, user : %d", group, user))

	var ma = &model.PrimaryAgents{
		GroupId:        group,
		AgentId:        user,
		LastAccessTime: time.Now().UTC(),
	}

	api.DB.LogMode(true)

	api.DB.Where(&model.PrimaryAgents{
		GroupId: group,
		AgentId: user,
	}).First(&ma)

	logger.Debug(ma)

	api.DB.Model(&ma).Updates(model.PrimaryAgents{LastAccessTime: time.Now().UTC()})

	// uri := "/v1/kv/klevr/groups/" + group + "/users/" + user + "/zones/" + zone + "/platforms/" + platform + "/primary_ack"
	// communicator.Put_http(ConsulURL+uri, ack, API_key_string)
}

// HostpoolMgt Old hostlist purge
func HostpoolMgt(group, user, zone, platform string) string {
	// /// Define variables
	// var arr []string
	// var queue, target_key string
	// Host_purge_result = "\n"

	// /// Get Hostlist with Keys
	// dataJson := communicator.Get_http(ConsulURL+"/v1/kv/klevr/groups/"+group+"/users/"+user+"/zones/"+zone+"/platforms/"+platform+"/hosts/?keys", API_key_string)
	// _ = json.Unmarshal([]byte(dataJson), &arr)
	// for i := 0; i < len(arr); i++ {
	// 	var target_txt, time_arry []string
	// 	var time_string string
	// 	endpoint := arr[i][strings.LastIndex(arr[i], "/")+1:]
	// 	if endpoint == "health" {
	// 		queue = communicator.Get_http(ConsulURL+"/v1/kv/"+arr[i]+"?raw=1", API_key_string) /// Endpoing value will be "~/health" part of API
	// 		get_data := arr[i]

	// 		/// Get value of each hosts
	// 		target_key = ConsulURL + "/v1/kv/" + get_data
	// 		println("target_key=", target_key) ///////////  Test output
	// 		/// Parsing the Key/value of host_info
	// 		target_txt = strings.Split(string(queue), "&")
	// 		time_arry = strings.Split(target_txt[0], "=")

	// 		/// Parsing the Key/value for Unix Time
	// 		time_string = string(time_arry[1])
	// 		time_parsing, err := strconv.ParseInt(time_string, 10, 64)
	// 		if err != nil {
	// 			log.Println(err)
	// 		}
	// 		/// Duration check
	// 		tm := time.Unix(time_parsing, 0)
	// 		if time.Since(tm).Hours() > 1 {
	// 			/// Delete old host via API server
	// 			Host_purge_result = Host_purge_result + "Overtime: " + get_data + "\n"
	// 			communicator.Delete_http(ConsulURL+"/v1/kv/"+get_data, API_key_string)
	// 		} else {
	// 			Host_purge_result = Host_purge_result + "It's ok: " + get_data + "\n"
	// 		}
	// 	}
	// }
	// return Host_purge_result

	return ""
}

// ClientReceiver ..
func ClientReceiver(group, user, zone, hostname, hostIP, platform, hostAlive, primaryAlive string) string {
	// uri := "/v1/kv/klevr/groups/" + group + "/users/" + user + "/zones/" + zone + "/platforms/" + platform + "/hosts/" + hostname + "/health"
	// data := "last_check=" + host_alive + "&ip=" + host_ip + "&clientType=" + platform + "&primaryConnection=" + primary_alive
	// communicator.Put_http(ConsulURL+uri, data, API_key_string)
	// Buffer_result = data
	// return Buffer_result

	return ""
}

// PutHostinfo ..
func PutHostinfo(group, user, zone, platform, hostname, body string) string {
	// uri := "/v1/kv/klevr/groups/" + group + "/users/" + user + "/zones/" + zone + "/platforms/" + platform + "/hosts/" + hostname + "/hostinfo"
	// data := body
	// communicator.Put_http(ConsulURL+uri, data, API_key_string)
	// Buffer_result = data
	// return Buffer_result

	return ""
}
