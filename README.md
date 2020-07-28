![klevr_logo.png](https://raw.githubusercontent.com/Klevry/klevr/master/assets/klevr_logo.png)
# Klevr: Kloud-native everywhere
## Interconnector for the Platform based SaaS delivery
 * Asynchronous distributed infrastructure management console and agent for separated networks.
 * Supports for:
   * Baremetal server in the On-premise datacenter
   * PC/Workstation in the Office/intranet
   * Laptop at everywhere
   * Public-cloud

## Kickstart for webconsole & KVstore
* docker-compose command
```
git clone https://github.com/ralfyang/klevr.git
docker-compose up -d
```

## Diagram Overview
 * Image click to Youtube:
 * [![Diagram Overview](https://raw.githubusercontent.com/Klevry/klevr/master/assets/Klevr_diagram_overview.png)](https://youtu.be/xLkqm1vEmd0)

## Features
 * **[Agent](./agent/)**
   * Provisioning: Docker, Kubernetes, Micro K8s(on Linux laptop) with Vagrant & VirtualBox, Prometheus 
   * Get & Run: Hypervisor(via libvirt container), Terraform, Prometheus, Beacon
   * Metric data aggregate & delivery
  * **[Web console](./webconsole/)**
   * Host pool management
   * Resource management
   * Primary host management 
   * Task management 
   * Service catalog management
   * Service delivery to Dev./Stg./Prod.
 * **Docker images**
   * [Webconsole](./Dockerfile/klevr_websonsole)(Webserver): [klevry:webconsole:latest](https://hub.docker.com/repository/docker/klevry/webconsole)
   * ~~[Beacon](./Dockerfile/beacon)(Primary agent health checker): [klevry/beacon:latest](https://hub.docker.com/repository/docker/klevry/beacon)~~
   * [Libvirt](./Dockerfile/libvirt)(Hypervisor): [klevry/libvirt:latest](https://hub.docker.com/repository/docker/klevry/libvirt)
   * Prometheus(Container monitoring)
   * Metric crawler
   * Task manager
 * **KV store([Consul](https://github.com/hashicorp/consul))**
   

## Requirement for use
 * [ ] Docker/Docker-compose/Docker-registry
   * [x] ~~Beacon~~
   * [x] Libvirt
   * [ ] Task manager to terraform
 * [ ] Terraform of container
 * [x] KVM(libvirt)
 * [ ] ~~Micro K8s~~
 * [ ] K3s
 * [x] ~~Consul~~
 * [ ] Prometheus 
 * [x] ~~Vagrant~~
 * [ ] Halm
 * [ ] Vault(maybe)
 * [ ] Packer(maybe)


## Description for Directories and files
```
.
├── README.md                   // This Screen as you see. :)
├── docker-compose.yml          // Kickstarter: Bootstraping by docker-compose
├── Dockerfile                  // Directory for docker image build
│   ├── libvirt
│   └── manager                 // Actual binary file of manager will be move to this linke directory for the docker build
├── assets
│   └── [Images & Contents]
├── cmd                         // Actual artifacts fpr Klevr agent & manager(webserver) 
│   ├── klevr-agent
│   │   ├── Makefile
│   │   ├── agent_installer.sh  // Remote installer via curl command as a generated script by Manger
│   │   ├── klevr               // Actual `Klevr` agent binary
│   │   └── main.go             // main source code of the Agent
│   └── klevr-manager
│       ├── Docker -> ../../Dockerfile/manager  // Binary artifact send to this directory for Docker build  
│       ├── Makefile
│       └── main.go             // main source code of the Manager
├── conf
│   ├── Dump20200720.sql        // Database for Manager initialinzing & running
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
