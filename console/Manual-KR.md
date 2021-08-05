# Klevr Console 매뉴얼

React로 구현된 Web 기반의 관리 도구  

# 로그인/사용자 활성화
로그인 페이지에 진입하면, Klevr Console 사용자의 활성화 상태 를 확인하여 활성화 및 로그인을 진행 합니다.

* 사용자 비활성화 상태이면 → Activator 페이지 이동  
  ![activate](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_activate.png)
* 사용자 활성화 상태이면 → 로그인 진행  
  ![signin](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_signin.png)  

# GNB _Zone Selector 
![GNB Zone Selector](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_gnb_zone_selector.png)
* Zone List 에서 Zone 을 선택하면 → Zone 별로 Resource 들 (task,agent...) 을 조회합니다. 

# GNB _Zone Page
## 페이지 소개
![zone intro](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_zone_intro.png)
* 해당 존의 API Key 와, Klevr 에 추가 되어있는 모든 Zone List 가 조회됩니다.
* API Key 는 각 zone 마다 가지고 있는 값이며, 입력된 적이 없으면 입력 후 ADD KEY 를 눌러 추가합니다.
* 현재는 API Key 값을 등록할 수만 있고, 새로 업데이트 할 수 는 없습니다.
* 이 API Key 는 Agent 를 추가 할 때 사용 되어는 키 입니다.
* 이 페이지에서 zone 을 ADD , DELETE 할 수 있습니다.

## Zone 추가
![zone add](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_zone_add.png)
* Zone Name , Platform 을 선택하여 zone 을 추가합니다.

## Zone 삭제
![zone del](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_zone_del.png)
* 현재 활성화 되어있는 zone 은 삭제할 수 없습니다.  
  (하나 뿐인 zone 을 삭제했을 경우와, 활성화된 zone 을 삭제했을 때 → 어떤 zone으로 이동시켜야 하냐는 문제를 고려하여 우선 비활성화 하는 방법을 선택했습니다.)



# Overview page
## 페이지 소개
![Overview intro](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_overview_intro.png)
* Overview 페이지에는 Task 리스트 , Credential 리스트 , Agent 리스트 가 있으며 데이터가 실시간으로 자동갱신이 되지 않기 때문에 각 리소스마다 refresh 버튼을 눌러 데이터를 갱신할 수 있습니다.
* Task 와 Credential 은 따로 페이지가 있어서 overview 화면에서는 일부만 잘라서 보여지며, 전체 리스트는 해당 페이지에서 확인 가능합니다. (view all 버튼을 누르거나, 왼쪽 사이드 메뉴탭에서 진입)

# Task Page
## 페이지 소개
![task intro](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_task_intro.png)
* Task 페이지에서는 Task List 와 Task Log 들을 조회할 수 있습니다.
* TaskList 는 전체 / Order(atOnce) / Scheduler(iteration) / Provisioning(longterm) 별로 구분해서 볼 수 있습니다.
* Task의  Status 가 변경될 때마다 화면이 자동으로 갱신되지 않기 때문에 refresh 버튼을 눌러 확인합니다.

## Task 추가
![task add](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_task_add_1.png)
![task add](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_task_add_2.png)
* Task 는 Order , Scheduler , Provision 세 가지 타입으로 추가할 수 있으며
* Scheduler 타입으로 선택하는 경우에는 iteration period 를 지정해줘야하는데 cron format 형식으로 입력받습니다.
* Target Agent 는 해당 존에 설치된 agent 들 중에서 고를 수 있고, 만약 특정 agent 를 선택하지 않는다면 (=none) 자동으로 agent 를 지정합니다.
* Command Type 은 inline 과 reserved 두가지가 있으며 inline 은 그냥 직접 커맨드를 입력하면 되고, reserved 는 Klevr Manger 에 설정되어있는 커맨드를 불러와 사용할 수 있는 타입입니다.

## Task 취소
* Task 의 Status 가 scheduled 또는 wait-polling 일 때에만 취소를 할 수 있습니다.
![task cancel](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_task_cancel_1.png)
* 정상적으로 취소가 되면 Success Message 가 나타나는데,
![task cancel](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_task_cancel_2.png)
* Task 의 Status 가 실시간으로 화면에 새로 보여지지 않기 때문에, Cancel 버튼을 누르는 순간 실제로는 Task 삭제할 수 없는 Status 로 변경이 이미 되어있을 수 있습니다. 그럼 Task 취소가 실패하게 되므로 이 때에는 Fail Message 를 보여주고 직접 refresh 를 해서 확인하도록 알려줍니다.
![task cancel](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_task_cancel_3.png)

# Credential Page
## Page 소개
![credential intro](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_credential_intro.png)
* Credential List 를 조회할 수 있는데, Key 는 그냥 보여지고 / value 는 hash 값으로 조회합니다.
* Credential Key, value 를 ADD , UPDATE(value값만) , DELETE 할 수 있습니다.

## Credential 추가
![credential add](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_credential_add.png)
* Key 와 Value 를 입력하고 Value 를 추가한 후에는 hash 로 밖에 조회할 수 없습니다.

## Credential 업데이트 (=수정)
![credential update](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_credential_update.png)
* Key 는 수정이 불가능하고, Value 만 업데이트 가능합니다.

# Agent Page
## Page 소개
![agent intro](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_agent_intro.png)
* Agent List 를 조회 / ADD 할 수 있습니다.
* Agent 는 하나의 Primary agent 와 다수의 Secondary agent 들이 있고 Role 을 통해 확인됩니다.
## Agent 추가
![agent add](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_agent_add.png)
* Agent 를 추가할 수 있는 스크립트를 만드려면 →  API Key , target Platform , Manager 주소 , Zone Id 가 필요한데
  * API KEY 는 Zones 페이지에서 등록 후 사용할 수 있고
  * Platform 은 모달에서 선택해야하며 (linux, baremetal, kubernetes)
  * Manger 주소는 자동으로 가져오고
  * Zone Id 도 자동으로 가져와 채워져 있습니다.
* 위 내용을 다 채워서 Create agent setup script 버튼을 누르면 스크립트가 생성되고 copy 하여 설치하면 됩니다.
![agent add modal](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_agent_add_modal.png)

## Agent List
![agent list](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_agent_list.png)
* 3-2 에서 추가한 Agent 들은 이렇게 리스트에 나타나게 되고 우측상단 새로고침 버튼으로 데이터를 갱신할 수 있습니다.

# Logs Page
## Page 소개
![logs intro](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_logs_intro.png)
* task 수행중 생긴 Log 들을 조회할 수 있습니다.
* 새로운 Log 가 생겨도 화면이 자동으로 갱신되지 않기 때문에 refresh 버튼을 눌러 확인합니다.

# Settings Page
## 페이지 소개
![settings intro](https://raw.githubusercontent.com/Klevry/klevr/master/assets/manual_console_settings_intro.png)
* 연결되어있는 Klevr Manger URL 을 조회할 수 있습니다.
* Password UPDATE 를 해줄 수 있으며, 현재 비밀번호가 일치했을 때 New password 에 입력된 값으로 비밀번호가 바뀝니다.

