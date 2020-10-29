package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gorilla/mux"
)

// CustomHeader name constants
const (
	CHeaderAPIKey         = "X-API-KEY"
	CHeaderAgentKey       = "X-AGENT-KEY"
	CHeaderHashCode       = "X-HASH-CODE"
	CHeaderZoneID         = "X-ZONE-ID"
	CHeaderSupportVersion = "X-SUPPORT-AGENT-VERSION"
	CHeaderTimestamp      = "X-TIMESTAMP"
)

type agentAPI int

// InitAgent initialize agent API
// @title Klevr-Manager API
// @version 1.0
// @description
// @contact.name mrchopa
// @contact.email ys3gods@gmail.com
// @BasePath /
func (api *API) InitAgent(agent *mux.Router) {
	logger.Debug("API InitAgent - init URI")

	agentAPI := agentAPI(0)

	registURI(agent, PUT, "/handshake", agentAPI.receiveHandshake)
	registURI(agent, PUT, "/{agentKey}", agentAPI.receivePolling)
	registURI(agent, GET, "/reports/{agentKey}", agentAPI.checkPrimaryInfo)

	// agent API 핸들러 추가
	agent.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ch := parseCustomHeader(r)
			ctx := CtxGetFromRequest(r)

			// TODO: Support agent version 입력 추가

			// APIKey 인증
			logger.Debug(r.RequestURI)
			if !authenticate(ctx, ch.ZoneID, ch.APIKey) {
				return
			}

			// TODO: hash 검증 로직 추가

			next.ServeHTTP(w, r)

			// TODO: 전송구간 암호화 로직 추가

			// TODO: hash 생성 로직 추가

			// response header 설정
			h := w.Header()

			h.Set(CHeaderAgentKey, ch.AgentKey)
			h.Set(CHeaderHashCode, ch.HashCode)
			h.Set(CHeaderSupportVersion, ch.SupportVersion)
			h.Set(CHeaderTimestamp, string(time.Now().UTC().Unix()))
		})
	})
}

func authenticate(ctx *common.Context, zoneID uint64, apiKey string) bool {
	if !GetDBConn(ctx).existAPIKey(apiKey, zoneID) {
		panic(common.NewHTTPError(401, "authentication failed"))
	}

	return true
}

func parseCustomHeader(r *http.Request) *common.CustomHeader {
	zoneID, _ := strconv.ParseUint(strings.Join(r.Header.Values(CHeaderZoneID), ""), 10, 64)
	ts, _ := strconv.ParseInt(strings.Join(r.Header.Values(CHeaderTimestamp), ""), 10, 64)

	h := &common.CustomHeader{
		APIKey:         strings.Join(r.Header.Values(CHeaderAPIKey), ""),
		AgentKey:       strings.Join(r.Header.Values(CHeaderAgentKey), ""),
		HashCode:       strings.Join(r.Header.Values(CHeaderHashCode), ""),
		ZoneID:         uint64(zoneID),
		SupportVersion: strings.Join(r.Header.Values(CHeaderSupportVersion), ""),
		Timestamp:      ts,
	}

	ctx := *CtxGetFromRequest(r)
	ctx.Put(common.CustomHeaderName, h)

	return h
}

func UpdateTaskStatus(tx *Tx, zoneID uint64, tasks *[]common.KlevrTask) {
	len := len(*tasks)

	if len > 0 {
		for _, t := range *tasks {

			tx.updateTask(&Tasks{
				Id:          t.ID,
				ExeAgentKey: t.AgentKey,
				Status:      t.Status,
				TaskDetail: &TaskDetail{
					Result: t.Result,
				},
			})
		}
	}
}

