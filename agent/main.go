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
	netutil "k8s.io/apimachinery/pkg/util/net"
) 


var klevr_agent_id_file = "/tmp/klevr_agent.id"
var klevr_task_dir = "/tmp/klevr_task"
var klevr_agent_conf_file = "/tmp/klevr_agent.conf"
var klevr_agent_id_string string

var klevr_server_addr = "localhost:8080"
var klevr_console = "http://"+klevr_server_addr
var api_key_string string
var local_ip_add string
var account_n string
var svc_provider string
var installer string
var Master_ip string

/// debug_mode = dev or not
var debug_mode string = "dev" 


/// Function for Debuging
func Debug(output string){
	if debug_mode == "dev"{
		log.Println("DEBUG:",output)
	}
}


func check(e error) {
	if e != nil {
		panic(e)
//		log.Printf(" - unknown error")
	}
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
	err := ioutil.WriteFile(klevr_agent_id_file, []byte(hex.EncodeToString(hashed) + "\n"), 0644)
	check(err)
}


// Find out the IP mac_addess
func Check_variable() string{
	// get Local IP address automatically
	default_ip,err := netutil.ChooseHostInterface()
        if err != nil {
                log.Fatalf("Failed to get IP address: %v", err)
        }

	// Flag options
	userid := flag.String("user", "", "Account key from Klevr service")
	provider := flag.String("provider", "", "[baremetal|aws] - Service Provider for Host build up")
	local_ip := flag.String("ip", default_ip.String(), "local IP address for networking")
	klevr_addr := flag.String("webconsole", klevr_server_addr, "Klevr webconsole(server) address (URL or IP, Optional: Port) for connect")

	//var klevr_server_addr = "localhost:8080"

	flag.Parse() // Important for parsing

	// Need to switch for the slave-server list update to API
	local_ip_add = *local_ip

	// Check the null data from CLI
	if len(*userid) == 0 {
		fmt.Println("Please insert an AccountID")
		os.Exit(0)
	}
	if len(*provider) == 0 {
		fmt.Println("Please make sure the provider")
		os.Exit(0)
	}
	if len(*local_ip) == 0 {
		local_ip_add = default_ip.String()
	}
	if len(*klevr_addr) > 0 {
		klevr_server_addr = *klevr_addr
	}// else if len(*klevr_addr) == 0 {
	//	klevr_server_addr = klevr_server_addr
//	}

	// Check for the Print
	account_n = *userid
	fmt.Println("Account:",account_n)
	mca := Get_mac()
	//base_info := "User Account ID + MAC address as a HW + local IP address"
	base_info := *userid + mca + *local_ip
	_, err = ioutil.ReadFile(klevr_agent_id_file)
	if err != nil{
		hash_create(base_info)
	}
	svc_provider = string(*provider)

	return svc_provider
	return local_ip_add
	return account_n
	return klevr_server_addr
	return api_key_string
}

func Klevr_agent_id_get() string{
	klevr_agent_id, _ := ioutil.ReadFile(klevr_agent_id_file)
	string_parse := strings.Split(string(klevr_agent_id),"\n")
	klevr_agent_id_string = string_parse[0] 
	return klevr_agent_id_string
}

func Basement(){
	os.MkdirAll(klevr_task_dir, 600)
}

func Chk_inst() string{
	cmm := exec.Command("which","apt-get")
	err := cmm.Run()
	if err != nil {
		installer = "yum"
	} else {
		installer = "apt-get"
	}
	return installer
}


func Chk_pkg(pkg string){
	Chk_inst()
	cmm := exec.Command("which", pkg)
	cmm.Env = append(os.Environ())
	if err := cmm.Run(); err != nil {
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
        if installer == "apt-get" {
                log.Printf("- Please wait for the %s update",installer)
                update := exec.Command("sudo",installer,"update")
                update.Run()
        }
        log.Printf("- Please wait for Installing the %s Package....", packs)
        cmd := exec.Command("sudo",installer,"install","-y",packs)
        err := cmd.Run()
        if err != nil{
                log.Printf("- Command finished with error for %s: %v", packs, err)
        }else {
                log.Printf("- \"%s\" package has been installed",packs)
        }
}

func Alive_chk_to_mgm(fail_chk string) {
	now_time := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	uri := fmt.Sprint(klevr_console+"/user/"+account_n+"/hostname/"+klevr_agent_id_string+"/"+local_ip_add+"/type/"+svc_provider+"/"+now_time+"/"+fail_chk)
	Debug(uri) /// log output
	communicator.Get_http(uri, api_key_string)
}

func Conf_manager() string{
	Master_ip = communicator.Get_http(klevr_console+"/user/"+account_n+"/masterinfo", api_key_string)
	Debug(Master_ip) /// log output
	return Master_ip
}


//func Turn_on (){
//	Alive_chk_to_mgm("ok")
//}


func main(){
	Check_variable()
	Basement()
	Klevr_agent_id_get()
	Chk_pkg("docker")
	Conf_manager()
	//Chk_pkg("asciinema") /// for test
	println("provider: ", svc_provider)
	println("local_ip_add:", local_ip_add)
	println("Agent UniqID:", klevr_agent_id_string)
	println("Master:", Master_ip)
	s := gocron.NewScheduler()
	s.Every(1).Seconds().Do(Conf_manager)
//	s.Every(1).Seconds().Do(Turn_on)
	<- s.Start()
}


