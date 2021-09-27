package agent

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"google.golang.org/grpc"

	pb "github.com/Klevry/klevr/pkg/agent/protobuf"
)

func (agent *KlevrAgent) checkPrimary(prim string) bool {
	if prim == localIPAddress(agent.NetworkInterfaceName) {
		logger.Debug("I am Primary")

		return true
	} else {
		logger.Debug("I am Secondary")

		return false
	}
}

/*func (agent *KlevrAgent) primaryTaskSend(ip string, task []byte) *common.KlevrTask {
	serverAddr := net.JoinHostPort(ip, agent.grpcPort)
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		logger.Errorf("did not connect :%v", err)
		return nil
	}

	defer conn.Close()

	state := conn.GetState()
	logger.Debug(state.String())
	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	c := pb.NewTaskSendClient(conn)

	// send to secondary
	r, resErr := c.SendTask(ctx, &pb.Message{Task: task})
	if resErr != nil {
		logger.Errorf("could not response: %v", resErr)
		return nil
	}

	logger.Debugf("this is response: %v", r)

	var t common.KlevrTask

	err = json.Unmarshal(r.Task, &t)
	if err != nil {
		logger.Debugf("%v", string(r.Task))
		logger.Error(err)
	}

	return &t
}*/

// ZoneStatusCheck는 현재 소속된 zone의 agent들의 상태 정보를 확인
func (agent *KlevrAgent) checkZoneStatus() {
	for i, n := range agent.Agents {
		if (agent.Primary.IP == n.IP) && (agent.Primary.AgentKey == n.AgentKey) {
			agent.Agents[i].LastAliveCheckTime = &common.JSONTime{Time: time.Now().UTC()}
			agent.Agents[i].IsActive = true
		} else {
			serverAddr := net.JoinHostPort(n.IP, agent.grpcPort)
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				logger.Errorf("did not connect :%v", err)
			}

			ctx, _ := context.WithTimeout(context.Background(), time.Second)
			c := pb.NewTaskSendClient(conn)

			s, resErr := c.StatusCheck(ctx, &pb.Status{})
			if resErr == nil {
				var agentStatus common.AgentStatus
				json.Unmarshal(s.Status, &agentStatus)

				logger.Debugf("AgentStatus: %v", agentStatus)

				if n.AgentKey == agentStatus.AgentKey {
					agent.Agents[i].Core = agentStatus.Core
					agent.Agents[i].Memory = agentStatus.Memory
					agent.Agents[i].Disk = agentStatus.Disk
					agent.Agents[i].FreeMemory = agentStatus.FreeMemory
					agent.Agents[i].FreeDisk = agentStatus.FreeDisk
					agent.Agents[i].LastAliveCheckTime = &common.JSONTime{Time: time.Now().UTC()}
					agent.Agents[i].IsActive = true
				} else {
					agent.Agents[i].IsActive = false
				}
			} else {
				logger.Debugf("zoneStatusCheck error: %v", resErr)
				agent.Agents[i].IsActive = false
			}

			conn.Close()
		}
	}
}
