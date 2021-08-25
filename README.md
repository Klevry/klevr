![klevr_logo.png](https://raw.githubusercontent.com/Klevry/klevr/master/assets/klevr_logo.png)
<br><a href="https://opensource.klevr.dev">https://opensource.klevr.dev</a>
 * Image click to Youtube:
 * [![Klevr concept diagram](https://user-images.githubusercontent.com/4043594/130567515-05e5b863-9117-4473-907e-89ae9c45229a.png)](https://youtu.be/W3rH7XI-G1A)


# Klevr: Kloud-native everywhere
## Hyper-connected Cloud-native delivery solution for SaaS
 * Asynchronous distributed infrastructure management console and agent for separated networks.
 * Supports for:
   * Baremetal server in the On-premise datacenter
   * PC/Workstation in the Office/intranet
   * Laptop at everywhere
   * Public-cloud

## Kickstart for klevr demo
* docker-compose command
  ```
  git clone https://github.com/klevry/klevr.git
  sudo docker-compose -f klevr/docker-compose.yml up -d
  ```
* Lending page (default ID/PW = admin/admin)
![image](https://user-images.githubusercontent.com/4043594/130207252-da0b5572-9f31-4a04-a55c-1f88b65a6c3d.png)

* Main page
![image](https://user-images.githubusercontent.com/4043594/130207444-8c49a724-fac9-4b5d-8a67-6d63c80fcda5.png)


## Diagram Overview
 * Image click to Youtube:
 * [![Diagram Overview](https://user-images.githubusercontent.com/4043594/130544379-ca032ecb-d1f7-468c-a289-3b3ab3c671d1.png)](https://youtu.be/xLkqm1vEmd0)

## Features
 * **Agent**
   * Provisioning: Docker, Kubernetes, Micro K8s(on Linux laptop) with Vagrant & VirtualBox, Prometheus 
   * Get & Run: Hypervisor(via libvirt container or Multipass), Terraform, Prometheus, Beacon, Helm chart
   * Metric data aggregate & delivery
  * **Manager**
   * Host pool management
   * Resource management
   * Primary host management 
   * Task management(To be)
   * Service catalog management(To be)
   * Service delivery to Dev./Stg./Prod.(To be)
 * **Docker images**
   * [Agent](./Dockerfiles/agent)(user's infrastructure management agent): [klevry/agent:latest](https://hub.docker.com/repository/docker/klevry/klevr-agent)
   * [Manager](./Dockerfiles/manager)(management console): [klevry/manager:latest](https://hub.docker.com/repository/docker/klevry/klevr-manager)
   * ~~[Beacon](./Dockerfiles/beacon)(Primary agent health checker): [klevry/beacon:latest](https://hub.docker.com/repository/docker/klevry/beacon)~~
   * ~~[Libvirt](./Dockerfiles/libvirt)(Hypervisor): [klevry/libvirt:latest](https://hub.docker.com/repository/docker/klevry/libvirt)~~
   * ~~Prometheus operator(Service discovery)~~
   * [ProvBee](https://github.com/NexClipper/provbee)(nexclipper/provbee)  
   * ~~Metric crawler~~
   * Task manager
 * ~~**KV store([Consul](https://github.com/hashicorp/consul))**~~
   
## Simple logic of asynchronous task management - (Click to Youtube for details)
 * [![Primary election of agent](https://raw.githubusercontent.com/Klevry/klevr/master/assets/Klevr_Agent_primary_election_n_delivery_logic.png)](https://www.youtube.com/watch?v=hyMaVsCcgbA&t=2s)



## Architecture
### Scheme of Database
 * AGENT_GROUPS: Zone(A.K.A Group) information of Agents. Task will be seperated by Zone base
 * AGENTS: Manage the status of Agents allowed access to the Manager and information on the Zone to which the Agent belongs
 * API_AUTHENTICATIONS: API key management to go through the authentication process to access the functions provided by the Manager
 * TASK_LOCK: Informs the Manager that it can provide the function of the job by preempting the lock.
 * TASKS: Overall task and status management
 * TASK_DETAIL: Detailed setting contents of each task
 * TASK_STEPS: Manage the steps that perform the actual work of the task
 * TASK_LOGS: Task log

### Structure
 * Klevr has a web-based management tool (console) implemented in React.
   * The user manual of the console can be viewed at [here] (./console/Manual-KR.md).
   * It provides user (admin) authentication and can manage Task, Credential, Zone, Agent, and API Key.
   * By setting "REACT_APP_API_URL" in the ".env" file, you can specify the Manager you want to connect to from the console.
 * Klevr consists of Manager, Agent and DB.
   ![Klevr Elements](https://raw.githubusercontent.com/Klevry/klevr/master/assets/klevr_elements.png)
 * Background tasks to manage tasks and agents in Manager
   * Lock:Check the lock status by scheduler 
     ![background-1](https://raw.githubusercontent.com/Klevry/klevr/master/assets/background_1.png)
   * EventHandler: Notification of task change status with WebHook
     ![background-2](https://raw.githubusercontent.com/Klevry/klevr/master/assets/background_2.png)
   * AgentStatus: Continuously check and change the current status of the Agent
     ![background-3](https://raw.githubusercontent.com/Klevry/klevr/master/assets/background_3.png)
   * ScheduledTask: Task status is Scheduled and the status of tasks before the scheduled time is changed to waitPolling status.
     ![background-4](https://raw.githubusercontent.com/Klevry/klevr/master/assets/background_4.png)
   * TaskHandOverUpdater: DB status change of tasks whose status is HandOver
     ![background-5](https://raw.githubusercontent.com/Klevry/klevr/master/assets/background_5.png)
* Manage Tasks and Agents in Manager
  * Task execution
    ![task execution](https://raw.githubusercontent.com/Klevry/klevr/master/assets/task_execution.png) 
  * Task state transition
    ![task status](https://raw.githubusercontent.com/Klevry/klevr/master/assets/task_status.png)
  * Primary Agent management
    * The Agent who initially requested HandShake from the Manager is selected as Primary
    * After this, Agents requesting HandShake are selected as Secondary
    * Secondary Agents monitor the status of Primary Agent. The first secondary agent that detects an abnormality in the primary agent reports the status of the primary agent to the manager and is selected as the primary agent.


## Usage
### Swagger-UI API
* API Dashboard URL : http://localhost:8090/swagger/index.html
### 1. Zone 
* Create
  * [POST] /inner/groups
* Listing
  * [GET] /inner/groups/{groupID}
* Delete
  * [DELETE] /inner/groups/{groupID}
### 2. API KEY
* Create 
  * [POST] /inner/groups/{groupID}/apikey
* Listing
  * [GET] /inner/groups/{groupID}/apikey
* Modify
  * [PUT] /inner/groups/{groupID}/apikey
### 3. TASK 
* Create
  * [POST] /inner/tasks
* Show
  * [GET] /inner/tasks
* Listing
  * [GET] /inner/tasks/{taskID}
* Cancel
  * [DELETE] /inner/tasks/{taskID}
* Reserved word command information
  * [GET] /inner/commands


## Requirement for use
 * [x] Docker/Docker-compose/Docker-registry
   * [x] ~~Beacon~~
   * [x] ~~Libvirt~~
   * [x] Task manage to [ProvBee](https://github.com/NexClipper/provbee)
 * [x] Terraform of container by [ProvBee](https://github.com/NexClipper/provbee)
 * [x] KVM(libvirt) by [ProvBee](https://github.com/NexClipper/provbee)
 * [x] Multipass for Hosted Virtual-machine
 * [x] ~~Micro K8s~~ K3s
 * [x] Prometheus by [ProvBee](https://github.com/NexClipper/provbee)
 * [x] Grafana by [ProvBee](https://github.com/NexClipper/provbee)
 * [x] Helm by [ProvBee](https://github.com/NexClipper/provbee)
 * [ ] ~~Vault(maybe)~~
 * [ ] ~~Packer(maybe)~~
 * [x] ~~Vagrant~~
 * [x] ~~Consul~~ 

## Description for Directories and files
```
.
├── README.md                   // This Screen as you see. :)
├── docker-compose.yml          // Kickstarter: Bootstraping by docker-compose
├── Dockerfiles                  // Directory for docker image build
│   ├── libvirt
│   └── manager                 // Actual binary file of manager will be move to this link directory for the docker build
├── assets
│   └── [Images & Contents]
├── cmd                         // Actual artifacts fpr Klevr agent & manager(webserver) 
│   ├── klevr-agent
│   │   ├── Makefile
│   │   ├── agent_installer.sh  // Remote installer via curl command as a generated script by Manager
│   │   ├── klevr               // Actual `Klevr` agent binary
│   │   └── main.go             // main source code of the Agent
│   └── klevr-manager
│       ├── Docker -> ../../Dockerfiles/manager  // Binary artifact send to this directory for Docker build  
│       ├── Makefile
│       └── main.go             // main source code of the Manager
├── conf
│   ├── klevr-manager-db.sql.create        // Database for Manager initializing & running
│   ├── klevr-manager-db.sql.modify        // Database for 
│   └── klevr-manager-local.yml // Config file for Manager running
├── pkg
│   ├── common                  // 'common' package directory
│   │   ├── config.go
│   │   ├── error.go
│   │   ├── http.go
│   │   ├── log.go
│   │   └── orm.go
│   ├── communicator            // 'communicator' package directory
│   │   ├── README.md
│   │   └── communicator.go
│   └── manager                 // 'manager' package directory
│       ├── api.go
│       ├── api_agent.go
│       ├── api_install.go
│       ├── api_legacy.go
│       ├── api_model.go
│       ├── handler.go
│       ├── persist_model.go
│       ├── repository.go
│       └── server.go
├── go.mod
├── go.sum
└── scripts                    // Operation script for Provisioning
    └── [Provisioning scripts]

```
