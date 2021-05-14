package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gorilla/mux"
)

// CustomHeader name constants
const (
	CHeaderAPIKey           = "X-API-KEY"
	CHeaderAgentKey         = "X-AGENT-KEY"
	CHeaderHashCode         = "X-HASH-CODE"
	CHeaderZoneID           = "X-ZONE-ID"
	CHeaderSupportVersion   = "X-SUPPORT-AGENT-VERSION"
	CHeaderTimestamp        = "X-TIMESTAMP"
	CHeaderPayloadEncrypted = "X-PAYLOAD-ENC"
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
			if !api.authenticate(ctx, ch.ZoneID, ch.APIKey, ch.AgentKey) {
				logger.Debug(fmt.Sprintf("failed authenticate: %v", ch))
				w.WriteHeader(401)
				fmt.Fprintf(w, "%s", "failed authenticate")
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
			h.Set(CHeaderTimestamp, strconv.FormatInt(time.Now().UTC().Unix(), 10))
		})
	})
}

func (api *API) authenticate(ctx *common.Context, zoneID uint64, apiKey, agentKey string) bool {
	logger.Debug(fmt.Sprintf("[authenticate info] zoneID:%d, apiKey:%s, agentKey:%s", zoneID, apiKey, agentKey))

	value, bExist := api.BlockKeyMap.Get(agentKey)
	if bExist && apiKey == value.(string) {
		logger.Debugf("[BlockKeyMap(Get)] zoneID(%d), apiKey(%s), agentKey(%s)", zoneID, apiKey, agentKey)
		return false
	}

	apiKeyMap := api.APIKeyMap

	if !apiKeyMap.Has(strconv.FormatUint(zoneID, 10)) {
		tx := GetDBConn(ctx)
		apiKey, ok := tx.getAPIKey(zoneID)

		if ok && apiKey.GroupId > 0 {
			manager := ctx.Get(CtxServer).(*KlevrManager)

			apiKeyMap.Set(strconv.FormatUint(zoneID, 10), manager.decrypt(apiKey.ApiKey))
			logger.Debugf("[apiKeyMap(Set)] zoneID(%d), apiKey(%s)", zoneID, apiKey.ApiKey)
		}
	}

	ifval, aExist := apiKeyMap.Get(strconv.FormatUint(zoneID, 10))

	if aExist {
		val := ifval.(string)
		logger.Debugf("[apiKeyMap(Get)] apiKey.db(%s), apiKey.in(%s)", val, apiKey)

		if apiKey != "" && val == apiKey {
			return true
		}
	}

	api.BlockKeyMap.Set(agentKey, apiKey)
	logger.Debugf("[BlockKeyMap(Set)] zoneID(%d), apiKey(%s), agentKey(%s)", zoneID, apiKey, agentKey)

	logger.Warningf("API key not matched - apiKey(%s), agentKey(%s)", apiKey, agentKey)
	panic(common.NewHTTPError(401, "authentication failed"))
}