// ReceiveHandshake godoc
// @Summary 에이전트의 handshake 요청을 받아 처리한다.
// @Description 에이전트 프로세스가 기동시 최초 한번 handshake를 요청하여 에이전트 정보 등록 및 에이전트 실행에 필요한 실행 정보를 반환한다.
// @Tags agents
// @Accept json
// @Produce json
// @Router /agents/handshake [put]
// @Param X-API-KEY header string true "API KEY"
// @Param X-AGENT-KEY header string true "AGENT KEY"
// @Param X-ZONE-ID header string true "ZONE ID"
// @Param b body common.Body true "agent 정보"
// @Success 200 {object} common.Body
func (api *agentAPI) receiveHandshake(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	ch := ctx.Get(common.CustomHeaderName).(*common.CustomHeader)
	// var cr = &common.Request{r}

	tx := GetDBConn(ctx)
	var requestBody common.Body
	var paramAgent common.Me

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	paramAgent = requestBody.Me
	logger.Debug(fmt.Sprintf("CustomHeader : %v", ch))

	_, exist := tx.getAgentGroup(ch.ZoneID)

	if !exist {
		common.WriteHTTPError(400, w, nil, fmt.Sprintf("Does not exist zone for zoneId : %d", ch.ZoneID))
		return
	}

	agent := tx.getAgentByAgentKey(ch.AgentKey, ch.ZoneID)

	// agent 생성 or 수정
	upsertAgent(tx, agent, ch, &paramAgent)

	tx.Commit()

	// response 데이터 생성
	rb := &common.Body{}

	// primary 조회
	var oldPrimaryAgentKey string
	rb.Agent.Primary, oldPrimaryAgentKey = getPrimary(ctx, tx, ch.ZoneID, agent)


	// 접속한 agent 정보
	me := &rb.Me

	me.HmacKey = agent.HmacKey
	me.EncKey = agent.EncKey
	// me.CallCycle = 0 // seconds
	// me.LogLevel = "DEBUG"

	// Primary agent인 경우 node 정보 추가
	if ch.AgentKey == rb.Agent.Primary.AgentKey {
		rb.Agent.Nodes = getNodes(tx, ch.ZoneID)
	}

	b, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)

	AddEvent(&KlevrEvent{
		EventType: AgentConnect,
		AgentKey:  agent.AgentKey,
		GroupID:   agent.GroupId,
		EventTime: &common.JSONTime{Time: time.Now().UTC()},
	})

	if ch.AgentKey == rb.Agent.Primary.AgentKey {
		AddEvent(&KlevrEvent{
			EventType: PrimaryElected,
			AgentKey:  agent.AgentKey,
			GroupID:   agent.GroupId,
			EventTime: &common.JSONTime{Time: time.Now().UTC()},
		})
	}

	if oldPrimaryAgentKey != "" && oldPrimaryAgentKey != rb.Agent.Primary.AgentKey {
		AddEvent(&KlevrEvent{
			EventType: PrimaryRetire,
			AgentKey:  oldPrimaryAgentKey,
			GroupID:   agent.GroupId,
			EventTime: &common.JSONTime{Time: time.Now().UTC()},
		})
	}
}

