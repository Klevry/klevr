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


## Features
 * Agent([./agent/klevr_agent](agent/))
   * Provisioning: Docker, Micro K8s
   * Get & Run: Hypervisor(via libvirt container), Terraform, Prometheus, Vagrant
   * Metric data aggregate & delivery
 * API([Consul](https://github.com/hashicorp/consul))
 * Web console([./webconsole/klevr_webconsole](./webconsole/))
   * Host pool management
   * Resource management
   * Master node management 
   * Task management 
   * Service catalog management
   * Service delivery to Dev./Stg./Prod.
   

## Requirement for use
 * Docker/Docker-registry
 * Terraform 
 * libvirt container
 * Micro K8s
 * Consul
 * Prometheus 
 * Vagrant
 * Halm
 * Vault(maybe)
 * Packer(maybe)
