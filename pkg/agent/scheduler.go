package agent

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/model"
	"github.com/Klevry/klevr/pkg/serialize"
	"github.com/NexClipper/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"

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

// ZoneStatusCheck는 현재 소속된 zone의 agent들의 상태 정보를 확인
func (agent *KlevrAgent) checkZoneStatus() {
	for i, n := range agent.Agents {
		if (agent.Primary.IP == n.IP) && (agent.Primary.AgentKey == n.AgentKey) {
			agent.Agents[i].LastAliveCheckTime = &serialize.JSONTime{Time: time.Now().UTC()}
			agent.Agents[i].IsActive = true
		} else {
			serverAddr := net.JoinHostPort(n.IP, agent.grpcPort)
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				logger.Errorf("did not connect :%v", err)
			}

			state := conn.GetState()

			if err != nil || !(state == connectivity.Ready || state == connectivity.Idle) {
				logger.Debugf("zoneStatusCheck dial error(%v): %s", err, state.String())
				agent.Agents[i].IsActive = false
			} else {
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
						agent.Agents[i].LastAliveCheckTime = &serialize.JSONTime{Time: time.Now().UTC()}
						agent.Agents[i].IsActive = true
					} else {
						agent.Agents[i].IsActive = false
					}
				} else {
					logger.Debugf("zoneStatusCheck error: %v", resErr)
					agent.Agents[i].IsActive = false
				}
			}

			conn.Close()
		}
	}
}

func (agent *KlevrAgent) getRemoteUpdatedTasks() []model.KlevrTask {
	remoteTasks := make([]model.KlevrTask, 0)

	for _, n := range agent.Agents {
		if !((agent.Primary.IP == n.IP) && (agent.Primary.AgentKey == n.AgentKey)) {
			serverAddr := net.JoinHostPort(n.IP, agent.grpcPort)
			conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
			if err != nil {
				logger.Errorf("did not connect :%v", err)
			}

			state := conn.GetState()

			if err != nil || !(state == connectivity.Ready || state == connectivity.Idle) {
				logger.Debugf("getRemoteUpdatedTasks dial error(%v): %s", err, state.String())
			} else {
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				c := pb.NewTaskSendClient(conn)

				s, resErr := c.GetUpdatedTasks(ctx, &pb.Message{})
				if resErr == nil {
					tasks := make([]model.KlevrTask, 0)
					json.Unmarshal(s.Task, &tasks)

					logger.Debugf("tasks: %v", tasks)

					for _, item := range tasks {
						remoteTasks = append(remoteTasks, item)
					}
				} else {
					logger.Debugf("getRemoteUpdatedTasks error: %v", resErr)
				}
			}

			conn.Close()
		}
	}

	return remoteTasks
}
