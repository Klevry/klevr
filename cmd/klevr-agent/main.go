package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"github.com/Klevry/klevr/pkg/agent"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	//"bytes"
	//"crypto/sha1"
	//"encoding/base64"
	//"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/mackerelio/go-osstat/memory"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	//"net"
	//"net/http"
	"os"
	//"os/exec"
	_ "regexp"
	"runtime"
	//"strconv"
	//"strings"
	"syscall"
	//"time"

	"github.com/Klevry/klevr/pkg/communicator"
	netutil "k8s.io/apimachinery/pkg/util/net"
	//"github.com/mackerelio/go-osstat/cpu"
	//"github.com/mackerelio/go-osstat/disk"
)

var AGENT_VERSION = "0.0.1"

var Klevr_agent_id_file = "/tmp/klevr_agent.id"
var Klevr_task_dir = "/tmp/klevr_task"
var Klevr_agent_conf_file = "/tmp/klevr_agent.conf"
var Primary_communication_result = "/tmp/communication_result.stmp"
var Prov_script = "https://raw.githubusercontent.com/Klevry/klevr/master/scripts"
var Klevr_primary_info = "/tmp/klevr_primary_info"
var Primary_alivecheck = "/tmp/primary_alivecheck_timestamp"

//var Prov_script = "https://raw.githubusercontent.com/folimy/klevr/master/provisioning_lists"
var Timestamp_from_Primary = "/tmp/timestamp_from_primary.stmp"
var Klevr_tmp_manager string
var Cluster_info = "/tmp/cluster_info"
var SSH_provbee = "ssh provbee-service "
var Commands = "/tmp/command"


var Klevr_agent_id_string string

var Klevr_manager string
var Api_key_string string
var Local_ip_add string
var API_key_id string
var Platform_type string
var Klevr_zone string
var Klevr_company string

var Installer string
var Primary_ip string
var AM_I_PRIMARY string
var System_info string
var Error_buffer string
var Result_buffer string

var Body common.Body
var Primary_alivecheck_time int64
var ping bool
var doOnce sync.Once

///// Mode_debug = dev or not
//var Mode_debug string = "dev"
//

func Command_checker(cmd, msg string) (string, error) {
	chk_command := exec.Command("sh", "-c", cmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	chk_command.Stdout = &out
	chk_command.Stderr = &stderr
	err := chk_command.Run()
	if err != nil {
		logger.Debugf("%v", err)
		//		panic(msg)
	}
	Result_buffer = out.String()
	Error_buffer = msg
	return Error_buffer, err
	return Result_buffer, err
}

//func Required_env_chk(){
//	Command_checker("egrep '(vmx|svm)' /proc/cpuinfo", "Error: Required VT-X. Please check the BIOS or check the other machine.")
//	Command_checker("echo 'options kvm_intel nested=1' >> /etc/modprobe.d/kvm-nested.conf;modprobe -r kvm_intel && modprobe kvm_intel", "Error: Required apply of modprobe command." )
//	Command_checker("cat /sys/module/kvm_intel/parameters/nested", "Error: Required check for this file - /sys/module/kvm_intel/parameters/nested for \"Y\"")
//}

func Get_mac() (mac_add string) {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, i := range interfaces {
			if i.Flags&net.FlagUp != 0 && bytes.Compare(i.HardwareAddr, nil) != 0 {
				mac_add = i.HardwareAddr.String()
				break
			}
		}
	}
	return mac_add
}

func hash_create(s string) {
	h := sha1.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	err := ioutil.WriteFile(Klevr_agent_id_file, []byte(hex.EncodeToString(hashed)+"\n"), 0644)
	logger.Errorf("%v", err)
}

