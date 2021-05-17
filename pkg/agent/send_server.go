package agent

import (
	"context"
	"encoding/json"
	"runtime"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/mackerelio/go-osstat/memory"

	pb "github.com/Klevry/klevr/pkg/agent/protobuf"
	"github.com/NexClipper/logger"
)

type sendServer struct {
	agentKey string
}

func (s sendServer) SendTask(ctx context.Context, in *pb.Message) (*pb.Message, error) {
	logger.Debugf("Receive message body from client: %v", string(in.Task))

	var t common.KlevrTask

	err := json.Unmarshal(in.Task, &t)
	if err != nil {
		logger.Debugf("%v", string(in.Task))
		logger.Error(err)
	}

	executor := common.GetTaskExecutor()
	executor.RunTask(s.agentKey, &t)

	result, _ := executor.GetUpdatedTasks()

	b := jsonMarshal(result)

	return &pb.Message{Task: b}, nil
}

func (s sendServer) StatusCheck(ctx context.Context, in *pb.Status) (*pb.Status, error) {
	logger.Debugf("Receive message body from client: %v", in.Status)

	disk := diskUsage("/")
	memory, _ := memory.Get()

	agentStatus := &common.AgentStatus{
		AgentKey: s.agentKey,
		Resource: &common.Resource{
			Core:       runtime.NumCPU(),
			Memory:     int(memory.Total / MB),
			Disk:       int(disk.All / MB),
			FreeMemory: int(memory.Free / MB),
			FreeDisk:   int(disk.Free / MB),
		},
	}

	b := jsonMarshal(agentStatus)

	return &pb.Status{Status: b}, nil
}
