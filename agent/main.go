package main

import (
	"os"
	"os/exec"
	"fmt"
        "flag"
	"net"
	"log"
	"time"
	"io/ioutil"
	"crypto/sha1"
	"encoding/hex"
	"bytes"
	"github.com/ralfyang/klevr/communicator"
	"strings"
	"strconv"
	"github.com/jasonlvhit/gocron"
	"net/http"
	_"regexp"
	netutil "k8s.io/apimachinery/pkg/util/net"
	"github.com/zcalusic/sysinfo"
	"encoding/json"
	//"github.com/mackerelio/go-osstat/memory"
	//"github.com/mackerelio/go-osstat/cpu"
	//"github.com/mackerelio/go-osstat/disk"
)


var Klevr_agent_id_file = "/tmp/klevr_agent.id"
var Klevr_task_dir = "/tmp/klevr_task"
var Klevr_agent_conf_file = "/tmp/klevr_agent.conf"
var Primary_communication_result = "/tmp/communication_result.stmp"

var Klevr_agent_id_string string

var Klevr_console string
var Api_key_string string
var Local_ip_add string
var User_account_id string
var Platform_type string
var Klevr_zone string
var Klevr_company string


var Installer string
var Primary_ip string
var AM_I_PRIMARY string
var System_info string
var Error_buffer string
var Result_buffer string

/// Mode_debug = dev or not
var Mode_debug string = "dev" 


/// Function for Debug
func Debug(output string){
	if Mode_debug == "dev"{
		log.Println("DEBUG:",output)
	}
}


func check(e error) {
	if e != nil {
		panic(e)
//		log.Printf(" - unknown error")
	}
}