// ReceivePolling godoc
// @Summary primary 에이전트의 polling 요청을 받아 처리한다.
// @Description primary 에이전트의 polling 요청을 받아 primary 에이전트의 실행정보 갱신, nodes 정보 갱신, task 할당 및 상태 업데이트를 수행한다.
// @Tags agents
// @Accept json
// @Produce json
// @Router /agents/{agentKey} [put]
// @Param X-API-KEY header string true "API KEY"
// @Param X-AGENT-KEY header string true "AGENT KEY"
// @Param X-ZONE-ID header string true "ZONE ID"
// @Param agentKey path string true "agent key"
// @Param b body common.Body true "agent 정보"
// @Success 200 {object} common.Body
func (api *agentAPI) receivePolling(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	ch := ctx.Get(common.CustomHeaderName).(*common.CustomHeader)
	// var cr = &common.Request{r}
	tx := GetDBConn(ctx)
	var param common.Body

	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	// response 데이터 생성
	rb := &common.Body{}

	// agent access 정보 갱신
	agent := updateAgentAccess(tx, ch.AgentKey, ch.ZoneID)
	logger.Debugf("%+v", agent)

	// 수행한 task 상태 정보 업데이트
	var taskLength = len(param.Task)

	if taskLength > 0 {
		var pTaskMap = make(map[uint64]*Tasks)
		var tIds = make([]uint64, len(param.Task))

		for i, task := range param.Task {
			tIds[i] = task.ID
		}

		pTasks, _ := tx.getTasksByIds(tIds)
		for _, pt := range *pTasks {
			pTaskMap[pt.Id] = &pt
		}

		updateTaskStatus(ctx, pTaskMap, &param.Task)
	}

	rb.Agent.Primary, _ = getPrimary(ctx, tx, ch.ZoneID, agent)

	if agent.AgentKey == rb.Agent.Primary.AgentKey {
		// TODO: primary agent 실행 파라미터 update
		// rb.Me.CallCycle =
		// rb.Me.LogLevel =

		// agent zone 상태 정보 업데이트
		nodes := param.Agent.Nodes
		nodeLength := len(nodes)
		arrAgent := make([]Agents, nodeLength)

		for i, a := range nodes {
			arrAgent[i].AgentKey = a.AgentKey
			arrAgent[i].LastAliveCheckTime = a.LastAliveCheckTime.Time
			arrAgent[i].Cpu = a.Core
			arrAgent[i].Memory = a.Memory
			arrAgent[i].Disk = a.Disk

			if agent.AgentKey == a.AgentKey {
				arrAgent[i].IsActive = 1
			} else {
				arrAgent[i].IsActive = boolToByte(a.IsActive)
			}
		}

		tx.updateZoneStatus(&arrAgent)

		tx.Commit()

		// 신규 task 할당
		nTasks, cnt := tx.getTasksWithSteps(ch.ZoneID, []string{string(common.WaitPolling), string(common.HandOver)})
		if cnt > 0 {
			var dtos []common.KlevrTask = make([]common.KlevrTask, len(*nTasks))

			for i, t := range *nTasks {
				dtos[i] = *TaskPersistToDto(&t)
			}

			rb.Task = dtos

			AddHandOverTasks(nTasks)
		}

		// node 정보 추가
		rb.Agent.Nodes = getNodes(tx, ch.ZoneID)
	}

	b, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

func updateTaskStatus(ctx *common.Context, oTasks map[uint64]*Tasks, uTasks *[]common.KlevrTask) {
	var length = len(*uTasks)
	var events = make([]KlevrEvent, 0, length*2)

	tx := GetDBConn(ctx)

	for _, t := range *uTasks {
		oTask := oTasks[t.ID]

		// Task 상태 이상으로 오류 종료 처리
		if t.Status == common.Scheduled || t.Status == common.WaitPolling || t.Status == common.HandOver {
			oTask.Status = common.Failed
			oTask.Logs.Logs = "Invalid Task Status Updated. - " + string(t.Status)

			events = append(events, KlevrEvent{
				EventType: TaskCallback,
				AgentKey:  oTask.AgentKey,
				GroupID:   oTask.ZoneId,
				Result:    NewKlevrEventTaskResultString(oTask, true, false, false, t.Result, t.Log, "Invalid Task Status", string(t.Status)),
				EventTime: &common.JSONTime{Time: time.Now().UTC()},
			})
		} else {
			var complete = false
			var success = false
			var isCommandError = false
			var sendEvent = true
			var errorMessage string

			switch t.Status {
			case common.WaitExec:
				sendEvent = false
			case common.Running:
				// 한 단계 이상이 완료 되어야 event 발송
				if t.CurrentStep < 2 {
					sendEvent = false
				}
			case common.Recovering:
				if t.FailedStep > 0 {
					isCommandError = true
					errorMessage = "Error occurred during task step execution"
				}
				oTask.TaskDetail.FailedStep = t.FailedStep
			case common.Complete:
				complete = true
				success = true
			case common.FailedRecover:
				oTask.TaskDetail.IsFailedRecover = true
				complete = true
			case common.Failed:
				if t.FailedStep > 0 {
					isCommandError = true
					errorMessage = "Error occurred during task step execution"
				}
				complete = true
			case common.Canceled:
				complete = true
			case common.Stopped:
				complete = true
			default:
				panic("invalid task status - " + t.Status)
			}

			oTask.TaskDetail.CurrentStep = t.CurrentStep
			oTask.Status = t.Status
			oTask.Logs.Logs = t.Log

			if sendEvent {
				events = append(events, KlevrEvent{
					EventType: TaskCallback,
					AgentKey:  oTask.AgentKey,
					GroupID:   oTask.ZoneId,
					Result:    NewKlevrEventTaskResultString(oTask, complete, success, isCommandError, t.Result, t.Log, errorMessage, t.Log),
					EventTime: &common.JSONTime{Time: time.Now().UTC()},
				})
			}
		}

		tx.updateTask(oTask)
	}

	AddEvents(&events)
}

// CheckPrimaryInfo godoc
// @Summary secondary 에이전트의 primary 상태 확인 요청을 처리한다.
// @Description secondary 에이전트의 primary 에이전트 상태 확인 요청을 받아 primary 재선출 및 primary 정보를 반환한다.
// @Tags agents
// @Accept json
// @Produce json
// @Router /agents/reports/{agentKey} [get]
// @Param X-API-KEY header string true "API KEY"
// @Param X-AGENT-KEY header string true "AGENT KEY"
// @Param X-ZONE-ID header string true "ZONE ID"
// @Param agentKey path string true "agent key"
// @Param b body common.Body true "agent 정보"
// @Success 200 {object} common.Body
func (api *agentAPI) checkPrimaryInfo(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	ch := ctx.Get(common.CustomHeaderName).(*common.CustomHeader)
	// var cr = &common.Request{r}
	tx := GetDBConn(ctx)
	var param common.Body

	err := json.NewDecoder(r.Body).Decode(&param)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	// response 데이터 생성
	rb := &common.Body{}

	// agent access 정보 갱신
	agent := updateAgentAccess(tx, ch.AgentKey, ch.ZoneID)

	var oldPrimaryAgentKey string
	rb.Agent.Primary, oldPrimaryAgentKey = getPrimary(ctx, tx, ch.ZoneID, agent)

	if ch.AgentKey == rb.Agent.Primary.AgentKey {
		rb.Agent.Nodes = getNodes(tx, ch.ZoneID)
	}

	b, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)

	if ch.AgentKey == rb.Agent.Primary.AgentKey {
		AddEvent(&KlevrEvent{
			EventType: PrimaryElected,
			AgentKey:  agent.AgentKey,
			GroupID:   agent.GroupId,
			EventTime: &common.JSONTime{Time: time.Now().UTC()},
		})
	}

	if oldPrimaryAgentKey != "" && oldPrimaryAgentKey != rb.Agent.Primary.AgentKey {
		AddEvent(&KlevrEvent{
			EventType: PrimaryRetire,
			AgentKey:  oldPrimaryAgentKey,
			GroupID:   agent.GroupId,
			EventTime: &common.JSONTime{Time: time.Now().UTC()},
		})
	}
}

