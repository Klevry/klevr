package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gorhill/cronexpr"
	"github.com/gorilla/mux"
)

type serversAPI int

// InitInner initialize inner server API
func (api *API) InitInner(inner *mux.Router) {
	logger.Debug("API InitAgent - init URI")

	serversAPI := serversAPI(0)

	registURI(inner, POST, "/groups", serversAPI.addGroup)
	registURI(inner, GET, "/groups", serversAPI.getGroups)
	registURI(inner, GET, "/groups/{groupID}", serversAPI.getGroup)
	registURI(inner, DELETE, "/groups/{groupID}", serversAPI.deleteGroup)
	registURI(inner, POST, "/groups/{groupID}/apikey", serversAPI.addAPIKey)
	registURI(inner, PUT, "/groups/{groupID}/apikey", serversAPI.updateAPIKey)
	registURI(inner, GET, "/groups/{groupID}/apikey", serversAPI.getAPIKey)
	registURI(inner, GET, "/variables", serversAPI.getKlevrVariables)
	registURI(inner, GET, "/groups/{groupID}/agents", serversAPI.getAgents)
	registURI(inner, GET, "/groups/{groupID}/primary", serversAPI.getPrimaryAgent)
	registURI(inner, GET, "/groups/{groupID}/credentials", serversAPI.getCredentials)
	registURI(inner, POST, "/tasks", serversAPI.addTask)
	registURI(inner, POST, "/tasks/{groupID}/simple/inline", serversAPI.addSimpleInlineTask)
	registURI(inner, POST, "/tasks/{groupID}/simple/reserved", serversAPI.addSimpleReservedTask)
	registURI(inner, DELETE, "/tasks/{taskID}", serversAPI.cancelTask)
	registURI(inner, GET, "/tasks/{taskID}", serversAPI.getTask)
	registURI(inner, GET, "/tasks", serversAPI.getTasks)
	registURI(inner, GET, "/commands", serversAPI.getReservedCommands)
	registURI(inner, GET, "/health", serversAPI.healthCheck)
	registURI(inner, PUT, "/loglevel", serversAPI.updateLogLevel)
	registURI(inner, GET, "/loglevel", serversAPI.getLogLevel)
	registURI(inner, POST, "/credentials", serversAPI.addCredential)
	registURI(inner, GET, "/credentials/{credentialID}", serversAPI.getCredential)
	registURI(inner, DELETE, "/credentials/{credentialID}", serversAPI.deleteCredential)
	registURI(inner, GET, "/users/agents", serversAPI.getTotalAgents)
}

