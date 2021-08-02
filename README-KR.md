![klevr_logo.png](https://raw.githubusercontent.com/Klevry/klevr/master/assets/klevr_logo.png)

# Klevr: Kloud-native everywhere
 * <a href="https://opensource.klevr.dev">https://opensource.klevr.dev</a>

## 플랫폼 기반 SaaS 전달을 위한 인터커넥터
 * 분리된 네트워크를 위한 비동기 분산 인프라 관리 콘솔 및 에이전트.
 * 지원:
   * 온-프레미스 데이터센터의 베어메탈 서버
   * 회사/인트라넷의 PC/워크스테이션
   * 노트북
   * 퍼블릭 클라우드

## Klevr 데모를 위한 시작하기
* docker-compose 명령어
  ```
  git clone https://github.com/klevry/klevr.git
  docker-compose -f docker-compose-demo.yml up -d
  ```

## 다이어그램 개요
 * 이미지를 클릭 시 유튜브로 이동:
   [![Diagram Overview](https://raw.githubusercontent.com/Klevry/klevr/master/assets/Klevr_diagram_overview.png)](https://youtu.be/xLkqm1vEmd0)

## 특징
 * **[Agent](./agent/)**
   * 프로비저닝: Docker, Kubernetes, Micro K8s(on Linux laptop) with Vagrant & VirtualBox, Prometheus 
   * 다운로드 및 실행: Hypervisor(via libvirt container), Terraform, Prometheus, Beacon
   * 메트릭 데이터 집계 및 전달
 * **[Manager](./manager)**
   * **[Web console](./console/)**
   * 호스트 풀 관리
   * 리소스 관리
   * 프라이머리 호스트 관리 
   * 작업 관리 
   * 서비스 카탈로그 관리
   * 개발/스테이징/프로덕션에 서비스 전달
 * **도커 이미지**
   * [Agent](./Dockerfile/agent)(user's infrastructure management agent): [klevry/agent:latest](https://hub.docker.com/repository/docker/klevry/klevr-agent)
   * [Manager](./Dockerfile/manager)(management console): [klevry/manager:latest](https://hub.docker.com/repository/docker/klevry/klevr-manager)
   * [Console](./Dockerfile/console)(web console): [klevry/console:latest](https://hub.docker.com/repository/docker/klevry/klevr-console)
## 아키텍쳐
### 데이터베이스 주요 스키마
 * AGENT_GROUPS: Agent들이 속해 있는 Zone(Group)의 정보. Zone(Group)을 기준으로 Task들이 Agent에 분배될 수 있습니다.
 * AGENTS: Manager에 접근이 허가된 Agents들의 상태 및 해당 Agent가 속해 있는 Zone의 정보등을 관리
 * API_AUTHENTICATIONS: Manager에서 제공하는 기능을 사용할 수 있는 agent들의 인증용 API Key를 관리
 * TASK_LOCK: Manager가 Task의 기능을 제공할 수 있음을 Lock을 선점함으로 알림
 * TASKS: Task의 전반적인 사항과 상태 관리
 * TASK_DETAIL: 각 Task의 세부 설정 내용
 * TASK_STEPS: Task의 실제 작업을 수행하는 Step을 관리
 * TASK_LOGS: Task의 로그

### 구조
 * Klevr는 React로 구현된 Web 기반의 관리 도구(console)을 갖고 있습니다.
   * Console의 사용자 메뉴얼은 [여기](./console/Manual-KR.md)에서 볼 수 있습니다.
   * 사용자(admin) 인증을 제공하고 있으며, Task, Credential, Zone, Agent 그리고 API Key를 관리 할 수 있습니다.
   * ".env" 파일에서 "REACT_APP_API_URL"를 설정함으로써 Console에서 연결되고자 하는 Manager를 지정할 수 있습니다.
 * Klevr는 Manager, Agent 그리고 DB로 구성되어 있습니다  
   ![Klevr Elements](https://raw.githubusercontent.com/Klevry/klevr/master/assets/klevr_elements.png)
 * Manager에서 Task와 Agent를 관리하기 위한 백그라운드 작업들
   * Lock: Lock 상태를 주기별로 확인  
     ![background-1](https://raw.githubusercontent.com/Klevry/klevr/master/assets/background_1.png)
   * EventHandler: WebHook으로 Task의 변경 상태에 대해서 알림  
     ![background-2](https://raw.githubusercontent.com/Klevry/klevr/master/assets/background_2.png)
   * AgentStatus: Agent의 현재 상태를 지속적으로 확인 변경  
     ![background-3](https://raw.githubusercontent.com/Klevry/klevr/master/assets/background_3.png)
   * ScheduledTask: Task 상태가 Scheduled 이고 예약 시간보다 이전인 Task들의 상태를 waitPolling 상태로 변경  
     ![background-4](https://raw.githubusercontent.com/Klevry/klevr/master/assets/background_4.png)
   * TaskHandOverUpdater: 상태가 HandOver인 Task들의 DB 상태 변경  
     ![background-5](https://raw.githubusercontent.com/Klevry/klevr/master/assets/background_5.png)
* Manager에서 Task 및 Agent 관리
  * Task 실행
    ![task execution](https://raw.githubusercontent.com/Klevry/klevr/master/assets/task_execution.png) 
  * Task 상태 변환
    ![task status](https://raw.githubusercontent.com/Klevry/klevr/master/assets/task_status.png)
  * Primary Agent 관리
    * 최초에Manager에게 HandShake를 요청한 Agent가 Primary로 선정
    * 이 후 HandShake 요청하는 Agent들은 Secondary로 선정
    * Secondary Agent들은 Primary Agent의 상태를 감시합니다. Primary Agent의 이상을 감지한 최초의 Secondary Agent가 Manager에게 Primary Agent의 상태를 보고 후 Primary Agent로 선정됩니다.

## 비동기 작업 관리의 간단 로직 - (클릭 시 유튜브로 이동합니다.)
 * [![Primary election of agent](https://raw.githubusercontent.com/Klevry/klevr/master/assets/Klevr_Agent_primary_election_n_delivery_logic.png)](https://www.youtube.com/watch?v=hyMaVsCcgbA&t=2s)


## 사용을 위해 필요한 것
 * [x] Docker/Docker-compose/Docker-registry
   * [x] ~~Beacon~~
   * [x] Libvirt
   * [x] Task manage to [ProvBee](https://github.com/NexClipper/provbee)
 * [x] Terraform of container by [ProvBee](https://github.com/NexClipper/provbee)
 * [x] KVM(libvirt)
 * [x] ~~Micro K8s~~ K3s
 * [x] Prometheus 
 * [x] Grafana
 * [ ] Helm
 * [ ] Vault(maybe)
 * [ ] ~~Packer(maybe)~~
 * [x] ~~Vagrant~~
 * [x] ~~Consul~~ 


## 디렉토리와 파일 설명 
```
.
├── CNAME
├── Dockerfile
│   ├── README.md
│   ├── agent
│   ├── console
│   ├── libvirt
│   └── manager
├── LICENSE
├── README.md                          // 지금 보고있는 화면입니다 :)
├── assets
│   ├── [Images & Contents]
├── build.sh
├── cmd                                // 실제 아티팩트 fpr Klevr 에이전트 및 관리자 (웹 서버)
│   ├── klevr-agent
│   │   ├── Makefile
│   │   ├── README.md
│   │   ├── agent_installer.sh         // Manager가 생성한 스크립트로 curl 명령을 통한 원격 설치 프로그램
│   │   └── main.go                    // 에이전트의 main 소스 코드
│   └── klevr-manager
│       ├── Dockerfile                 // Docker 빌드를 위해 바이너리 아티팩트를 이 디렉토리로 전송
│       ├── Makefile
│       ├── README.md
│       └── main.go                    // 매니저의 main 소스 코드
├── conf
│   ├── klevr-manager-compose.yml      // Manager 실행을 위한 환경 설정 파일
│   ├── klevr-manager-db.sql.create    // Manager 초기화 및 실행을 위한 데이터베이스
│   ├── klevr-manager-db.sql.modify
│   └── klevr-manager-local.yml
├── console                            // Klevr WebConsole
│   ├── Makefile
│   ├── README.md
│   ├── jsconfig.json
│   ├── package-lock.json
│   ├── package.json
│   ├── public
│   ├── src
│   │   ├── components
│   │   │   ├── common
│   │   │   ├── credentials
│   │   │   ├── overview
│   │   │   ├── settings
│   │   │   ├── store
│   │   │   ├── task
│   │   │   └── zones
│   │   ├── pages
│   │   ├── theme
│   │   └── utils
│   └── yarn.lock
├── docker-compose-agent.yml
├── docker-compose-console.yml
├── docker-compose-demo.yml
├── go.mod
├── go.sum
├── pkg
│   ├── agent                           // agent 패키지 디렉토리
│   │   ├── agent.go
│   │   ├── common.go
│   │   ├── handshake.go
│   │   ├── primary_status_report.go
│   │   ├── protobuf
│   │   │   └── task.proto
│   │   ├── scheduler.go
│   │   ├── scheduler_primary.go
│   │   ├── scheduler_secondary.go
│   │   └── send_server.go
│   ├── common                           // common 패키지 디렉토리
│   │   ├── api_agent_model.go
│   │   ├── commander.go
│   │   ├── commands.go
│   │   ├── const.go
│   │   ├── context.go
│   │   ├── credential_model.go
│   │   ├── encrypt.go
│   │   ├── error.go
│   │   ├── http.go
│   │   ├── json_time.go
│   │   ├── jwt.go
│   │   ├── log.go
│   │   ├── md5.go
│   │   ├── orm.go
│   │   ├── queue.go
│   │   ├── security.go
│   │   ├── task_executor.go
│   │   └── task_model.go
│   ├── communicator                     // http communicator 패키지 디렉토리
│   │   ├── README.md
│   │   └── communicator.go
│   ├── manager                          // manager 패키지 디렉토리
│   │   ├── api.go
│   │   ├── api_agent.go
│   │   ├── api_console.go
│   │   ├── api_inner.go
│   │   ├── api_inner_model.go
│   │   ├── api_install.go
│   │   ├── cache.go
│   │   ├── context_constants.go
│   │   ├── handler.go
│   │   ├── persist_model.go
│   │   ├── repository.go
│   │   ├── server.go
│   │   ├── server_test.go
│   │   └── storage.go
│   └── rabbitmq                        rabbitmq 패키지 디렉토리
│       └── rabbitmq.go
├── scripts
│   ├── baremetal
│   └── linux_laptop
├── swag.sh
├── test
│   ├── common
│   │   └── task_executor_test.go
│   └── repository_test.go
└── version.properties
```

## Usage
### Swagger-UI API
* API 대시보드 URL : http://localhost:8090/swagger/index.html
### 1. Zone 
* 등록
  * [POST] /inner/groups
* 조회
  * [GET] /inner/groups/{groupID}
* 삭제
  * [DELETE] /inner/groups/{groupID}
### 2. API KEY
* 등록 
  * [POST] /inner/groups/{groupID}/apikey
* 조회
  * [GET] /inner/groups/{groupID}/apikey
* 수정
  * [PUT] /inner/groups/{groupID}/apikey
### 3. TASK 
* 등록
  * [POST] /inner/tasks
* 리스트
  * [GET] /inner/tasks
* 조회
  * [GET] /inner/tasks/{taskID}
* 취소
  * [DELETE] /inner/tasks/{taskID}
* 예약어 커맨드 정보
  * [GET] /inner/commands