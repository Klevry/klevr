# Klevr Webserver console
## How to use
* Simple for localhost serving: ```./klevr_webconsole```
* Set for custom: ```klevr_webconsole -port=8080 -apiserver="http://192.168.2.100:8500"```
 * apiserver is consul address

## APIs
* Link for Agent Download: `[Web-console URL]/`
* Show API server info.: `[Web-console URL]/apiserver`
* Show hosts info.: `[Web-console URL]/user/[USERID]/hostsinfo` with Master of agent
* Purge old host: `[Web-console URL]/user/[USERID]/hostmgt`
