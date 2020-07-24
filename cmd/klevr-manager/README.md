# Klevr manager(web console)
## API
### API list
|API| URI |method |parameters | callback | description | example |

| agent setup script 생성 | /install/agents/bootstrap | POST | platform (string) zoneId (uint64) apiKey (string) managerUri(uri) | 설치 스크립트 curl -sL https:/asdfasdfasdf.com/file/agents/download | bash -c -zoneid -asdkljlasdkf | github에서 설치 스크립트 다운로드 | file download url 클릭 사용자에 맞게 제네레이티드된 주소값이 생성 복사 붙여넣기( curl -sL https:/asdfasdfasdf.com/file/agents/download?platform=baremetal&zoneId=1234&apiKey=4rjuifdhj93rnfkl) ->  curl -sL ljkasjdlfjas.com/down | bash -c -zoneid -asdkljlasdkf 
| agent 다운로드 | /install/agents/download | GET|| agent 바이너리 |github에서 설치 스크립트 다운로드 || 
|agent hand-shake | /agents/handshake | PUT |platform (string), zoneId (uint64), apiKey (string), agentKey (string), | agent 실행 정보 - 암호화키 - 로그 레벨 - 호출 주기 primary 정보 - primary 노드 정보 |최초 agent 기동 시 호출 |agentKey = zoneid(int64)+hash(MID+IP) 
|pooling & task mgmt |/agents/{agentKey} |PUT |hosts-ip, hosts-alivecheck hosts-resource(cpu,mem,disk) |task 정보, target secondary, agent, hosts list||| 
|primary status report| /agents/reports/{agentKey}| PUT |primary ip |primary 정보 - primary 노드 정보|||