func getNodes(tx *Tx, zoneID uint64) []common.Agent {
	cnt, agents := tx.getAgentsByGroupId(zoneID)
	nodes := make([]common.Agent, cnt)

	if cnt > 0 {
		for i, a := range *agents {
			nodes[i] = common.Agent{
				AgentKey: a.AgentKey,
				IP:       a.Ip,
				Port:     a.Port,
				Version:  a.Version,
				Resource: &common.Resource{
					Core:   a.Cpu,
					Memory: a.Memory,
					Disk:   a.Disk,
				},
			}
		}

		return nodes
	}

	return nil
}

func updateAgentAccess(tx *Tx, agentKey string, zoneID uint64) *Agents {
	agent := tx.getAgentByAgentKey(agentKey, zoneID)

	// agent 접속 시간 갱신
	// agent.IsActive = 1
	// agent.LastAccessTime = time.Now().UTC()
	// tx.updateAgent(agent)
	tx.updateAccessAgent(agent.Id, time.Now().UTC())

	return agent
}

func getPrimary(ctx *common.Context, tx *Tx, zoneID uint64, curAgent *Agents) (common.Primary, string) {

	// primary agent 정보
	groupPrimary := tx.getPrimaryAgent(zoneID)
	var primaryAgent *Agents
	var oldPrimaryAgentKey string

	if groupPrimary.AgentId == curAgent.Id {
		primaryAgent = curAgent
	} else if groupPrimary.AgentId == 0 {
		primaryAgent = electPrimary(ctx, zoneID, curAgent.Id, false)
	} else {
		primaryAgent = tx.getAgentByID(groupPrimary.AgentId)
		oldPrimaryAgentKey = primaryAgent.AgentKey

		logger.Debugf("primaryAgent : %+v", primaryAgent)

		if primaryAgent.Id == 0 || primaryAgent.IsActive == 0 {
			primaryAgent = electPrimary(ctx, zoneID, curAgent.Id, true)

			logger.Debugf("changed primaryAgent : %+v", primaryAgent)
		}
	}

	return common.Primary{
		AgentKey:       primaryAgent.AgentKey,
		IP:             primaryAgent.Ip,
		Port:           primaryAgent.Port,
		IsActive:       byteToBool(primaryAgent.IsActive),
		LastAccessTime: primaryAgent.LastAccessTime.UTC().Unix(),
	}, oldPrimaryAgentKey
}

