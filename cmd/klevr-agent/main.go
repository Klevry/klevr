package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"github.com/Klevry/klevr/pkg/agent"
	"github.com/jasonlvhit/gocron"
	"os/exec"
	"strings"

	//"bytes"
	//"crypto/sha1"
	//"encoding/base64"
	//"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/mackerelio/go-osstat/memory"

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
var Primary_alivecheck = "/tmp/primary_alivecheck_timestamp"

var Klevr_tmp_manager string
var SSH_provbee = "ssh provbee-service "

var Klevr_agent_id_string string

var Klevr_manager string
var Api_key_string string
var Local_ip_add string
var API_key_id string
var Platform_type string
var Klevr_zone string
var Primary_ip string
var AM_I_PRIMARY string


var Body common.Body
var Primary_alivecheck_time int64
var primScheduler = gocron.NewScheduler()


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
	base_info := *apikey + mca + default_ip.String()
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




// disk usage of path/disk
func DiskUsage(path string) (disk agent.DiskStatus) {
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

func PingToMaster() {
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
	if err2 != nil {
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
	body.Me.Resource.Memory = int(memory.Total / MB)
	body.Me.Resource.Disk = int(disk.All / MB)
}

/*
in: body.me
out: body.me, body.agent.primary
*/
func HandShake() {

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
	if err2 != nil {
		logger.Error(err2)
	}

	logger.Debugf("%v", Body.Agent.Primary)
	Primary_ip = Body.Agent.Primary.IP
}

/*
in: body.me, body.agent.nodes, body.task
out: body.me, body.task
*/
func TaskManagement() {
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
func PrimaryStatusReport() {
	//uri := Klevr_manager + "/agents/" + Klevr_agent_id_get()

	alivecheck, err := ioutil.ReadFile(Primary_alivecheck)
	if err != nil {
		logger.Error(err)
	}

	var alive agent.AliveCheck
	json.Unmarshal(alivecheck, &alive)

	rb := &common.Body{}

	SendMe(rb)
	//rb.Agent.Primary = GetPrimary(uri)

	if alive.IsActive {

	}

}

func printprimary() {
	if Check_primary() == "true" {
		logger.Debugf("-----------I am Primary")
	} else {
		logger.Debugf("-----------I am Secondary")
	}
	logger.Debugf("Primary ip : %s, My ip : %s", Primary_ip, Local_ip_add)
}

func getCommand() {
	uri := Klevr_manager + "/agents/commands/init"

	provcheck := exec.Command("sh", "-c", "ssh provbee-service busybee beestatus hello > /tmp/con")
	errcheck :=provcheck.Run()
	if errcheck != nil {
		logger.Error(errcheck)
	}

	by := readFile("/tmp/con")
	str := strings.TrimRight(string(by), "\n")

	if strings.Compare(str, "hi") == 0{
		result := communicator.Get_Json_http(uri, Klevr_agent_id_get(), API_key_id, Klevr_zone)

		err := json.Unmarshal(result, &Body)
		if err != nil {
			logger.Error(err)
		}

		coms := Body.Task[0].Command
		com := strings.Split(coms, "\n")

		filenum := len(com)

		for i := 0; i < filenum-1; i++ {
			// num := strconv.Itoa(i)

			// var read string

			// err := json.Unmarshal(readFile(Commands+num), &read)
			// if err != nil {
			// 	logger.Error(err)
			// }

			// if com[i] == read {
			// 	logger.Debugf("same command")
			// } else {
			// 	logger.Debugf("%d-----%s", i, com[i])
			// 	writeFile(Commands+num, com[i])

			execute := SSH_provbee + com[i]
			//execute := com[i]

			exe := exec.Command("sh", "-c", execute)
			errExe := exe.Run()
			if errExe != nil {
				logger.Error(errExe)
			}

			// }
		}


		if _, err := os.Stat("/tmp/grafana"); !os.IsNotExist(err) {
			data, err := ioutil.ReadFile("/tmp/grafana")
			if err != nil {
				logger.Error(err)
			}

			//logger.Debugf("%v", string(data))

			if string(data) != "" {

				da := strings.Split(string(data), "\n")

				logger.Debugf("%v", da[0])
				primaryInit(Body, coms, "done", da[0])
			}

		}

		primScheduler.Remove(getCommand)
	}

}

func primaryInit(bod common.Body, command string, status string, param string) []byte {
	uri := Klevr_manager + "/agents/zones/init"

	rb := &common.Body{}

	SendMe(rb)

	par := make(map[string]interface{})
	par["grafana"] = param

	rb.Task = make([]common.Task, 1)
	rb.Task[0].ID = bod.Task[0].ID
	rb.Task[0].AgentKey = bod.Task[0].AgentKey
	rb.Task[0].Command = command
	rb.Task[0].Status = status
	rb.Task[0].Params = par
	rb.Task[0].Result = bod.Task[0].Result
	rb.Task[0].Type = bod.Task[0].Type

	b, err := json.Marshal(rb)
	if err != nil {
		logger.Error(err)
	}

	logger.Debugf("request body : [%s]", b)

	result := communicator.Post_Json_http(uri, b, Klevr_agent_id_get(), API_key_id, Klevr_zone)
	return result
}

func writeFile(path string, data string) {
	d, _ := json.MarshalIndent(data, "", "  ")
	err := ioutil.WriteFile(path, d, os.FileMode(0644))
	if err != nil {
		logger.Error(err)
	}
}

func readFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error(err)
	}

	return data
}

func deleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		logger.Error(err)
	}

}

func main() {
	common.InitLogger(common.NewLoggerEnv())

	/// check the cli command with required options
	Check_variable()


	println("platform: ", Platform_type)
	println("Local_ip_add:", Local_ip_add)
	println("Agent UniqID:", Klevr_agent_id_string)
	println("Primary:", Primary_ip)

	HandShake()

	if Check_primary() == "true" {
		primScheduler.Every(5).Seconds().Do(printprimary)
		primScheduler.Every(5).Seconds().Do(getCommand)

		go func() {
			<-primScheduler.Start()
		}()
	} else {
		s := gocron.NewScheduler()
		s.Every(5).Seconds().Do(printprimary)

		go func() {
			<-s.Start()
		}()
	}

	http.ListenAndServe(":18800", nil)

}