func Command_checker(cmd, msg string) string{
	chk_command := exec.Command("sh","-c",cmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	chk_command.Stdout = &out
	chk_command.Stderr = &stderr
	err := chk_command.Run()
	if err != nil {
		log.Printf(msg)
//		panic(msg)
	}
	Result_buffer = out.String()
	Error_buffer = msg
	return Error_buffer
	return Result_buffer
}

func Required_env_chk(){
	Command_checker("egrep '(vmx|svm)' /proc/cpuinfo", "Error: Required VT-X. Please check the BIOS or check the other machine.")
	Command_checker("echo 'options kvm_intel nested=1' >> /etc/modprobe.d/kvm-nested.conf;modprobe -r kvm_intel && modprobe kvm_intel", "Error: Required apply of modprobe command." )
	Command_checker("cat /sys/module/kvm_intel/parameters/nested", "Error: Required check for this file - /sys/module/kvm_intel/parameters/nested for \"Y\"")
}


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


func hash_create(s string){
	h := sha1.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	err := ioutil.WriteFile(Klevr_agent_id_file, []byte(hex.EncodeToString(hashed) + "\n"), 0644)
	check(err)
}


// Find out the IP mac_addess
func Check_variable() string{
	// get Local IP address automatically
	default_ip,err := netutil.ChooseHostInterface()
	klevr_tmp_server := "localhost:8080"
        if err != nil {
                log.Fatalf("Failed to get IP address: %v", err)
        }

	// Flag options
	userid := flag.String("id", "", "Account ID from Klevr service")
	platform := flag.String("platform", "", "[baremetal|aws] - Service Platform for Host build up")
	company := flag.String("group", "", "Group name will be a uniq company name or team name")
	zone := flag.String("zone", "dev-zone", "zone will be a [Dev/Stg/Prod]")
	local_ip := flag.String("ip", default_ip.String(), "local IP address for networking")
	klevr_addr := flag.String("webconsole", klevr_tmp_server, "Klevr webconsole(server) address (URL or IP, Optional: Port) for connect")


	flag.Parse() // Important for parsing

	// Check the null data from CLI
	if len(*userid) == 0 {
		fmt.Println("Please insert an AccountID")
		os.Exit(0)
	}
	if len(*company) == 0 {
		fmt.Println("Please make sure the group name")
		os.Exit(0)
	}
	if len(*platform) == 0 {
		fmt.Println("Please make sure the platform")
		os.Exit(0)
	}
	if len(*local_ip) == 0 {
		Local_ip_add = default_ip.String()
	}else{
		Local_ip_add = *local_ip
	}


	if len(*klevr_addr) == 0 {
		klevr_tmp_server = klevr_tmp_server
	}else{
		klevr_tmp_server = *klevr_addr
	}

	Klevr_console = "http://"+klevr_tmp_server

	// Check for the Print
	User_account_id = *userid
	fmt.Println("Account:",User_account_id)
	mca := Get_mac()
	//base_info := "User Account ID + MAC address as a HW + local IP address"
	base_info := *userid + mca + *local_ip
	_, err = ioutil.ReadFile(Klevr_agent_id_file)
	if err != nil{
		hash_create(base_info)
	}
	Platform_type = string(*platform)
	Klevr_zone = string(*zone)
	Klevr_company = string(*company)

	return Platform_type
	return Local_ip_add
	return User_account_id
	return Klevr_console
	return Klevr_zone
	return Klevr_company

	return Api_key_string
}

func Klevr_agent_id_get() string{
	klevr_agent_id, _ := ioutil.ReadFile(Klevr_agent_id_file)
	string_parse := strings.Split(string(klevr_agent_id),"\n")
	Klevr_agent_id_string = string_parse[0]
	return Klevr_agent_id_string
}

func Set_basement(){
	os.MkdirAll(Klevr_task_dir, 600)
}

func Chk_inst() string{
	docker_ps_command := exec.Command("which","apt-get")
	err := docker_ps_command.Run()
	if err != nil {
		Installer = "yum"
	} else {
		Installer = "apt-get"
	}
	return Installer
}


func Check_package(pkg string){
	Chk_inst()
	docker_ps_command := exec.Command("which", pkg)
	docker_ps_command.Env = append(os.Environ())
	if err := docker_ps_command.Run(); err != nil {
		if pkg == "docker" {
			log.Printf("- Package install for %s", pkg)
			Manual_inst("https://bit.ly/startdocker", "docker")
		}else{
			Install_pkg(pkg)
		}
	}
}

func Manual_inst(uri, target string){
	exec_file := "/tmp/temporary_file_for_install.sh"
	m_down := exec.Command("curl","-sL",uri,"-o",exec_file)
	m_down.Run()
	if err := os.Chmod(exec_file, 0755); err != nil {
		check(err)
	}
	m_inst := exec.Command("bash",exec_file)
	m_inst.Stdout = os.Stdout
	m_inst.Run()

	check_command := exec.Command("which", target)
	if err := check_command.Run(); err != nil {
		log.Printf("- %s package has not been installed: Please install the package manually: %s", target, target)
		os.Exit(1)
	}else{
		log.Printf("- %s package has been installed", target)
	}
}


func Install_pkg(packs string){
        if Installer == "apt-get" {
                log.Printf("- Please wait for the %s update",Installer)
                update := exec.Command("sudo",Installer,"update")
                update.Run()
        }
        log.Printf("- Please wait for Installing the %s Package....", packs)
        cmd := exec.Command("sudo",Installer,"install","-y",packs)
        err := cmd.Run()
        if err != nil{
                log.Printf("- Command finished with error for %s: %v", packs, err)
        }else {
                log.Printf("- \"%s\" package has been installed",packs)
        }
}
//Klevr_company Klevr_zone
func Alive_chk_to_mgm(fail_chk string) {
	now_time := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	uri := fmt.Sprint(Klevr_console+"/group/"+Klevr_company+"/user/"+User_account_id+"/zone/"+Klevr_zone+"/platform/"+Platform_type+"/hostname/"+Klevr_agent_id_string+"/"+Local_ip_add+"/"+now_time+"/"+fail_chk)
	Debug(uri) /// log output
	communicator.Get_http(uri, Api_key_string)
}

func Get_primaryinfo() string{
	uri_result := strings.Split(communicator.Get_http(Klevr_console+"/group/"+Klevr_company+"/user/"+User_account_id+"/zone/"+Klevr_zone+"/platform/"+Platform_type+"/primaryinfo", Api_key_string), "=")
	Primary_ip = uri_result[1]
	Debug(Primary_ip) /// log output
	return Primary_ip
}


func Check_primary() string{
        if Primary_ip == "" {
                log.Printf("- Klevr task manager has not defined. Please wait for vote from webconsole")
        }else if Primary_ip == Local_ip_add {
                AM_I_PRIMARY = "PRIMARY"
                log.Printf("--------------------------------  Primary_ip=%s, Local_ip_add=%s",Primary_ip,Local_ip_add)
        }else if Primary_ip != Local_ip_add  {
                AM_I_PRIMARY = "0"
                log.Printf("--------------------------------  Primary_ip=%s, Local_ip_add=%s",Primary_ip,Local_ip_add)
        }
        return AM_I_PRIMARY
}


func Resource_chk_to_mgm() {
	uri := fmt.Sprint(Klevr_console+"/group/"+Klevr_company+"/user/"+User_account_id+"/zone/"+Klevr_zone+"/platform/"+Platform_type+"/hostname/"+Klevr_agent_id_string+"/hostinfo")
	Debug(uri) /// log output
	Resource_info()
	communicator.Put_http(uri, System_info, Api_key_string)
	Debug("System_info:"+System_info) /// log output
}

func Resource_info()string{
	var si sysinfo.SysInfo
	si.GetSysInfo()
	data, err := json.Marshal(&si)
	if err != nil {
	    log.Fatal(err)
	}
	System_info = fmt.Sprintf("%s",data)
	return System_info
}




//func Primary_ack_stamping(){
//	primary_ack_time := fmt.Sprint(time.Now().Unix())
//        err := ioutil.WriteFile(Primary_status_file, []byte(primary_ack_time), 0644)
//	println(err)
//}



func Secondary_scanner(){
	Secondary_raw_file, _ := ioutil.ReadFile(Primary_communication_result)
	raw_string_parse := strings.Split(string(Secondary_raw_file),"\n")
	var quee_host string
	for i := 1; i < len(raw_string_parse)-2; i++ {
		var fin_res string = ""
		target_raw := raw_string_parse[i]
		strr1 := strings.Split(target_raw, "&")
		raw_result_split := strings.Split(strr1[1], "=")

		Target_secondary_hosts := "http://"+raw_result_split[1]+":18800"
		fin_res = communicator.Get_http(Target_secondary_hosts+"/status", "")
		if fin_res == "OK" {
//			println("77777777777777777777777777777777777777 Secondary_raw_fileSecondary_raw_file: ", fin_res)  /// for test result
			if i == len(raw_string_parse)-3{
				quee_host = quee_host+Target_secondary_hosts+": "+fin_res
			}else{
				quee_host = quee_host+Target_secondary_hosts+": "+fin_res+"\n"
			}
		}
	}
//	regex, _ := regexp.Compile("\n\n")
//	flat_quee_host := regex.ReplaceAllString(quee_host, "\n")
	flat_quee_host := strings.Replace(quee_host, "\n\n", "\n", -1)
	println("8888888888888888888888888888888888888888888888888888888888")
	println(flat_quee_host)
	println("9999999999999999999999999999999999999999999999999999999999")
}



func RnR(){
	Check_primary()
	if AM_I_PRIMARY == "PRIMARY" {
		// Put master alive time to stamp
		ack_timecheck_from_api := communicator.Get_http(Klevr_console+"/group/"+Klevr_company+"/user/"+User_account_id+"/zone/"+Klevr_zone+"/platform/"+Platform_type+"/ackprimary", Api_key_string)

		// Write done the information about of Final result time & hostlists
		ioutil.WriteFile(Primary_communication_result, []byte(ack_timecheck_from_api), 0644)

		Secondary_scanner()

		Alive_chk_to_mgm("ok")
		if Platform_type == "baremetal" {
//			println ("Docker_runner here - klevr_beacon_img")
			//Docker_runner("klevry/beacon:latest", "primary_beacon", "-p 18800:18800 -v /tmp/status:/info") // no use anymore. process has been changed to goroutin.
			println ("Docker_runner here - klevr_taskmanager_img")
			println ("Get_task_from_here for baremetal")
		} else if Platform_type == "aws" {
			println ("Get_task_from_here for AWS")
		}
		println ("Get_task_excution_from_here")
		Debug("I am Primary")
		Resource_info() /// test
		Resource_chk_to_mgm()
	}else{
		url := "http://"+Primary_ip+":18800/status"
	        req, _ := http.NewRequest("GET", url, nil)
	        req.Header.Add("cache-control", "no-cache")
	        _, err := http.DefaultClient.Do(req)
		if err != nil {
			Alive_chk_to_mgm("failed")
		}else{
			Alive_chk_to_mgm("ok")
		}
		// Primary error checker here - 2020/6/25 
		Debug("I am Secondary")
//		Resource_info() /// test
		Resource_chk_to_mgm()
//		Debug(aaa)
	}
}



// Docker image pull
func Docker_pull(image_name string){
	log.Printf("- %s docker image pulling now. Please wait...", image_name)
	pulling_image := exec.Command("docker", "pull", image_name)
	pulling_image.Stdout = os.Stdout
	err := pulling_image.Run()
	if err != nil {
		log.Printf("- %s docker image not existed in the registry. Please check the image name or network connection.", image_name)
		os.Exit(1)
	}else{
		log.Printf("- Docker image has been pulled.")
	}
}


// Docker image runner
func Docker_runner(image_name, service_name, options string){
	docker_ps_command := "docker ps | grep " + image_name + "|egrep -v CONTAINER | head -1"
	Command_checker(docker_ps_command, "Error: Docker running process check failed")
	if len(Result_buffer) != 0 {
		Debug(image_name+" docker container is running now.")
	}else{
		Docker_pull(image_name)
		Command_checker("docker run -d --name "+service_name+" "+options+" "+image_name, "\"- %s container already existed. Please check the docker process.\", image_name")
	}
}


func main(){
	// Docker image define
	var libvirtd = "klevry/libvirt:latest"

	/// check the cli command with required options 
	Check_variable()

	/// Checks env. for baremetal to Hypervisor provisioning
	if Platform_type == "baremetal" {
		Required_env_chk()
	}

	/// Set up the Task & configuration directory
	Set_basement()

	/// Uniq ID create & get
	Klevr_agent_id_get()

	/// Requirement package check
	Check_package("docker")
	Check_package("curl")

	if Platform_type == "baremetal" {
		Docker_runner(libvirtd, "nested_kvm", "--privileged -d -e 'container=docker' -p 18002:22 -p 18001:16509 -p 18003:5900  -v /sys/fs/cgroup:/sys/fs/cgroup:rw")
        }

	/// Check for primary info
	Alive_chk_to_mgm("ok")
	Resource_chk_to_mgm()
	Get_primaryinfo()

	println("platform: ", Platform_type)
	println("Local_ip_add:", Local_ip_add)
	println("Agent UniqID:", Klevr_agent_id_string)
	println("Primary:", Primary_ip)


	/// Scheduler
	s := gocron.NewScheduler()
	s.Every(1).Seconds().Do(Get_primaryinfo)
//	s.Every(1).Seconds().Do(Turn_on)
	s.Every(2).Seconds().Do(RnR)

	go func(){
		<- s.Start()
	}()

	///Http listen for host info get
	http.HandleFunc("/info", func(w http.ResponseWriter, req *http.Request) {
		Resource_info()
		w.Write([]byte(System_info))
	})

	///Http listen for beacon
	http.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("OK"))
	})
	http.ListenAndServe(":18800", nil)

}