// primary agent 선출
func electPrimary(ctx *common.Context, zoneID uint64, agentID uint64, oldDel bool) *Agents {
	logger.Debugf("electPrimary for %d", zoneID)

	var tx *Tx
	var agent *Agents

	common.Block{
		Try: func() {
			tx = &Tx{CtxGetDbConn(ctx).NewSession()}

			if oldDel {
				tx.deletePrimaryAgent(zoneID)
			}

			pa := &PrimaryAgents{
				GroupId: zoneID,
				AgentId: agentID,
			}

			cnt, err := tx.insertPrimaryAgent(pa)

			if err != nil {
				pa = tx.getPrimaryAgent(zoneID)
			} else if cnt != 1 {
				logger.Warning(fmt.Sprintf("insert primary agent cnt : %d", cnt))
				common.Throw(common.NewStandardError("elect primary failed."))
			}

			if pa.AgentId == 0 {
				logger.Debugf("invalid primary agent : %v", pa)
				common.Throw(common.NewStandardError("elect primary failed."))
			}

			agent = tx.getAgentByID(pa.AgentId)

			if agent.Id == 0 {
				logger.Warning(fmt.Sprintf("primary agent not exist for id : %d, [%v]", agent.Id, agent))
				common.Throw(common.NewStandardError("elect primary failed."))
			}

			tx.Commit()
		},
		Catch: func(e error) {
			if tx != nil {
				tx.Rollback()
			}

			logger.Warningf("%+v", e)
			common.Throw(e)
		},
		Finally: func() {
			if tx != nil && !tx.IsClosed() {
				tx.Close()
			}
		},
	}.Do()

	return agent
}

func upsertAgent(tx *Tx, agent *Agents, ch *common.CustomHeader, paramAgent *common.Me) {
	if agent.AgentKey == "" { // 처음 접속하는 에이전트일 경우 신규 등록
		agent.AgentKey = ch.AgentKey
		agent.GroupId = ch.ZoneID
		agent.IsActive = 1
		agent.LastAccessTime = time.Now().UTC()
		agent.Ip = paramAgent.IP
		agent.Port = paramAgent.Port
		agent.Cpu = paramAgent.Resource.Core
		agent.Memory = paramAgent.Resource.Memory
		agent.Disk = paramAgent.Resource.Disk
		agent.HmacKey = common.GetKey(16)
		agent.EncKey = common.GetKey(32)

		tx.addAgent(agent)
	} else { // 기존에 등록된 에이전트 재접속일 경우 접속 정보 업데이트
		agent.IsActive = 1
		agent.LastAccessTime = time.Now().UTC()
		agent.Ip = paramAgent.IP
		agent.Port = paramAgent.Port
		agent.Cpu = paramAgent.Resource.Core
		agent.Memory = paramAgent.Resource.Memory
		agent.Disk = paramAgent.Resource.Disk
		agent.HmacKey = common.GetKey(16)
		agent.EncKey = common.GetKey(32)

		tx.updateAgent(agent)
	}
}

func byteToBool(b byte) bool {
	if b == 0 {
		return false
	}

	return true
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}

	return 0
}