// Find out the IP mac_addess
func Check_variable() string {
	// get Local IP address automatically
	default_ip, err := netutil.ChooseHostInterface()
	if err != nil {
		log.Fatalf("Failed to get IP address: %v", err)
	}

	// Flag options
	// Sample: -apiKey=\"{apiKey}\" -platform={platform} -manager=\"{managerUrl}\" -zoneId={zoneId}
	apikey := flag.String("apiKey", "", "API Key from Klevr service")
	platform := flag.String("platform", "", "[baremetal|aws] - Service Platform for Host build up")
	zone := flag.String("zoneId", "dev-zone", "zone will be a [Dev/Stg/Prod]")
	local_ip := flag.String("ip", default_ip.String(), "local IP address for networking")
	klevr_addr := flag.String("manager", Klevr_tmp_manager, "Klevr webconsole(server) address (URL or IP, Optional: Port) for connect")

	flag.Parse() // Important for parsing

	// Check the null data from CLI
	if len(*apikey) == 0 {
		fmt.Println("Please insert an API Key")
		os.Exit(0)
	}
	if len(*platform) == 0 {
		fmt.Println("Please make sure the platform")
		os.Exit(0)
	}
	if len(*local_ip) == 0 {
		Local_ip_add = default_ip.String()
	} else {
		Local_ip_add = *local_ip
	}

	if len(*klevr_addr) == 0 {
		Klevr_tmp_manager = Klevr_tmp_manager
	} else {
		Klevr_tmp_manager = *klevr_addr
	}

	Klevr_manager = Klevr_tmp_manager

	// Check for the Print
	API_key_id = *apikey
	fmt.Println("Account:", API_key_id)
	mca := Get_mac()
	//base_info := "User Account ID + MAC address as a HW + local IP address"
	base_info := *apikey + mca + *local_ip
	_, err = ioutil.ReadFile(Klevr_agent_id_file)
	if err != nil {
		hash_create(base_info)
	}
	Platform_type = string(*platform)
	Klevr_zone = string(*zone)

	return Platform_type
	return Local_ip_add
	return API_key_id
	return Klevr_manager
	return Klevr_zone

	return Api_key_string
}

func Klevr_agent_id_get() string {
	klevr_agent_id, _ := ioutil.ReadFile(Klevr_agent_id_file)
	string_parse := strings.Split(string(klevr_agent_id), "\n")
	Klevr_agent_id_string = string_parse[0]
	return Klevr_agent_id_string
}

//Provisioning file download
func Get_provisionig_script() {
	urli := Prov_script + "/" + Platform_type
	Get_script := communicator.Get_http(urli, Api_key_string)
	//Command_checker(Get_script_arr, "Error: Provisioning has been failed")
	Get_script_arr := strings.Split(strings.Replace(Get_script, "\n\n", "\n", -1), "\n")
	println("%%%%%%%%%%%%%%%%%%%: ", len(Get_script_arr))

	for i := 0; i < len(Get_script_arr); i++ {
		if len(Get_script_arr[i]) > 1 {
			fin_arr := strings.Split(Get_script_arr[i], ",")
			// println("::::::::::::::::::: eval "+fin_arr[0], fin_arr[1])
			_, err := Command_checker(fin_arr[0], fin_arr[1])
			if err != nil {
				os.Exit(1)
			}
		}

	}
}

////Klevr_company Klevr_zone
//func Alive_chk_to_mgm(fail_chk string) {
//	now_time := strconv.FormatInt(time.Now().UTC().Unix(), 10)
//	uri := fmt.Sprint(Klevr_console + "/group/"  + "/user/" + API_key_id + "/zone/" + Klevr_zone + "/platform/" + Platform_type + "/hostname/" + Klevr_agent_id_string + "/" + Local_ip_add + "/" + now_time + "/" + fail_chk)
//	Debug(uri) /// log output
//	communicator.Get_http(uri, Api_key_string)
//}

//func Resource_chk_to_mgm() {
//	uri := fmt.Sprint(Klevr_console + "/group/"  + "/user/" + API_key_id + "/zone/" + Klevr_zone + "/platform/" + Platform_type + "/hostname/" + Klevr_agent_id_string + "/hostinfo")
//	Debug(uri) /// log output
//	//Resource_info()
//	communicator.Put_http(uri, System_info, Api_key_string)
//	Debug("System_info:" + System_info) /// log output
//}

//func Resource_info() string {
//	var si sysinfo.SysInfo
//	si.GetSysInfo()
//	data, err := json.Marshal(&si)
//	if err != nil {
//		log.Fatal(err)
//	}
//	System_info = fmt.Sprintf("%s", data)
//	return System_info
//}

//func Primary_ack_stamping(){
//	primary_ack_time := fmt.Sprint(time.Now().Unix())
//        err := ioutil.WriteFile(Primary_status_file, []byte(primary_ack_time), 0644)
//	println(err)
//}


