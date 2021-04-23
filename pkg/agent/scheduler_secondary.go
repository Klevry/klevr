package agent

import (
	"net"
	"strconv"
	"time"

	"github.com/NexClipper/logger"
	"google.golang.org/grpc"

	pb "github.com/Klevry/klevr/pkg/agent/protobuf"
)

func (agent *KlevrAgent) secondaryServer() {
	logger.Debugf("GRPC SERVER START!!!!")

	var errLis error

	addressAndPort := net.JoinHostPort(localIPAddress(agent.NetworkInterfaceName), agent.grpcPort)
	_, err := net.DialTimeout("tcp", addressAndPort, time.Second)
	if err != nil {
		logger.Errorf("not open port!@#!@#!@@#")

		// grpc server start
		if agent.NetworkInterfaceName == "" {
			agent.connect, errLis = net.Listen("tcp", ":"+agent.grpcPort)
		} else {
			agent.connect, errLis = net.Listen("tcp", addressAndPort)
		}
		if errLis != nil {
			logger.Fatalf("failed to liesten: %v", err)
		}
	}

	grpcServer := grpc.NewServer()

	pb.RegisterTaskSendServer(grpcServer, sendServer{agentKey: agent.AgentKey})

	if err := grpcServer.Serve(agent.connect); err != nil {
		logger.Fatalf("failed to serve: %s", err)
	}
}

func statusCheck(agent *KlevrAgent) {
	_, err := net.DialTimeout("tcp", agent.Primary.IP+":"+strconv.Itoa(agent.Primary.Port), 3*time.Second)
	if err != nil {
		logger.Errorf("%v", err)
		agent.primaryStatusReport()
	}
}