func parseCustomHeader(r *http.Request) *common.CustomHeader {
	zoneID, _ := strconv.ParseUint(strings.Join(r.Header.Values(CHeaderZoneID), ""), 10, 64)
	ts, _ := strconv.ParseInt(strings.Join(r.Header.Values(CHeaderTimestamp), ""), 10, 64)
	payloadEncrypted, _ := strconv.ParseBool(strings.Join(r.Header.Values(CHeaderPayloadEncrypted), ""))

	h := &common.CustomHeader{
		APIKey:           strings.Join(r.Header.Values(CHeaderAPIKey), ""),
		AgentKey:         strings.Join(r.Header.Values(CHeaderAgentKey), ""),
		HashCode:         strings.Join(r.Header.Values(CHeaderHashCode), ""),
		ZoneID:           uint64(zoneID),
		SupportVersion:   strings.Join(r.Header.Values(CHeaderSupportVersion), ""),
		Timestamp:        ts,
		PayloadEncrypted: payloadEncrypted,
	}

	ctx := *CtxGetFromRequest(r)
	ctx.Put(common.CustomHeaderName, h)

	return h
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

	logger.Debug(fmt.Sprintf("Handshake Body: %+v", r.Body))
	logger.Debug(fmt.Sprintf("Agent: %v", requestBody.Me))
	logger.Debug(fmt.Sprintf("CustomHeader: %v", ch))

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	paramAgent = requestBody.Me

	_, exist := tx.getAgentGroup(ch.ZoneID)

	if !exist {
		common.WriteHTTPError(400, w, nil, fmt.Sprintf("Does not exist zone for zoneId : %d", ch.ZoneID))
		return
	}

	//agent := tx.getAgentByAgentKey(ch.AgentKey, ch.ZoneID)
	txManager := NewAgentStorage()
	agent := txManager.GetAgentByAgentKey(tx, ch.AgentKey, ch.ZoneID)

	// agent 생성 or 수정
	upsertAgent(ctx, tx, agent, ch, &paramAgent)

	tx.Commit()

	// response 데이터 생성
	rb := &common.Body{}

	// primary 조회
	var oldPrimaryAgentKey string
	rb.Agent.Primary, oldPrimaryAgentKey = getPrimary(ctx, tx, ch.ZoneID, agent)

	// 접속한 agent 정보
	me := &rb.Me

	manager := ctx.Get(CtxServer).(*KlevrManager)

	me.HmacKey = manager.decrypt(agent.HmacKey)
	me.EncKey = manager.decrypt(agent.EncKey)
	me.CallCycle = manager.Config.Agent.CallCycle // seconds
	// me.LogLevel = "DEBUG"

	// Primary agent인 경우 node 정보 추가
	if ch.AgentKey == rb.Agent.Primary.AgentKey {
		rb.Agent.Nodes = getNodes(ctx, tx, ch.ZoneID)
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

	if oldPrimaryAgentKey != "" && oldPrimaryAgentKey != rb.Agent.Primary.AgentKey {
		AddEvent(&KlevrEvent{
			EventType: PrimaryRetire,
			AgentKey:  oldPrimaryAgentKey,
			GroupID:   agent.GroupId,
			EventTime: &common.JSONTime{Time: time.Now().UTC()},
		})
	}

	if ch.AgentKey == rb.Agent.Primary.AgentKey {
		AddEvent(&KlevrEvent{
			EventType: PrimaryElected,
			AgentKey:  agent.AgentKey,
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

	logger.Debugf("polling received data : [%+v]", param)

	// response 데이터 생성
	rb := &common.Body{}

	// agent access 정보 갱신
	agent := updateAgentAccess(tx, ch.AgentKey, ch.ZoneID)
	logger.Debugf("%+v", agent)

	// 수행한 task 상태 정보 업데이트
	var taskLength = len(param.Task)

	manager := ctx.Get(CtxServer).(*KlevrManager)

	if taskLength > 0 {
		var pTaskMap = make(map[uint64]Tasks)
		var tIds = make([]uint64, len(param.Task))

		for i, task := range param.Task {
			tIds[i] = task.ID
		}

		pTasks, _ := tx.getTasksByIds(manager, tIds)
		for _, pt := range *pTasks {
			pTaskMap[pt.Id] = pt
		}

		logger.Debugf("map for update - [%+v]", pTaskMap)

		updateTaskStatus(ctx, pTaskMap, &param.Task)
	}

	rb.Agent.Primary, _ = getPrimary(ctx, tx, ch.ZoneID, agent)

	if agent.AgentKey == rb.Agent.Primary.AgentKey {
		// TODO: primary agent 실행 파라미터 update
		rb.Me.CallCycle = manager.Config.Agent.CallCycle
		// rb.Me.LogLevel =

		// agent zone 상태 정보 업데이트
		nodes := param.Agent.Nodes
		nodeLength := len(nodes)
		arrAgent := make([]Agents, nodeLength)

		manager := ctx.Get(CtxServer).(*KlevrManager)

		agentKeys := make([]string, nodeLength)
		taskIDs := make([]uint64, nodeLength)
		for i, a := range nodes {
			arrAgent[i].AgentKey = a.AgentKey
			if a.LastAliveCheckTime != nil {
				arrAgent[i].LastAliveCheckTime = a.LastAliveCheckTime.Time
			}
			arrAgent[i].Cpu = manager.encrypt(strconv.Itoa(a.Core))
			arrAgent[i].Memory = manager.encrypt(strconv.Itoa(a.Memory))
			arrAgent[i].Disk = manager.encrypt(strconv.Itoa(a.Disk))
			arrAgent[i].FreeMemory = manager.encrypt(strconv.Itoa(a.FreeMemory))
			arrAgent[i].FreeDisk = manager.encrypt(strconv.Itoa(a.FreeDisk))

			if agent.AgentKey == a.AgentKey {
				arrAgent[i].IsActive = 1
			} else {
				arrAgent[i].IsActive = boolToByte(a.IsActive)
				if a.IsActive == false {
					if tid, ok := CheckShutdownTask(a.AgentKey); ok {
						agentKeys = append(agentKeys, a.AgentKey)
						taskIDs = append(taskIDs, tid)
					}
				}
			}
		}

		//tx.updateZoneStatus(&arrAgent)
		NewAgentStorage().UpdateZoneStatus(tx, ch.ZoneID, arrAgent)
		tx.updateShutdownTasks(taskIDs)

		RemoveShutdownTask(agentKeys)

		tx.Commit()

		// Credential 조회
		nCredentials, cnt := tx.getCredentials(ch.ZoneID)

		// 신규 task 할당
		nTasks, cnt := tx.getTasksWithSteps(manager, ch.ZoneID, []string{string(common.WaitPolling), string(common.HandOver)})
		if cnt > 0 {
			var dtos []common.KlevrTask = make([]common.KlevrTask, len(*nTasks))

			for i, t := range *nTasks {
				t = TaskMatchingCredential(manager, t, nCredentials)
				dtos[i] = *TaskPersistToDto(&t)
			}

			rb.Task = dtos

			AddHandOverTasks(nTasks)
		}

		// node 정보 추가
		rb.Agent.Nodes = getNodes(ctx, tx, ch.ZoneID)
	}

	b, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)

	logger.Debug(string(b))
}

func updateTaskStatus(ctx *common.Context, oTasks map[uint64]Tasks, uTasks *[]common.KlevrTask) {
	var length = len(*uTasks)
	var events = make([]KlevrEvent, 0, length*2)

	tx := GetDBConn(ctx)

	manager := ctx.Get(CtxServer).(*KlevrManager)

	for _, t := range *uTasks {
		oTask := oTasks[t.ID]
		oTask.ExeAgentKey = t.ExeAgentKey

		logger.Debugf("updateTaskStatus : [%+v]", oTask)

		// Task 상태 이상으로 오류 종료 처리
		if t.Status == common.Scheduled || t.Status == common.WaitPolling || t.Status == common.HandOver {
			if t.EventHookSendingType != common.EventHookWithSuccess {
				oTask.Status = common.Failed
				oTask.Logs.Logs = "Invalid Task Status Updated. - " + string(t.Status)

				events = append(events, KlevrEvent{
					EventType: TaskCallback,
					AgentKey:  oTask.AgentKey,
					GroupID:   oTask.ZoneId,
					Result:    NewKlevrEventTaskResultString(&oTask, true, false, false, t.Result, t.Log, "Invalid Task Status", string(t.Status)),
					EventTime: &common.JSONTime{Time: time.Now().UTC()},
				})
			}
		} else {
			var complete = false
			var success = false
			var isCommandError = false
			var sendEvent = false
			var errorMessage string

			switch t.Status {
			case common.WaitExec:
			case common.Started:
				if t.EventHookSendingType == common.EventHookWithBothEnds {
					sendEvent = true
				}
			case common.Running:
				if t.EventHookSendingType == common.EventHookWithEachSteps {
					sendEvent = true
				}
			case common.WaitInterationSchedule:
				if t.EventHookSendingType == common.EventHookWithSuccess {
					sendEvent = true
				}

				success = true
			case common.Recovering:
				if t.EventHookSendingType == common.EventHookWithEachSteps {
					sendEvent = true
				}

				if t.FailedStep > 0 {
					isCommandError = true
					errorMessage = "Error occurred during task step execution"
				}
				oTask.TaskDetail.FailedStep = t.FailedStep
			case common.Complete:
				if t.EventHookSendingType == common.EventHookWithConclusion ||
					t.EventHookSendingType == common.EventHookWithBothEnds ||
					t.EventHookSendingType == common.EventHookWithSuccess {
					sendEvent = true
				}

				complete = true
				success = true
			case common.FailedRecover:
				if t.EventHookSendingType == common.EventHookWithConclusion ||
					t.EventHookSendingType == common.EventHookWithBothEnds ||
					t.EventHookSendingType == common.EventHookWithFailed {
					sendEvent = true
				}

				oTask.TaskDetail.IsFailedRecover = true
				isCommandError = true
				complete = true
			case common.Failed:
				if t.EventHookSendingType == common.EventHookWithConclusion ||
					t.EventHookSendingType == common.EventHookWithBothEnds ||
					t.EventHookSendingType == common.EventHookWithFailed {
					sendEvent = true
				}

				if t.FailedStep > 0 {
					isCommandError = true
					errorMessage = "Error occurred during task step execution"
				}
				complete = true
			case common.Canceled:
				if t.EventHookSendingType == common.EventHookWithConclusion ||
					t.EventHookSendingType == common.EventHookWithBothEnds {
					sendEvent = true
				}

				complete = true
			case common.Stopped:
				if t.EventHookSendingType == common.EventHookWithConclusion ||
					t.EventHookSendingType == common.EventHookWithBothEnds {
					sendEvent = true
				}

				complete = true
			default:
				panic("invalid task status - " + t.Status)
			}

			if t.EventHookSendingType == common.EventHookWithAll || t.EventHookSendingType == "" {
				sendEvent = true
			} else if t.EventHookSendingType == common.EventHookWithChangedResult && t.IsChangedResult {
				sendEvent = true
			}

			oTask.TaskDetail.CurrentStep = t.CurrentStep
			oTask.Status = t.Status
			oTask.Logs.Logs = t.Log
			oTask.TaskDetail.Result = t.Result

			if sendEvent {
				events = append(events, KlevrEvent{
					EventType: TaskCallback,
					AgentKey:  oTask.AgentKey,
					GroupID:   oTask.ZoneId,
					Result:    NewKlevrEventTaskResultString(&oTask, complete, success, isCommandError, t.Result, t.Log, errorMessage, t.Log),
					EventTime: &common.JSONTime{Time: time.Now().UTC()},
				})
			}
		}

		tx.updateTask(manager, &oTask)
		tx.Commit()
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

	// response 데이터 생성
	rb := &common.Body{}

	// agent access 정보 갱신
	agent := updateAgentAccess(tx, ch.AgentKey, ch.ZoneID)

	var oldPrimaryAgentKey string
	rb.Agent.Primary, oldPrimaryAgentKey = getPrimary(ctx, tx, ch.ZoneID, agent)

	if ch.AgentKey == rb.Agent.Primary.AgentKey {
		rb.Agent.Nodes = getNodes(ctx, tx, ch.ZoneID)
	}

	b, err := json.Marshal(rb)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)

	if oldPrimaryAgentKey != "" && oldPrimaryAgentKey != rb.Agent.Primary.AgentKey {
		AddEvent(&KlevrEvent{
			EventType: PrimaryRetire,
			AgentKey:  oldPrimaryAgentKey,
			GroupID:   agent.GroupId,
			EventTime: &common.JSONTime{Time: time.Now().UTC()},
		})
	}

	if ch.AgentKey == rb.Agent.Primary.AgentKey {
		AddEvent(&KlevrEvent{
			EventType: PrimaryElected,
			AgentKey:  agent.AgentKey,
			GroupID:   agent.GroupId,
			EventTime: &common.JSONTime{Time: time.Now().UTC()},
		})
	}
}

func getNodes(ctx *common.Context, tx *Tx, zoneID uint64) []common.Agent {
	//cnt, agents := tx.getAgentsByGroupId(zoneID)
	txManager := NewAgentStorage()
	cnt, agents := txManager.GetAgentsByZoneID(tx, zoneID)

	nodes := make([]common.Agent, cnt)

	manager := ctx.Get(CtxServer).(*KlevrManager)

	if cnt > 0 {
		for i, a := range *agents {
			nodes[i] = common.Agent{
				AgentKey:           a.AgentKey,
				IP:                 a.Ip,
				Port:               a.Port,
				Version:            a.Version,
				IsActive:           byteToBool(a.IsActive),
				LastAliveCheckTime: &common.JSONTime{Time: a.LastAliveCheckTime},
				Resource:           &common.Resource{},
			}

			core, _ := strconv.Atoi(manager.decrypt(a.Cpu))
			memory, _ := strconv.Atoi(manager.decrypt(a.Memory))
			disk, _ := strconv.Atoi(manager.decrypt(a.Disk))
			freeMemory, _ := strconv.Atoi(manager.decrypt(a.FreeMemory))
			freeDisk, _ := strconv.Atoi(manager.decrypt(a.FreeDisk))

			nodes[i].Core = core
			nodes[i].Memory = memory
			nodes[i].Disk = disk
			nodes[i].FreeMemory = freeMemory
			nodes[i].FreeDisk = freeDisk
		}

		return nodes
	}

	return nil
}

func updateAgentAccess(tx *Tx, agentKey string, zoneID uint64) *Agents {
	txManager := NewAgentStorage()

	//agent := tx.getAgentByAgentKey(agentKey, zoneID)
	agent := txManager.GetAgentByAgentKey(tx, agentKey, zoneID)

	oldStatus := agent.IsActive

	curTime := time.Now().UTC()

	// agent 접속 시간 갱신
	// agent.IsActive = 1
	// agent.LastAccessTime = time.Now().UTC()
	// tx.updateAgent(agent)

	//cnt := tx.updateAccessAgent(agentKey, curTime)
	cnt := txManager.UpdateAccessAgent(tx, zoneID, agentKey, curTime)

	if oldStatus == 0 && cnt > 0 {
		AddEvent(&KlevrEvent{
			EventType: AgentConnect,
			AgentKey:  agentKey,
			GroupID:   agent.GroupId,
			EventTime: &common.JSONTime{Time: curTime},
		})
	}

	tx.Commit()

	return agent
}

func getPrimary(ctx *common.Context, tx *Tx, zoneID uint64, curAgent *Agents) (common.Primary, string) {
	primaryMutex := ctx.Get(CtxPrimary).(*sync.Mutex)

	primaryMutex.Lock()
	defer primaryMutex.Unlock()

	// primary agent 정보
	groupPrimary, _ := tx.getPrimaryAgent(zoneID)
	var primaryAgent *Agents
	var oldPrimaryAgentKey string

	if groupPrimary.AgentId == curAgent.Id {
		primaryAgent = curAgent
	} else if groupPrimary.GroupId == 0 && groupPrimary.AgentId == 0 {
		primaryAgent = electPrimary(ctx, zoneID, curAgent.Id, false)
	} else {
		//primaryAgent = tx.getAgentByID(groupPrimary.AgentId)
		txManager := NewAgentStorage()
		primaryAgent = txManager.GetAgentByID(tx, zoneID, groupPrimary.AgentId)
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
	logger.Debugf("electPrimary for zone(%d), agent(%d)", zoneID, agentID)

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
				pa, _ = tx.getPrimaryAgent(zoneID)
			} else if cnt != 1 {
				logger.Warning(fmt.Sprintf("insert primary agent cnt : %d", cnt))
				common.Throw(common.NewStandardError("elect primary failed."))
			}

			if pa.AgentId == 0 {
				logger.Debugf("invalid primary agent : %v", pa)
				common.Throw(common.NewStandardError("elect primary failed."))
			}

			//agent = tx.getAgentByID(pa.AgentId)
			txManager := NewAgentStorage()
			agent = txManager.GetAgentByID(tx, zoneID, pa.AgentId)

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

func upsertAgent(ctx *common.Context, tx *Tx, agent *Agents, ch *common.CustomHeader, paramAgent *common.Me) {
	manager := ctx.Get(CtxServer).(*KlevrManager)

	txManager := NewAgentStorage()

	if agent.AgentKey == "" { // 처음 접속하는 에이전트일 경우 신규 등록
		agent.AgentKey = ch.AgentKey
		agent.GroupId = ch.ZoneID
		agent.IsActive = 1
		agent.LastAccessTime = time.Now().UTC()
		agent.Ip = paramAgent.IP
		agent.Port = paramAgent.Port
		agent.Cpu = manager.encrypt(strconv.Itoa(paramAgent.Resource.Core))
		agent.Memory = manager.encrypt(strconv.Itoa(paramAgent.Resource.Memory))
		agent.Disk = manager.encrypt(strconv.Itoa(paramAgent.Resource.Disk))
		agent.FreeMemory = manager.encrypt(strconv.Itoa(paramAgent.Resource.FreeMemory))
		agent.FreeDisk = manager.encrypt(strconv.Itoa(paramAgent.Resource.FreeDisk))
		agent.HmacKey = manager.encrypt(common.GetKey(16))
		agent.EncKey = manager.encrypt(common.GetKey(32))

		//tx.addAgent(agent)
		txManager.AddAgent(tx, agent)
	} else { // 기존에 등록된 에이전트 재접속일 경우 접속 정보 업데이트
		agent.IsActive = 1
		agent.LastAccessTime = time.Now().UTC()
		agent.Ip = paramAgent.IP
		agent.Port = paramAgent.Port
		agent.Cpu = manager.encrypt(strconv.Itoa(paramAgent.Resource.Core))
		agent.Memory = manager.encrypt(strconv.Itoa(paramAgent.Resource.Memory))
		agent.Disk = manager.encrypt(strconv.Itoa(paramAgent.Resource.Disk))
		agent.FreeMemory = manager.encrypt(strconv.Itoa(paramAgent.Resource.FreeMemory))
		agent.FreeDisk = manager.encrypt(strconv.Itoa(paramAgent.Resource.FreeDisk))
		agent.HmacKey = manager.encrypt(common.GetKey(16))
		agent.EncKey = manager.encrypt(common.GetKey(32))

		//tx.updateAgent(agent)
		txManager.UpdateAgent(tx, agent)
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
