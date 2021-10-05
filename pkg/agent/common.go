package agent

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/NexClipper/logger"
	"github.com/denisbrodbeck/machineid"
	"github.com/shirou/gopsutil/disk"
	netutil "k8s.io/apimachinery/pkg/util/net"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

// get local ip address
func getIPAddress() string {
	// get Local IP address automatically
	default_ip, err := netutil.ChooseHostInterface()
	if err != nil {
		log.Fatalf("Failed to get IP address: %v", err)
	}

	return default_ip.String()
}

func localIPAddress(networkInterfaceName string) string {
	var ipAddress string
	if networkInterfaceName == "" {
		ipAddress = getIPAddress()
	} else {
		nifs, err := net.Interfaces()
		if err != nil {
			log.Fatalf("Failed to get Network Interfaces: %v", err)
		}

		for _, ni := range nifs {
			if ni.Name == networkInterfaceName {
				addrs, err := ni.Addrs()
				if err != nil {
					log.Fatalf("Failed to get Address: %v", err)
				}

				for _, a := range addrs {
					v := a.(*net.IPNet)
					if ipv4addr := v.IP.To4(); ipv4addr != nil {
						ipAddress = ipv4addr.String()
						break
					}
				}
			}
		}
	}
	if ipAddress == "" {
		log.Fatalf("Failed to get IP address")
	}

	return ipAddress
}

// disk usage of path/disk
func diskUsage(path string) (d DiskStatus) {
	u, err := disk.Usage(path)
	if err != nil {
		return
	}

	d.All = u.Total
	d.Free = u.Free
	d.Used = u.Used

	return
}

func readFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error(err)
	}

	return data
}

// generate agent key
func agentKeyGen() (string, error) {

	/*
		nowTime := strconv.FormatInt(time.Now().UTC().Unix(), 10)

		uuid := make([]byte, 16)
		n, err := io.ReadFull(rand.Reader, uuid)
		if n != len(uuid) || err != nil {
			return "", err
		}
		// variant bits; see section 4.1.1
		uuid[8] = uuid[8]&^0xc0 | 0x80
		// version 4 (pseudo-random); see section 4.1.3
		uuid[6] = uuid[6]&^0xf0 | 0x40

		key := hex.EncodeToString(uuid) + nowTime
	*/

	// machineid.ID(): key => len(32)
	// machineid.ProtectedID(): key => len(64)
	key, err := machineid.ID()
	if err != nil {
		return "", err
	}

	return key, nil
}

func checkAgentKey() string {
	var agentIdFile = "/tmp/klevr_agent.id"
	var agentIdString string

	//if agent file exist
	if _, err := os.Stat(agentIdFile); !os.IsNotExist(err) {
		data := readFile(agentIdFile)

		if string(data) != "" {
			agentIdString = string(data)
		} else {
			logger.Error("There is no agent ID")
		}

	} else {
		key, err := agentKeyGen()
		if err != nil {
			logger.Error(err)
		}

		err = ioutil.WriteFile(agentIdFile, []byte(key), os.FileMode(0644))
		if err != nil {
			logger.Error(err)
		}

		agentIdString = key
	}

	return agentIdString
}