//func Hosts_alive_list(alive_list string) {
//	//  Hosts alive list klevr/groups/klevr-a-team/users/ralf/zones/dev/platforms/baremetal/alive_hosts
//	uri := fmt.Sprint(Klevr_console + "/groups/"  + "/users/" + API_key_id + "/zones/" + Klevr_zone + "/platforms/" + Platform_type + "/aliveagent")
//	Debug(uri) /// log output
//	alive_conv := fmt.Sprintf("%s", alive_list)
//	communicator.Put_http(uri, alive_conv, Api_key_string)
//}

//func RnR() {
//	Check_primary()
//	if AM_I_PRIMARY == "PRIMARY" {
//		// Put primary alive time to stamp
//		ack_timecheck_from_api := communicator.Get_http(Klevr_console+"/group/"+Klevr_company+"/user/"+API_key_id+"/zone/"+Klevr_zone+"/platform/"+Platform_type+"/ackprimary", Api_key_string)
//
//		// Write done the information about of Final result time & hostlists
//		ioutil.WriteFile(Primary_communication_result, []byte(ack_timecheck_from_api), 0644)
//
//		Secondary_scanner()
//
//		Alive_chk_to_mgm("ok")
//		if Platform_type == "baremetal" {
//			//			println ("Docker_runner here - klevr_beacon_img")
//			//Docker_runner("klevry/beacon:latest", "primary_beacon", "-p 18800:18800 -v /tmp/status:/info") // no use anymore. process has been changed to goroutin.
//			println("Docker_runner here - klevr_taskmanager_img")
//			println("Get_task_from_here for baremetal")
//		} else if Platform_type == "aws" {
//			println("Get_task_from_here for AWS")
//		}
//		println("Get_task_excution_from_here")
//		Debug("I am Primary")
//		//Resource_info() /// test
//		Resource_chk_to_mgm()
//	} else {
//		/// http://192.168.1.22:18800/primaryworks
//		// url := "http://"+Primary_ip+":18800/status"
//		url := "http://" + Primary_ip + ":18800/primaryworks"
//		primary_time_check := communicator.Get_http(url, Api_key_string)
//
//		//		fmt.Printf("++++++++++++++++++++++++++++++++++++++++++++++++ %s ]]]\n", primary_time_check)
//		/// Duration check
//		//		primary_time, _ := strconv.Atoi(primary_time_check)
//		primary_time, _ := strconv.ParseInt(primary_time_check, 10, 64)
//
//		fmt.Printf("++++++++++++++++++++++++++++++++++++++++++++++++ %d ]]]\n", primary_time)
//		var Host_purge_result string
//
//		/// Primary Last working time stamp
//		if primary_time != 0 {
//			ioutil.WriteFile(Timestamp_from_Primary, []byte(primary_time_check), 0644)
//		}
//
//		primary_time_result, _ := ioutil.ReadFile(Timestamp_from_Primary)
//		prim_string := string(primary_time_result)
//		primary_int, _ := strconv.ParseInt(prim_string, 10, 64)
//
//		tm := time.Unix(primary_int, 0)
//		if time.Since(tm).Minutes() > 1 {
//			/// Delete old host via API server
//			Host_purge_result = Primary_ip + ": Primary agent is not working!!\n"
//		} else {
//			//Host_purge_result = Host_purge_result+"It's ok: "+get_data+"\n"
//			Host_purge_result = Primary_ip + ": Primary agent is working hard :) \n"
//		}
//
//		println("Error check for Debug:", Host_purge_result)
//		// Primary error checker here - 2020/6/25
//		Debug("I am Secondary")
//		//		Resource_info() /// test
//		Resource_chk_to_mgm()
//		//		Debug(aaa)
//	}
//}

// Docker image pull
func Docker_pull(image_name string) {
	log.Printf("- %s docker image pulling now. Please wait...", image_name)
	pulling_image := exec.Command("docker", "pull", image_name)
	pulling_image.Stdout = os.Stdout
	err := pulling_image.Run()
	if err != nil {
		log.Printf("- %s docker image not existed in the registry. Please check the image name or network connection.", image_name)
		os.Exit(1)
	} else {
		log.Printf("- Docker image has been pulled.")
	}
}

//// Docker image runner
//func Docker_runner(image_name, service_name, options string) {
//	docker_ps_command := "docker ps | grep " + image_name + "|egrep -v CONTAINER | head -1"
//	Command_checker(docker_ps_command, "Error: Docker running process check failed")
//	if len(Result_buffer) != 0 {
//		Debug(image_name + " docker container is running now.")
//	} else {
//		Docker_pull(image_name)
//		Command_checker("docker run -d --name "+service_name+" "+options+" "+image_name, "\"- %s container already existed. Please check the docker process.\", image_name")
//	}
//}

