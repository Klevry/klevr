# klevr
* Asynchronous distributed infrastructure management console and agent for separated network.
* Supports for:
  * Baremetal server in the On-premise datacenter
  * PC/Workstation in the Office/intranet
  * Laptop at everywhere
  * Public-cloud

## Diagram Overview
 * Image click to Youtube:
 * [![Diagram Overview](/Klevr_diagram_overview.png)](https://www.youtube.com/watch?v=3dhf-Pzc13Y)


## Architecture
 * Agent([./agent/klevr_agent](agent/)) - API([Consul](https://github.com/hashicorp/consul)) - Web Server([./webconsole/klevr_webconsole](./webconsole/))

## Requirement for use
 * Docker/Docker-registry
 * Terraform 
 * libvirt container
 * Micro K8s
 * Consul
 * golang
 * Prometheus 
 * Vagrant
 * Vault(maybe)
 * Packer(maybe)
