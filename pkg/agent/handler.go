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

const (
	port = "9350"
)

type server struct{}

func (s *server) SendTask(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	logger.Debugf("Receive message body from client: %v", string(in.Task))

	var t common.KlevrTask

	err := json.Unmarshal(in.Task, &t)
	if err != nil {
		logger.Debugf("%v", string(in.Task))
		logger.Error(err)
	}

	executor.RunTask(&t)

	result, _ := executor.GetUpdatedTasks()

	b := JsonMarshal(result)

	return &pb.Message{Task: b}, nil
}

func (s *server) StatusCheck(ctx context.Context, in *pb.Status) (*pb.Status, error) {
	logger.Debugf("Receive message body from client: %v", in.Status)

	return &pb.Status{Status: "OK"}, nil
}

func (agent *KlevrAgent) SecondaryServer() {
	logger.Debugf("GRPC SERVER START!!!!")

	var errLis error

	addressAndPort := net.JoinHostPort(LocalIPAddress(agent.NetworkInterfaceName), port)
	_, err := net.DialTimeout("tcp", addressAndPort, time.Second)
	if err != nil {
		logger.Errorf("not open port!@#!@#!@@#")

		// grpc server start
		if agent.NetworkInterfaceName == "" {
			agent.connect, errLis = net.Listen("tcp", ":"+port)
		} else {
			agent.connect, errLis = net.Listen("tcp", addressAndPort)
		}
		if errLis != nil {
			logger.Fatalf("failed to liesten: %v", err)
		}
	}

	grpcServer := grpc.NewServer()

	pb.RegisterTaskSendServer(grpcServer, &server{})

	if err := grpcServer.Serve(agent.connect); err != nil {
		logger.Fatalf("failed to serve: %s", err)
	}
}

func (agent *KlevrAgent) PrimaryTaskSend(ip string, task []byte) {
	serverAddr := ip + port
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		logger.Errorf("did not connect :%v", err)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second)
	c := pb.NewTaskSendClient(conn)

	// send to secondary
	r, resErr := c.SendTask(ctx, &pb.Message{Task: task})
	if resErr != nil {
		logger.Errorf("could not response: %v", resErr)
	}

	logger.Debugf("this is response: %v", r)
}
