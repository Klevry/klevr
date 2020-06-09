package main

import (
	"os"
	"fmt"
        "flag"
	"net"
	"log"
	"io/ioutil"
	"crypto/sha1"
	"encoding/hex"
	"bytes"
	"github.com/ralfyang/klevr/communicator"
	netutil "k8s.io/apimachinery/pkg/util/net"
) 


var klevr_agent_id_file = "/tmp/klevr_agent.id"
var klevr_task_dir = "/tmp/klevr_tasks"
var klevr_agent_conf_file = "/tmp/klevr_agent.conf"



var api_server string
var klevr_server = "http://192.168.10.11:8080"
var api_key_string string
var local_ip_add string
var account_n string
var svc_provider string


func check(e error) {
	if e != nil {
		panic(e)
//		log.Printf(" - unknown error")
	}
}


func Get_apikey() string{
	api_key_string = communicator.Get_http(klevr_server+"/apikey", "" )
	return api_key_string
}

func Get_apiserver_info() string{
	api_server = communicator.Get_http(klevr_server+"/apiserver", api_key_string )
	return api_server
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
	return api_key_string
}

func main(){
	Get_apikey()
	Get_apiserver_info()
	Check_variable()
	println("apiserver :", api_server)
	println("apikey :", api_key_string)
	println("provider: ", svc_provider)
	println("local_ip_add", local_ip_add)
	println("account", account_n)

}