///// Primary last working time checker
//func Primary_works_check() string {
//	var primary_latest_check string
//	primary_raw_file, _ := ioutil.ReadFile(Primary_communication_result)
//	raw_string_parse := strings.Split(string(primary_raw_file), "\n")
//	if strings.Contains(raw_string_parse[0], "get_timestamp") == true {
//		strr1 := strings.Split(raw_string_parse[0], ": ")
//		primary_latest_check = strr1[1]
//	} else {
//		log.Println("Primary uptime is not recognized")
//		primary_latest_check = ""
//	}
//	return primary_latest_check
//}



/*
=======================================================================
*/

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func Check_primary() string {
	if Primary_ip == Local_ip_add {
		AM_I_PRIMARY = "true"
		log.Printf("--------------------------------  Primary_ip=%s, Local_ip_add=%s", Primary_ip, Local_ip_add)
	} else if Primary_ip != Local_ip_add {
		AM_I_PRIMARY = "false"
		log.Printf("--------------------------------  Primary_ip=%s, Local_ip_add=%s", Primary_ip, Local_ip_add)
	}
	return AM_I_PRIMARY
}

func PingToMaster(){
	timeout := time.Duration(1 * time.Second)

	alive := agent.AliveCheck{}

	_, err := net.DialTimeout("tcp", Primary_ip+":18800", timeout)
	if err != nil {
		logger.Error(err)
		alive.IsActive = false
	} else {
		alive.IsActive = true
		Primary_alivecheck_time = time.Now().Unix()
		alive.Time = Primary_alivecheck_time
	}

	m, _ := json.MarshalIndent(alive, "", "  ")
	err2 := ioutil.WriteFile(Primary_alivecheck, m, os.FileMode(0644))
	if err2 != nil{
		logger.Debugf("%v", err)
	}

	logger.Debug("%v", Primary_ip)
}

func SendMe(body *common.Body) {
	body.Me.IP = Local_ip_add
	body.Me.IP = Local_ip_add
	body.Me.Port = 18800
	body.Me.Version = AGENT_VERSION

	disk := DiskUsage("/")

	memory, _ := memory.Get()

	body.Me.Resource.Core = runtime.NumCPU()
	body.Me.Resource.Memory = int(memory.Total/MB)
	body.Me.Resource.Disk = int(disk.All/MB)
}

/*
in: body.me
out: body.me, body.agent.primary
 */
func HandShake(){

	uri := Klevr_manager + "/agents/handshake"

	rb := &common.Body{}

	SendMe(rb)

	logger.Debugf("%v", rb)

	b, err := json.Marshal(rb)
	if err != nil {
		logger.Error(err)
	}

	result := communicator.Put_Json_http(uri, b, Klevr_agent_id_get(), API_key_id, Klevr_zone)

	err2 := json.Unmarshal(result, &Body)
	if err2 != nil{
		logger.Error(err2)
	}

	logger.Debugf("%v", Body.Agent.Primary)
	Primary_ip = Body.Agent.Primary.IP
}

/*
in: body.me, body.agent.nodes, body.task
out: body.me, body.task
 */
func TaskManagement(){
	//uri := Klevr_manager + "/agents/" + Klevr_agent_id_get()

	PingToMaster()

	rb := &common.Body{}

	SendMe(rb)

	//rb.Agent.Nodes = GetNodes(uri)

	logger.Debugf("%v", rb.Agent.Nodes)

	//b, err := json.Marshal(rb)
	//if err != nil {
	//	logger.Error(err)
	//}

	//communicator.Put_Json_http(uri, b, Klevr_agent_id_get())
}

/*
in: body.me, body.agent.primary
out: body.me, body.agent.primary
 */
func PrimaryStatusReport(){
	//uri := Klevr_manager + "/agents/" + Klevr_agent_id_get()

	alivecheck, err := ioutil.ReadFile(Primary_alivecheck)
	if err != nil{
		logger.Error(err)
	}

	var alive agent.AliveCheck
	json.Unmarshal(alivecheck, &alive)

	rb := &common.Body{}

	SendMe(rb)
	//rb.Agent.Primary = GetPrimary(uri)

	if alive.IsActive{

	}

}

