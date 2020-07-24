# Klevr Webserver console
## How to use
* Simple for localhost serving: ```./klevr_webconsole```
* Set for custom: ```klevr_webconsole -port=8080 -apiserver="http://192.168.2.100:8500"```
 * apiserver is consul address

## Platform type Example
 * Baremetal
 * Linux laptop(micro k8s)
 * Kubernetes
 * Prometheus
 * AWS/Azure/GCP


## APIs
* Link for Agent Download: `[Web-console URL]/`
* Show hosts info.: `[Web-console URL]/user/[USERID]/hostsinfo` with Primary of agent
* Purge old host: `[Web-console URL]/user/[USERID]/hostsmgt`
* Client info. receiver: `[Web-console URL]/user/[USERID]/hostname/{HOSTNAME}/[IP]/type/[baremetal/aws]/[TTL]/[PRIMARY_STATUS]`
* Host system info. receiver: `[Web-console URL]/user/[USERID]/hostname/{HOSTNAME}/hostinfo`
