# klevr-agent
## Introduction
Klevr manager와 통신하며 manager로부터 task를 전달받아 수행하고 그 결과 에러를 매니저에게 전달한다.

프라이머리와 세컨더리는 매니저가 정해준다.

세컨더리는 프라이머리의 동작여부를 확인한다.

프라이머리가 동적하지않으면 세컨더리는 매니저에게 그것을 알리고 알린 세컨더리가 프라이머리가 된다.

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


## Process
1. 에이전트 구동시 매니저에게 handshake를 보낸다. 처음 handshake를 받은 에이전트가 primary가 된다.
2. primary는 매니저에게 command를 받아와 provbee에게 전달한다.
3. 설치된 grafana의 url을 manager로 전송한다.

## How to use
- Configuration

| parameters | description | default |
| --- | ---- | --- | 
| apikey | Account ID from Klevr service |  | 
| platform | [baremetal, k8s, aws] - Service Platform for Host build up |  |
| zoneId | Zone ID from Klevr service |  | 
| manager | IP address from Klevr manager |  | 