// addSimpleReservedTask godoc
// @Summary reserved simple TASK를 등록한다.
// @Description 간단하게 실행할 수 있는 reserved simple TASK를 등록한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/tasks/{groupID}/simple/reserved [post]
// @Param groupID path uint64 true "ZONE ID"
// @Param b body manager.SimpleReservedCommand true "TASK"
// @Success 200 {object} common.KlevrTask
func (api *serversAPI) addSimpleReservedTask(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	vars := mux.Vars(r)

	logger.Debugf("request variables : [%+v]", vars)

	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	var rc SimpleReservedCommand

	err = json.NewDecoder(r.Body).Decode(&rc)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	t := common.KlevrTask{
		ZoneID:         groupID,
		Name:           "simple",
		TaskType:       common.AtOnce,
		TotalStepCount: 1,
		Parameter:      rc.Parameter,
		Steps: []*common.KlevrTaskStep{&common.KlevrTaskStep{
			Seq:         1,
			CommandName: "simple",
			CommandType: common.RESERVED,
			Command:     rc.Command,
			IsRecover:   false,
		}},
		EventHookSendingType: common.EventHookWithAll,
	}

	// Task 상태 설정
	t = *common.TaskStatusAdd(&t)

	// DTO -> entity
	persistTask := *TaskDtoToPerist(&t)

	manager := ctx.Get(CtxServer).(*KlevrManager)

	// DB insert
	persistTask = *tx.insertTask(manager, &persistTask)

	task, _ := tx.getTask(manager, persistTask.Id)

	dto := TaskPersistToDto(task)

	b, err := json.Marshal(dto)
	if err != nil {
		panic(err)
	}

	logger.Debugf("response : [%s]", string(b))

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

// addSimpleInlineTask godoc
// @Summary inline simple TASK를 등록한다.
// @Description 간단하게 실행할 수 있는 inline script 형태의 simple TASK를 등록한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/tasks/{groupID}/simple/inline [post]
// @Param groupID path uint64 true "ZONE ID"
// @Param b body string true "inline script"
// @Success 200 {object} common.KlevrTask
func (api *serversAPI) addSimpleInlineTask(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	vars := mux.Vars(r)

	logger.Debugf("request variables : [%+v]", vars)

	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	nr := &common.Request{Request: r}

	t := common.KlevrTask{
		ZoneID:         groupID,
		Name:           "simple",
		TaskType:       common.AtOnce,
		TotalStepCount: 1,
		Steps: []*common.KlevrTaskStep{&common.KlevrTaskStep{
			Seq:         1,
			CommandName: "simple",
			CommandType: common.INLINE,
			Command:     nr.BodyToString(),
			IsRecover:   false,
		}},
		EventHookSendingType: common.EventHookWithAll,
	}

	// Task 상태 설정
	t = *common.TaskStatusAdd(&t)

	// DTO -> entity
	persistTask := *TaskDtoToPerist(&t)

	manager := ctx.Get(CtxServer).(*KlevrManager)

	// DB insert
	persistTask = *tx.insertTask(manager, &persistTask)

	task, _ := tx.getTask(manager, persistTask.Id)

	dto := TaskPersistToDto(task)

	b, err := json.Marshal(dto)
	if err != nil {
		panic(err)
	}

	logger.Debugf("response : [%s]", string(b))

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

// getPrimaryAgent godoc
// @Summary primary agent 정보를 반환한다.
// @Description groupID에 해당하는 klevr zone의 primary agent 정보를 반환한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/groups/{groupID}/primary [get]
// @Param groupID path uint64 true "ZONE ID"
// @Success 200 {object} common.Agent
func (api *serversAPI) getPrimaryAgent(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	vars := mux.Vars(r)

	logger.Debugf("request variables : [%+v]", vars)

	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	primary, exist := tx.getPrimaryAgent(groupID)
	var agent common.Agent

	if exist {
		//a := tx.getAgentByID(primary.AgentId)
		txManager := NewAgentStorage()
		a := txManager.GetAgentByID(tx, groupID, primary.AgentId)

		agent = common.Agent{
			AgentKey:           a.AgentKey,
			IsActive:           byteToBool(a.IsActive),
			LastAliveCheckTime: &common.JSONTime{a.LastAliveCheckTime},
			IP:                 a.Ip,
			Port:               a.Port,
			Version:            a.Version,
			Resource:           &common.Resource{},
		}

		manager := ctx.Get(CtxServer).(*KlevrManager)

		core, _ := strconv.Atoi(manager.decrypt(a.Cpu))
		memory, _ := strconv.Atoi(manager.decrypt(a.Memory))
		disk, _ := strconv.Atoi(manager.decrypt(a.Disk))
		freeMemory, _ := strconv.Atoi(manager.decrypt(a.FreeMemory))
		freeDisk, _ := strconv.Atoi(manager.decrypt(a.FreeDisk))

		agent.Core = core
		agent.Memory = memory
		agent.Disk = disk
		agent.FreeMemory = freeMemory
		agent.FreeDisk = freeDisk
	}

	b, err := json.Marshal(agent)
	if err != nil {
		panic(err)
	}

	logger.Debugf("response : [%s]", string(b))

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

// getReservedCommands godoc
// @Summary 예약어 커맨드 정보를 반환한다.
// @Description Klevr에서 사용할 수 있는 예약어 커맨드 정보를 반환한다. 사용자는 이 정보를 토대로 task를 생성하여 요청할 수 있다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/commands [get]
// @Success 200 {object} []ReservedCommand
func (api *serversAPI) getReservedCommands(w http.ResponseWriter, r *http.Request) {
	m := make(map[string]ReservedCommand)

	for k, v := range common.GetCommands() {
		m[k] = ReservedCommand{
			Description:    v.Description,
			ParameterModel: v.ParameterModel,
			ResultModel:    v.ResultModel,
			HasRecover:     v.Recover != nil,
		}
	}

	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

// healthCheck godoc
// @Summary klevr manager 확인용
// @Description klevr manager 상태 체크
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/health [get]
// @Success 200 {object} string "{\"health\":ok}"
func (api *serversAPI) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "{\"health\":\"ok\"}")
}

// updateLogLevel godoc
// @Summary klevr manager 로그 레벨
// @Description klevr manager의 로그 레벨 변경
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/loglevel [put]
// @Param b body string true "Log Level(debug, info, warn, error, fatal)"
// @Success 200 {object} string "{\"updated\":true|false}"
func (api *serversAPI) updateLogLevel(w http.ResponseWriter, r *http.Request) {
	nr := &common.Request{Request: r}
	targetLevel := nr.BodyToString()
	var level logger.Level

	switch strings.ToLower(targetLevel) {
	case "debug":
		level = 0
	case "info":
		level = 1
	case "warn", "warning":
		level = 2
	case "error":
		level = 3
	case "fatal":
		level = 4
	}

	logger.SetLevel(level)

	w.WriteHeader(200)
	fmt.Fprintf(w, "{\"updated\":ok}")
}

// getLogLevel godoc
// @Summary klevr manager 로그 레벨
// @Description klevr manager의 현재 로그 레벨
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/loglevel [get]
// @Success 200 {object} string"
func (api *serversAPI) getLogLevel(w http.ResponseWriter, r *http.Request) {
	level := logger.GetLevel()

	var levelValue string
	switch int(level) {
	case 0:
		levelValue = "debug"
	case 1:
		levelValue = "info"
	case 2:
		levelValue = "warn"
	case 3:
		levelValue = "error"
	case 4:
		levelValue = "fatal"
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, levelValue)
}

// getAgents godoc
// @Summary zone의 agent 목록을 반환한다.
// @Description groupID에 해당하는 klevr zone의 모든 agent 정보를 반환한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/groups/{groupID}/agents [get]
// @Param groupID path uint64 true "ZONE ID"
// @Success 200 {object} []common.Agent
func (api *serversAPI) getAgents(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	vars := mux.Vars(r)

	logger.Debugf("request variables : [%+v]", vars)

	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	//cnt, agents := tx.getAgentsByGroupId(groupID)
	txManager := NewAgentStorage()
	cnt, agents := txManager.GetAgentsByZoneID(tx, groupID)

	nodes := make([]Agent, cnt)

	manager := ctx.Get(CtxServer).(*KlevrManager)

	if cnt > 0 {
		for i, a := range *agents {
			nodes[i] = Agent{
				AgentKey:           a.AgentKey,
				IsActive:           byteToBool(a.IsActive),
				LastAliveCheckTime: &common.JSONTime{a.LastAliveCheckTime},
				LastAccessTime:     &common.JSONTime{a.LastAccessTime},
				IP:                 a.Ip,
				Port:               a.Port,
				Version:            a.Version,
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
	}

	b, err := json.Marshal(nodes)
	if err != nil {
		panic(err)
	}

	logger.Debugf("response : [%s]", string(b))

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

// getTotalAgents godoc
// @Summary 전체 agent 목록을 반환한다.
// @Description 모든 agent 정보를 반환한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/users/agents [get]
// @Success 200 {object} []common.Agent
func (api *serversAPI) getTotalAgents(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	cnt, agents := tx.getTotalAgents()
	nodes := make([]Agent, cnt)

	manager := ctx.Get(CtxServer).(*KlevrManager)

	if cnt > 0 {
		for i, a := range *agents {
			nodes[i] = Agent{
				AgentKey:           a.AgentKey,
				IsActive:           byteToBool(a.IsActive),
				LastAliveCheckTime: &common.JSONTime{a.LastAliveCheckTime},
				LastAccessTime:     &common.JSONTime{a.LastAccessTime},
				IP:                 a.Ip,
				Port:               a.Port,
				Version:            a.Version,
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
	}

	b, err := json.Marshal(nodes)
	if err != nil {
		panic(err)
	}

	logger.Debugf("response : [%s]", string(b))

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

// addTask godoc
// @Summary TASK를 등록한다.
// @Description KlevrTask 모델에 기입된 ZONE의 AGENT에서 실행할 TASK를 등록한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/tasks [post]
// @Param b body common.KlevrTask true "TASK"
// @Success 200 {object} common.KlevrTask
func (api *serversAPI) addTask(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	var tx = GetDBConn(ctx)
	var t common.KlevrTask

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	logger.Debugf("request add task : [%+v]", t)

	// Task 상태 설정
	t = *common.TaskStatusAdd(&t)

	// Task validation
	if common.Iteration == t.TaskType {
		_, err := cronexpr.Parse(t.Cron)

		if err != nil {
			common.WriteHTTPError(400, w, err, "Invalid cron expression - "+t.Cron)
		}
	}

	if t.EventHookSendingType == "" {
		t.EventHookSendingType = common.EventHookWithAll
	}

	// DTO -> entity
	persistTask := *TaskDtoToPerist(&t)

	manager := ctx.Get(CtxServer).(*KlevrManager)

	// DB insert
	persistTask = *tx.insertTask(manager, &persistTask)

	task, _ := tx.getTask(manager, persistTask.Id)

	dto := TaskPersistToDto(task)

	b, err := json.Marshal(dto)
	if err != nil {
		panic(err)
	}

	logger.Debugf("response : [%s]", string(b))

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

// removeTask godoc
// @Summary TASK를 취소한다.
// @Description agent에 전달되지 않은(hand-over 이전) task를 취소한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/tasks/{taskID} [delete]
// @Param taskID path uint64 true "task id"
// @Success 200 {object} string "{\"canceld\":true/false}"
func (api *serversAPI) cancelTask(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	var tx = GetDBConn(ctx)

	vars := mux.Vars(r)

	taskID, err := strconv.ParseUint(vars["taskID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid task id : %+v", vars["taskID"]))
		return
	}

	canceled := tx.cancelTask(taskID)

	// task cacel이 성공하면
	if canceled {
		// task 부가 데이터 삭제
		tx.purgeTask(taskID)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "{\"canceled\":%v}", canceled)
}

// getTask godoc
// @Summary TASK를 조회한다.
// @Description taskID에 해당하는 TASK를 조회한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/tasks/{taskID} [get]
// @Param taskID path uint64 true "task id"
// @Success 200 {object} common.KlevrTask
func (api *serversAPI) getTask(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	var tx = GetDBConn(ctx)

	vars := mux.Vars(r)

	taskID, err := strconv.ParseUint(vars["taskID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid task id : %+v", vars["taskID"]))
		return
	}

	manager := ctx.Get(CtxServer).(*KlevrManager)

	task, exist := tx.getTask(manager, taskID)

	if exist {
		dto := TaskPersistToDto(task)

		b, err := json.Marshal(dto)
		if err != nil {
			panic(err)
		}

		logger.Debugf("response : [%s]", string(b))

		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", b)
	} else {
		common.NewHTTPError(404, fmt.Sprintf("Not exist task for ID - %d", taskID))
	}
}

// getTasks godoc
// @Summary TASK 목록을 반환한다.
// @Description 검색조건에 해당하는 TASK 목록을 반환한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/tasks [get]
// @Param groupID query []uint64 true "ZONE ID 배열"
// @Param status query []string false "STATUS 배열"
// @Param agentKey query []string false "AGENT KEY 배열"
// @Param name query []string false "TASK NAME 배열"
// @Success 200 {object} []common.KlevrTask
func (api *serversAPI) getTasks(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	var tx = GetDBConn(ctx)

	groupIDs := r.URL.Query()["groupID"]
	statuses := r.URL.Query()["status"]
	agentKeys := r.URL.Query()["agentKey"]
	taskNames := r.URL.Query()["name"]

	logger.Debugf("request URI - [%s]", r.RequestURI)
	logger.Debugf("%+v", groupIDs)
	logger.Debugf("%+v", statuses)
	logger.Debugf("%+v", agentKeys)
	logger.Debugf("%+v", taskNames)

	logger.Debugf("%d", len(groupIDs))
	logger.Debugf("%d", len(statuses))
	logger.Debugf("%d", len(agentKeys))
	logger.Debugf("%d", len(taskNames))

	if groupIDs == nil || len(groupIDs) == 0 {
		common.WriteHTTPError(400, w, nil, "Query parameter groupID is required.")
		return
	}

	var iGroupIDs []uint64 = make([]uint64, len(groupIDs))
	var err error

	for i, id := range groupIDs {
		iGroupIDs[i], err = strconv.ParseUint(id, 0, 64)
		if err != nil {
			common.WriteHTTPError(400, w, err, fmt.Sprintf("invalid groupID - [%s]", id))
			return
		}
	}

	tasks, exist := tx.getTasks(iGroupIDs, statuses, agentKeys, taskNames)

	var b []byte

	if exist {
		var dtos []common.KlevrTask = make([]common.KlevrTask, len(*tasks))

		for i, t := range *tasks {
			dtos[i] = *TaskPersistToDto(&t)
		}

		b, err = json.Marshal(dtos)
		if err != nil {
			panic(err)
		}

		logger.Debugf("response : [%s]", string(b))
	}

	// fmt.Fprintf(w, "%s", b)
	w.Write(b)
	w.WriteHeader(200)

	logger.Debugf("response : [%s]", string(b))
}

// getKlevrVariables godoc
// @Summary Klevr에서 제공하는 시스템 변수 목록을 조회한다.
// @Description TASK inline command에서 사용할 수 있는 시스템 변수 목록을 조회한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/variables [get]
// @Success 200 {object} []KlevrVariable
func (api *serversAPI) getKlevrVariables(w http.ResponseWriter, r *http.Request) {
	// ctx := CtxGetFromRequest(r)

	var variables []KlevrVariable = []KlevrVariable{
		KlevrVariable{
			Name:        "KLEVR.HOST",
			Type:        "string",
			Length:      "-",
			Description: "klevr host url",
			Example:     "echo {KLEVR.HOST} => echo http://klevr.io",
		},
		KlevrVariable{
			Name:        "KLEVR.PORT",
			Type:        "int",
			Length:      "-",
			Description: "klevr service port",
			Example:     "echo {KLEVR.HOST} => echo http://klevr.io",
		},
	}

	b, err := json.Marshal(&variables)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

// addAPIKey godoc
// @Summary 사용자 그룹에 API key를 등록한다.
// @Description agent가 zone에 접속할 수 있는 API KEY를 등록한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/groups/{groupID}/apikey [post]
// @Param groupID path uint64 true "ZONE ID"
// @Param b body string true "API KEY"
// @Success 200
func (api *serversAPI) addAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	nr := &common.Request{Request: r}

	vars := mux.Vars(r)
	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	manager := ctx.Get(CtxServer).(*KlevrManager)

	auth := &ApiAuthentications{
		ApiKey:  manager.encrypt(nr.BodyToString()),
		GroupId: groupID,
	}

	tx.addAPIKey(auth)

	w.WriteHeader(200)
}

// updateAPIKey godoc
// @Summary 사용자 그룹의 API key를 수정한다.
// @Description agent가 zone에 접속할 수 있는 API KEY를 수정한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/groups/{groupID}/apikey [put]
// @Param groupID path uint64 true "ZONE ID"
// @Param b body string true "API KEY"
// @Success 200
func (api *serversAPI) updateAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	nr := &common.Request{Request: r}

	vars := mux.Vars(r)
	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	manager := ctx.Get(CtxServer).(*KlevrManager)
	apiKey := nr.BodyToString()

	auth := &ApiAuthentications{
		ApiKey:  manager.encrypt(apiKey),
		GroupId: groupID,
	}

	tx.updateAPIKey(auth)

	ctxAPI := ctx.Get(CtxAPI).(*API)
	ctxAPI.APIKeyMap.Set(strconv.FormatUint(groupID, 10), apiKey)

	w.WriteHeader(200)
}

// getAPIKey godoc
// @Summary 사용자 그룹의 API key를 조회한다.
// @Description agent가 zone에 접속할 수 있는 API KEY를 조회한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/groups/{groupID}/apikey [get]
// @Param groupID path uint64 true "ZONE ID"
// @Success 200 {object} string
func (api *serversAPI) getAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	vars := mux.Vars(r)
	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	auth, exist := tx.getAPIKey(groupID)
	if !exist {
		common.WriteHTTPError(400, w, nil, fmt.Sprintf("Does not exist APIKey for groupId : %d", groupID))
		return
	}

	manager := ctx.Get(CtxServer).(*KlevrManager)

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", manager.decrypt(auth.ApiKey))
}

// addGroup godoc
// @Summary ZONE을 추가한다.
// @Description KLEVR ZONE을 생성한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/groups [post]
// @Param b body AgentGroups true "AgentGroups model"
// @Success 200 {object} AgentGroups
func (api *serversAPI) addGroup(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)
	var ag AgentGroups

	err := json.NewDecoder(r.Body).Decode(&ag)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	logger.Debugf("request AgentGroup : %+v", &ag)
	// logger.Debug("%v", time.Now().UTC())

	tx.addAgentGroup(&ag)

	logger.Debugf("response AgentGroup : %+v", &ag)

	b, err := json.Marshal(&ag)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

// getGroups godoc
// @Summary ZONE 목록을 조회한다.
// @Description KLEVR ZONE 목록을 조회한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/groups [get]
// @Success 200 {object} []AgentGroups
func (api *serversAPI) getGroups(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	var tx = GetDBConn(ctx)

	ags := tx.getAgentGroups()

	b, err := json.Marshal(&ags)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

// getGroup godoc
// @Summary ZONE을 조회한다.
// @Description KLEVR ZONE을 조회한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/groups/{groupID} [get]
// @Param groupID path uint64 true "ZONE ID"
// @Success 200 {object} AgentGroups
func (api *serversAPI) getGroup(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	var tx = GetDBConn(ctx)

	vars := mux.Vars(r)
	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	ag, exist := tx.getAgentGroup(groupID)
	if !exist {
		common.WriteHTTPError(400, w, nil, fmt.Sprintf("Does not exist zone for groupId : %d", groupID))
		return
	}

	b, err := json.Marshal(&ag)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

// deleteGroup godoc
// @Summary ZONE을 삭제한다.
// @Description KLEVR ZONE을 삭제한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/groups/{groupID} [delete]
// @Param groupID path uint64 true "ZONE ID"
// @Success 200 {object} string "{\"deleted\":true/false}"
func (api *serversAPI) deleteGroup(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	qryOperation := r.URL.Query()["op"]
	if qryOperation != nil && qryOperation[0] == "db" {
		logger.Debug("db balue was enterted as an option.")
	}

	vars := mux.Vars(r)
	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid zone id : %+v", vars["groupID"]))
		return
	}

	ctxAPI := ctx.Get(CtxAPI).(*API)
	ctxAPI.APIKeyMap.Remove(strconv.FormatUint(groupID, 10))

	// logger.Debug("%v", time.Now().UTC())
	err = api.deletegroup(tx, groupID)
	if err != nil {
		common.WriteHTTPError(500, w, err, err.Error())
		return
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "{\"deleted\":%v}", true)
}

func (api *serversAPI) deletegroup(tx *Tx, id uint64) error {
	tx.deletePrimaryAgent(id)
	_, ok := tx.getPrimaryAgent(id)
	if ok == true {
		return fmt.Errorf("It cannot remove the zone(primaryagent) of the zoneid: %d", id)
	}

	tx.deleteApiAuthentication(id)
	cnt, _ := tx.getApiAuthenticationsByGroupId(id)
	if cnt > 0 {
		return fmt.Errorf("It cannot remove the zone(apiauthentication) of the zoneid: %d", id)
	}

	//tx.deleteAgent(id)
	//cnt, _ = tx.getAgentsByGroupId(id)
	txManager := NewAgentStorage()
	txManager.DeleteAgent(tx, id)
	cnt, _ = txManager.GetAgentsByZoneID(tx, id)
	if cnt > 0 {
		return fmt.Errorf("It cannot remove the zone of the zoneid: %d", id)
	}

	tx.deleteAgentGroup(id)

	return nil
}

func TaskDtoToPerist(dto *common.KlevrTask) *Tasks {
	persist := &Tasks{
		Id:          dto.ID,
		ZoneId:      dto.ZoneID,
		Name:        dto.Name,
		TaskType:    dto.TaskType,
		Schedule:    dto.Schedule.Time,
		AgentKey:    dto.AgentKey,
		ExeAgentKey: dto.ExeAgentKey,
		Status:      dto.Status,
		TaskDetail: &TaskDetail{
			TaskId:             dto.ID,
			Cron:               dto.Cron,
			UntilRun:           dto.UntilRun.Time,
			Timeout:            dto.Timeout,
			ExeAgentChangeable: dto.ExeAgentChangeable,
			TotalStepCount:     dto.TotalStepCount,
			CurrentStep:        dto.CurrentStep,
			HasRecover:         dto.HasRecover,
			Parameter:          dto.Parameter,
			CallbackUrl:        dto.CallbackURL,
			Result:             dto.Result,
			FailedStep:         dto.FailedStep,
			IsFailedRecover:    dto.IsFailedRecover,
			ShowLog:            dto.ShowLog,
		},
	}

	stepLen := len(dto.Steps)

	if stepLen > 0 {
		steps := make([]TaskSteps, stepLen)

		for i, dtoStep := range dto.Steps {
			steps[i] = TaskSteps{
				Id:          dtoStep.ID,
				Seq:         dtoStep.Seq,
				TaskId:      dto.ID,
				CommandName: dtoStep.CommandName,
				CommandType: dtoStep.CommandType,
				IsRecover:   dtoStep.IsRecover,
			}

			if dtoStep.CommandType == common.RESERVED {
				steps[i].ReservedCommand = dtoStep.Command
			} else if dtoStep.CommandType == common.INLINE {
				steps[i].InlineScript = dtoStep.Command
			} else {
				panic(fmt.Sprintf("Invalid Task Step CommandType : [%s]", dtoStep.CommandType))
			}
		}

		persist.TaskSteps = &steps
	}

	if dto.Log != "" {
		persist.Logs = &TaskLogs{
			TaskId: persist.Id,
			Logs:   dto.Log,
		}
	}

	logger.Debugf("TaskDtoToPerist \ndto : [%+v]\npersist : [%+v]", dto, persist)

	return persist
}

func TaskPersistToDto(persist *Tasks) *common.KlevrTask {
	detail := persist.TaskDetail

	logger.Debugf("detail [%+v]", detail)

	dto := &common.KlevrTask{
		ID:          persist.Id,
		ZoneID:      persist.ZoneId,
		Name:        persist.Name,
		TaskType:    persist.TaskType,
		Schedule:    common.JSONTime{Time: persist.Schedule},
		AgentKey:    persist.AgentKey,
		ExeAgentKey: persist.ExeAgentKey,
		Status:      persist.Status,
		CreatedAt:   common.JSONTime{Time: persist.CreatedAt},
		UpdatedAt:   common.JSONTime{Time: persist.UpdatedAt},
	}

	if detail != nil {
		dto.Cron = detail.Cron
		dto.UntilRun = common.JSONTime{Time: detail.UntilRun}
		dto.Timeout = detail.Timeout
		dto.ExeAgentChangeable = detail.ExeAgentChangeable
		dto.TotalStepCount = detail.TotalStepCount
		dto.CurrentStep = detail.CurrentStep
		dto.HasRecover = detail.HasRecover
		dto.Parameter = detail.Parameter
		dto.CallbackURL = detail.CallbackUrl
		dto.Result = detail.Result
		dto.FailedStep = detail.FailedStep
		dto.IsFailedRecover = detail.IsFailedRecover
		dto.ShowLog = detail.ShowLog

		logger.Debugf("dto cron : [%s], detail cron : [%s]", dto.Cron, detail.Cron)
	}

	stepLen := 0
	if persist.TaskSteps != nil {
		stepLen = len(*persist.TaskSteps)
	}

	if stepLen > 0 {
		steps := make([]*common.KlevrTaskStep, stepLen)

		for i, step := range *persist.TaskSteps {
			steps[i] = &common.KlevrTaskStep{
				ID:          step.Id,
				Seq:         step.Seq,
				CommandName: step.CommandName,
				CommandType: step.CommandType,
				IsRecover:   step.IsRecover,
			}

			if step.CommandType == common.RESERVED {
				steps[i].Command = step.ReservedCommand
			} else if step.CommandType == common.INLINE {
				steps[i].Command = step.InlineScript
			} else {
				panic(fmt.Sprintf("Invalid Task Step CommandType : [%s]", step.CommandType))
			}

			dto.Steps = steps
		}
	}

	if persist.Logs != nil {
		dto.Log = persist.Logs.Logs
	}

	logger.Debugf("TaskPersistToDto \npersist : [%+v]\ndto : [%+v]", persist, dto)

	return dto
}

func TaskMatchingCredential(manager *KlevrManager, task Tasks, credential *[]Credentials) Tasks {
	if len(*credential) == 0 {
		return task
	}

	if len(task.TaskDetail.Parameter) == 0 {
		return task
	}

	r := regexp.MustCompile("{{2}[a-zA-Z0-9]*}{2}")
	isMatch := r.MatchString(task.TaskDetail.Parameter)
	if isMatch == false {
		return task
	}

	for _, c := range *credential {
		pattern := fmt.Sprintf("{{2}%s}{2}", c.Name)
		v := manager.decrypt(c.Value)

		re := regexp.MustCompile(pattern)
		task.TaskDetail.Parameter = fmt.Sprintf("%s", re.ReplaceAllString(task.TaskDetail.Parameter, v))
		logger.Debugf("Apply Credential : %s", task.TaskDetail.Parameter)
	}

	return task
}

// addCredential godoc
// @Summary Credential을 등록한다.
// @Description KlevrCredential 모델에 기입된 ZONE에서 사용할 Credential을 등록한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/credentials [post]
// @Param b body common.KlevrCredential true "Credential"
// @Success 200 {object} common.KlevrCredential
func (api *serversAPI) addCredential(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	var tx = GetDBConn(ctx)
	var c common.KlevrCredential

	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	logger.Debugf("request add credential : [%+v]", c)

	// DTO -> entity
	persistCredential := *CredentialDtoToPerist(&c)

	manager := ctx.Get(CtxServer).(*KlevrManager)

	// DB insert
	persistCredential = *tx.insertCredential(manager, &persistCredential)

	task, _ := tx.getTask(manager, persistCredential.Id)

	dto := TaskPersistToDto(task)

	b, err := json.Marshal(dto)
	if err != nil {
		panic(err)
	}

	logger.Debugf("response : [%s]", string(b))

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

// getCredential godoc
// @Summary Credential를 조회한다.
// @Description credentialID에 해당하는 Credential를 조회한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/credentials/{credentialID} [get]
// @Param credentialID path uint64 true "credential id"
// @Success 200 {object} common.KlevrCredential
func (api *serversAPI) getCredential(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	var tx = GetDBConn(ctx)

	vars := mux.Vars(r)

	credentialID, err := strconv.ParseUint(vars["credentialID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid credential id : %+v", vars["credentialID"]))
		return
	}

	manager := ctx.Get(CtxServer).(*KlevrManager)

	credential, exist := tx.getCredential(manager, credentialID)

	if exist {
		dto := CredentialPersistToDto(credential)

		b, err := json.Marshal(dto)
		if err != nil {
			panic(err)
		}

		logger.Debugf("response : [%s]", string(b))

		w.WriteHeader(200)
		fmt.Fprintf(w, "%s", b)
	} else {
		common.NewHTTPError(404, fmt.Sprintf("Not exist task for ID - %d", credentialID))
	}
}

// getCredentials godoc
// @Summary Credential 목록을 반환한다.
// @Description 검색조건에 해당하는 Credential 목록을 반환한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/groups/{groupID}/credentials [get]
// @Param groupID path uint64 true "ZONE ID"
// @Success 200 {object} []common.KlevrCredential
func (api *serversAPI) getCredentials(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	var tx = GetDBConn(ctx)

	vars := mux.Vars(r)
	logger.Debugf("request variables: [%+v]", vars)

	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id: %+v", vars["groupID"]))
		return
	}

	credentials, cnt := tx.getCredentials(groupID)

	var b []byte

	if cnt > 0 {
		var dtos []common.KlevrCredential = make([]common.KlevrCredential, len(*credentials))

		for i, c := range *credentials {
			dtos[i] = *CredentialPersistToDto(&c)
		}

		b, err = json.Marshal(dtos)
		if err != nil {
			panic(err)
		}

		logger.Debugf("response : [%s]", string(b))
	}

	// fmt.Fprintf(w, "%s", b)
	w.Write(b)
	w.WriteHeader(200)

	logger.Debugf("response : [%s]", string(b))
}

// deleteCredential godoc
// @Summary Credential을 삭제한다.
// @Description credentialID에 해당하는 credential을 삭제한다.
// @Tags servers
// @Accept json
// @Produce json
// @Router /inner/credentials/{credentialID} [delete]
// @Param credentialID path uint64 true "credential id"
// @Success 200 {object} string "{\"deleted\":true}"
func (api *serversAPI) deleteCredential(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	var tx = GetDBConn(ctx)

	vars := mux.Vars(r)

	credentialID, err := strconv.ParseUint(vars["credentialID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid credential id : %+v", vars["credentialID"]))
		return
	}

	tx.deleteCredential(credentialID)

	w.WriteHeader(200)
	fmt.Fprint(w, "{\"deletedd\": true}")
}

func CredentialDtoToPerist(dto *common.KlevrCredential) *Credentials {
	persist := &Credentials{
		Id:     dto.ID,
		ZoneId: dto.ZoneID,
		Name:   dto.Name,
		Value:  dto.Value,
	}

	logger.Debugf("CredentialDtoToPerist \ndto : [%+v]\npersist : [%+v]", dto, persist)

	return persist
}

func CredentialPersistToDto(persist *Credentials) *common.KlevrCredential {
	dto := &common.KlevrCredential{
		ID:        persist.Id,
		ZoneID:    persist.ZoneId,
		Name:      persist.Name,
		Value:     persist.Value,
		CreatedAt: common.JSONTime{Time: persist.CreatedAt},
		UpdatedAt: common.JSONTime{Time: persist.UpdatedAt},
	}

	logger.Debugf("CredentialPersistToDto \npersist : [%+v]\ndto : [%+v]", persist, dto)

	return dto
}
