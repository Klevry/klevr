# klevr-agent
## Introduction
Klevr manager, Provbee와 통신을 하며 할당된 작업을 수행하는 Klevr agent

## Process
1. 에이전트 구동시 매니저에게 handshake를 보낸다. 
1. 처음 handshake를 받은 에이전트가 primary가 된다.
1. primary는 매니저에게 command를 받아와 provbee에게 전달한다.
1. provbee는 명령을 수행하고 결과를 primary에게 전달한다.
1. primary agent는 전달받은 결과값을 manager에게 전달한다

> 현재는 https://github.com/NexClipper/klevr_tasks/blob/master/queue 에서 명령어를 가져와서 실행하는 방식으로 개발되어있다.

### pkg
```shell script
.
├── agent_info.go             // Get agents info
├── agents.go                 // Agent main
├── file_rwd.go               // File read, write, delete functions
├── get_init_command.go       // Get tasks commands
├── handshake.go              // Handshake to manager
├── init_agent.go             // Create agentkey, initialize primary
├── primary_init_callback.go  
├── primary_status_report.go
└── task_mgmt.go
```

## How to use
- Configuration

| parameters | description | default |
| --- | ---- | --- | 
| apikey | Account ID from Klevr service |  | 
| platform | [baremetal, k8s, aws] - Service Platform for Host build up |  |
| zoneId | Zone ID from Klevr service |  | 
| manager | IP address from Klevr manager |  | 


- Useage
```shell script
go run ./main.go -apiKey="" -platform="" -manager="" -zoneId=""
```