func printprimary(){
	if(Check_primary() == "true"){
		logger.Debugf("-----------I am Primary")
	} else {
		logger.Debugf("-----------I am Secondary")
	}
	logger.Debugf("Primary ip : %s, My ip : %s", Primary_ip, Local_ip_add)
}


func getCommand(){
	uri := Klevr_manager + "/agents/commands/init"

	result := communicator.Get_Json_http(uri, Klevr_agent_id_get(), API_key_id, Klevr_zone)

	err := json.Unmarshal(result, &Body)
	if err != nil{
		logger.Error(err)
	}

	coms := Body.Task[0].Command
	com := strings.Split(coms, "\n")

	filenum := len(com)

	for i:=0; i<filenum-1; i++{
		num := strconv.Itoa(i)

		var read string

		err := json.Unmarshal(readFile(Commands+num), &read)
		if err != nil{
			logger.Error(err)
		}

		if(com[i] == read){
			logger.Debugf("same command")
		} else {
			logger.Debugf("%d-----%s", i, com[i])
			writeFile(Commands+num, com[i])

			execute := SSH_provbee + com[i]
			//execute := com[i]

			exe := exec.Command("sh", "-c", execute)
			errExe := exe.Run()
			if errExe != nil{
				logger.Error(errExe)
			} else {
				res := primaryInit(Body, coms, "done")
				logger.Debugf("%v", string(res))

			}

		}
	}
}

func primaryInit(bod common.Body, command string, status string) []byte{
	uri := Klevr_manager + "/agents/zones/init"

	rb := &common.Body{}

	SendMe(rb)

	rb.Task = make([]common.Task, 1)
	rb.Task[0].ID = bod.Task[0].ID
	rb.Task[0].AgentKey = bod.Task[0].AgentKey
	rb.Task[0].Command = command
	rb.Task[0].Status = status
	rb.Task[0].Params = bod.Task[0].Params
	rb.Task[0].Result = bod.Task[0].Result
	rb.Task[0].Type = bod.Task[0].Type

	b, err := json.Marshal(rb)
	if err != nil {
		logger.Error(err)
	}

	logger.Debugf("%v", rb)

	result := communicator.Post_Json_http(uri, b, Klevr_agent_id_get(), API_key_id, Klevr_zone)
	return result
}

func writeFile(path string, data string) {
	d, _ := json.MarshalIndent(data, "", "  ")
	err := ioutil.WriteFile(path, d, os.FileMode(0644))
	if err != nil{
		logger.Error(err)
	}
}

func readFile(path string) []byte{
	data, err := ioutil.ReadFile(path)
	if err != nil{
		logger.Error(err)
	}

	return data
}


func deleteFile(path string){
	err := os.Remove(path)
	if err != nil {
		logger.Error(err)
	}

}

func main() {
	/// check the cli command with required options
	Check_variable()
	///// Requirement package check
	//if Platform_type == "baremetal" {
	//	Check_package("curl")
	//	Check_package("docker")
	//}
	//
	/// Checks env. for baremetal to Hypervisor provisioning
	//Get_provisionig_script()
	//
	///// Set up the Task & configuration directory
	//Set_basement()
	//
	///// Uniq ID create & get
	//Klevr_agent_id_get()
	//
	///// Check for primary info
	//Alive_chk_to_mgm("ok")
	//Resource_chk_to_mgm()
	//Get_primaryinfo()

	println("platform: ", Platform_type)
	println("Local_ip_add:", Local_ip_add)
	println("Agent UniqID:", Klevr_agent_id_string)
	println("Primary:", Primary_ip)

	HandShake()
	getCommand()


	//if Check_primary() == "true"{
	//	s := gocron.NewScheduler()
	//	s.Every(5).Seconds().Do(printprimary)
	//	s.Every(5).Seconds().Do(getCommand)
	//
	//	go func() {
	//		<-s.Start()
	//	}()
	//} else {
	//	s := gocron.NewScheduler()
	//	//s.Every(1).Seconds().Do(PingToMaster)
	//	s.Every(1).Seconds().Do(printprimary)
	//	//s.Every(1).Seconds().Do(getCommand)
	//	//s.Every(1).Seconds().Do(slave)
	//
	//	go func() {
	//		<-s.Start()
	//	}()
	//}

	http.ListenAndServe(":18800", nil)

}