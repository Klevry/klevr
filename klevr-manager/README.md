# Klevr Manager Server
## Environment set-up
### Local
* Create Database
```
docker run --name mariadb -v [로컬 데이터베이스 경로]:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=[root 패스워드] -p [expose port]:3306 -d mariadb
```
* Initialize Database
```
docker exec -i mariadb sh -c 'exec mysql -uroot -p"[root 패스워드]"' < [프로젝트 root path]/sql/klevr-init.sql
```
* Create Redis
```
docker run --name redis -p 6379:6379 -d redis
```
## How to use
* Simple for localhost serving: ```./klevr_webconsole```
* Set for custom: ```klevr_webconsole -port=8080 -apiserver="http://192.168.2.100:8500"```
 * apiserver is consul address

## APIs
* Link for Agent Download: `[Web-console URL]/`
* Show hosts info.: `[Web-console URL]/user/[USERID]/hostsinfo` with Master of agent
* Purge old host: `[Web-console URL]/user/[USERID]/hostsmgt`
* Client info. receiver: `[Web-console URL]/user/[USERID]/hostname/{HOSTNAME}/[IP]/type/[baremetal/aws]/[TTL]/[MASTER_STATUS]`
* Host system info. receiver: `[Web-console URL]/user/[USERID]/hostname/{HOSTNAME}/hostinfo`
