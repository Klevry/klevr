
![klevr_logo.png](klevr_logo.png)
# Klevr: K(C)loud-native everywhere cleverly
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
 * [![Diagram Overview](/Klevr_diagram_overview.png)](https://www.youtube.com/watch?v=3dhf-Pzc13Y)


## Features
 * **Agent([./agent/klevr_agent](agent/))**
   * Provisioning: Docker, Micro K8s, Vagrant, VirtualBox
   * Get & Run: Hypervisor(via libvirt container), Terraform, Prometheus, Beacon
   * Metric data aggregate & delivery
  * **Web console([./webconsole/klevr_webconsole](./webconsole/))**
   * Host pool management
   * Resource management
   * Master node management 
   * Task management 
   * Service catalog management
   * Service delivery to Dev./Stg./Prod.
 * **Docker images**
   * Klever_console(Webserver)
   * Beacon(master health checker)
   * Libvirt(Hypervisor)
   * Prometheus(Container monitoring)
   * Metric crawler
   * Task manager
 * **KV store([Consul](https://github.com/hashicorp/consul))**
   

## Requirement for use
 * [ ] Docker/Docker-compose/Docker-registry
   * [x] Beacon: https://hub.docker.com/repository/docker/klevry/beacon
   * [x] Libvirt: https://hub.docker.com/repository/docker/klevry/libvirt
   * [ ] Task manager
 * [ ] Terraform 
 * [x] KVM(libvirt)
 * [ ] Micro K8s
 * [x] Consul
 * [ ] Prometheus 
 * [ ] Vagrant
 * [ ] Halm
 * [ ] Vault(maybe)
 * [ ] Packer(maybe)


