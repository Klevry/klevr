package agent

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"runtime"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/shirou/gopsutil/disk"
	netutil "k8s.io/apimachinery/pkg/util/net"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// get local ip address
func Local_ip_add() string {
	// get Local IP address automatically
	default_ip, err := netutil.ChooseHostInterface()
	if err != nil {
		log.Fatalf("Failed to get IP address: %v", err)
	}

	return default_ip.String()
}

func LocalIPAddress(networkInterfaceName string) string {
	var ipAddress string
	if networkInterfaceName == "" {
		ipAddress = Local_ip_add()
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
func DiskUsage(path string) (d DiskStatus) {
	u, err := disk.Usage(path)
	if err != nil {
		return
	}

	d.All = u.Total
	d.Free = u.Free
	d.Used = u.Used

	return
}

// send agent info to manager
func (agent *KlevrAgent) SendMe(body *common.Body) {
	body.Me.IP = LocalIPAddress(agent.NetworkInterfaceName)
	body.Me.Port = 18800
	body.Me.Version = "0.1.0"

	disk := DiskUsage("/")

	memory, _ := memory.Get()

	body.Me.Resource.Core = runtime.NumCPU()
	body.Me.Resource.Memory = int(memory.Total / MB)
	body.Me.Resource.Disk = int(disk.All / MB)
}

func JsonMarshal(a interface{}) []byte {
	b, err := json.Marshal(a)
	if err != nil {
		logger.Debugf("%v", string(b))
		logger.Error(err)
	}

	return b
}

func JsonUnmarshal(a []byte) (*common.Body, error) {
	var body common.Body

	err := json.Unmarshal(a, &body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

func ReadFile(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Error(err)
	}

	return data
}
