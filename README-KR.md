![klevr_logo.png](https://raw.githubusercontent.com/Klevry/klevr/master/assets/klevr_logo.png)
# Klevr: Kloud-native everywhere
## 플랫폼 기반 SaaS 전달을 위한 인터커넥터
 * 분리된 네트워크를 위한 비동기 분산 인프라 관리 콘솔 및 에이전트.
 * 지원:
   * 온-프레미스 데이터센터의 베어메탈 서버
   * 회사/인트라넷의 PC/워크스테이션
   * 노트북
   * 퍼블릭 클라우드

## 웹 콘솔 및 KVstore 킥스타터
* docker-compose 명령어
```
git clone https://github.com/ralfyang/klevr.git
docker-compose up -d
```

## 다이어그램 개요
 * 이미지를 클릭 시 유튜브로 이동:
 * [![Diagram Overview](https://raw.githubusercontent.com/Klevry/klevr/master/assets/Klevr_diagram_overview.png)](https://youtu.be/xLkqm1vEmd0)

## 특징
 * **[Agent](./agent/)**
   * 프로비저닝: Docker, Kubernetes, Micro K8s(on Linux laptop) with Vagrant & VirtualBox, Prometheus 
   * 다운로드 및 실행: Hypervisor(via libvirt container), Terraform, Prometheus, Beacon
   * 메트릭 데이터 집계 및 전달
  * **[Web console](./webconsole/)**
   * 호스트 풀 관리
   * 리소스 관리
   * 프라이머리 호스트 관리 
   * 작업 관리 
   * 서비스 카탈로그 관리
   * 개발/스테이징/프로덕션에 서비스 전달
 * **도커 이미지**
   * [Webconsole](./Dockerfile/klevr_websonsole)(Webserver): [klevry:webconsole:latest](https://hub.docker.com/repository/docker/klevry/webconsole)
   * ~~[Beacon](./Dockerfile/beacon)(Primary agent health checker): [klevry/beacon:latest](https://hub.docker.com/repository/docker/klevry/beacon)~~
   * [Libvirt](./Dockerfile/libvirt)(Hypervisor): [klevry/libvirt:latest](https://hub.docker.com/repository/docker/klevry/libvirt)
   * 프로메테우스(컨테이너 모니터링)
   * 메트릭 크롤러
   * 작업 관리
 * **KV store([Consul](https://github.com/hashicorp/consul))**
   
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
├── README.md                   // 지금 보고있는 화면입니다 :)
├── docker-compose.yml          // 킥스타터: docker-compose를 이용한 부트스트래핑
├── Dockerfile                  // Docker 이미지 빌드를 위한 디렉토리
│   ├── libvirt
│   └── manager                 // 관리자의 실제 바이너리 파일은 도커 빌드를 위해 이 링크 디렉토리로 이동합니다.
├── assets
│   └── [Images & Contents]
├── cmd                         // 실제 아티팩트 fpr Klevr 에이전트 및 관리자 (웹 서버) 
│   ├── klevr-agent
│   │   ├── Makefile
│   │   ├── agent_installer.sh  // Manager가 생성 한 스크립트로 curl 명령을 통한 원격 설치 프로그램
│   │   ├── klevr               // 실제 Klevr 에이전트 바이너리
│   │   └── main.go             // 에이전트의 주요 소스 코드
│   └── klevr-manager
│       ├── Docker -> ../../Dockerfile/manager  // Docker 빌드를 위해 바이너리 아티팩트를 이 디렉터리로 전송  
│       ├── Makefile
│       └── main.go             // 매니저의 주요 소스 코드
├── conf
│   ├── Dump20200720.sql        // Manager의 초기화 및 실행을 위한 데이터베이스
│   └── klevr-manager-local.yml // Manager 실행을 위한 환경 설정 파일
├── pkg
│   ├── common                  // 'common' 패키지 디렉토리
│   │   ├── config.go
│   │   ├── error.go
│   │   ├── http.go
│   │   ├── log.go
│   │   └── orm.go
│   ├── communicator            // 'communicator' 패키지 디렉토리
│   │   ├── README.md
│   │   └── communicator.go
│   └── manager                 // 'manager' 패키지 디렉토리
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
└── scripts                    // 프로비저닝을 위한 운영 스크립트
    └── [Provisioning scripts]

```
