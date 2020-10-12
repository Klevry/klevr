package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gorilla/mux"
)

type KlevrVariable struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Length      string `json:"length"`
	Description string `json:"description"`
	Example     string `json:"example"`
}

// InitInner initialize inner server API
func (api *API) InitInner(inner *mux.Router) {
	logger.Debug("API InitAgent - init URI")

	registURI(inner, POST, "/groups", addGroup)
	registURI(inner, GET, "/groups", getGroups)
	registURI(inner, GET, "/groups/{groupID}", getGroup)
	registURI(inner, POST, "/groups/{groupID}/apikey", addAPIKey)
	registURI(inner, PUT, "/groups/{groupID}/apikey", updateAPIKey)
	registURI(inner, GET, "/groups/{groupID}/apikey", getAPIKey)
	registURI(inner, GET, "/variables", getKlevrVariables)
	registURI(inner, GET, "/groups/{groupID}/agents", getAgents)
	registURI(inner, GET, "/groups/{groupID}/primary", getPrimaryAgent)
	registURI(inner, POST, "/tasks", addTask)
	registURI(inner, DELETE, "/tasks/{taskID}", removeTask)
	registURI(inner, GET, "/tasks/{taskID}", getTask)
	registURI(inner, GET, "/tasks", getTasks)
}

func getPrimaryAgent(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	vars := mux.Vars(r)

	logger.Debugf("request variables : [%+v]", vars)

	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	primary := tx.getPrimaryAgent(groupID)
	var agent common.Agent

	if primary != nil {
		a := tx.getAgentByID(primary.AgentId)

		agent = common.Agent{
			AgentKey:           a.AgentKey,
			IsActive:           a.IsActive,
			LastAliveCheckTime: a.LastAliveCheckTime.Unix(),
			IP:                 a.Ip,
			Port:               a.Port,
			Version:            a.Version,
			Resource: &common.Resource{
				Core:   a.Cpu,
				Memory: a.Memory,
				Disk:   a.Disk,
			},
		}
	}

	b, err := json.Marshal(agent)
	if err != nil {
		panic(err)
	}

	logger.Debugf("response : [%s]", string(b))

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

func getAgents(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	vars := mux.Vars(r)

	logger.Debugf("request variables : [%+v]", vars)

	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	cnt, agents := tx.getAgentsByGroupId(groupID)
	nodes := make([]common.Agent, cnt)

	if cnt > 0 {
		for i, a := range *agents {
			nodes[i] = common.Agent{
				AgentKey:           a.AgentKey,
				IsActive:           a.IsActive,
				LastAliveCheckTime: a.LastAliveCheckTime.Unix(),
				IP:                 a.Ip,
				Port:               a.Port,
				Version:            a.Version,
				Resource: &common.Resource{
					Core:   a.Cpu,
					Memory: a.Memory,
					Disk:   a.Disk,
				},
			}
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

func addTask(w http.ResponseWriter, r *http.Request) {
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

	// DTO -> entity
	persistTask := *TaskDtoToPerist(&t)

	// DB insert
	persistTask = *tx.insertTask(&persistTask)

	task, _ := tx.getTask(persistTask.Id)

	dto := TaskPersistToDto(task)

	b, err := json.Marshal(dto)
	if err != nil {
		panic(err)
	}

	logger.Debugf("response : [%s]", string(b))

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

func removeTask(w http.ResponseWriter, r *http.Request) {
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

func getTask(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	var tx = GetDBConn(ctx)

	vars := mux.Vars(r)

	taskID, err := strconv.ParseUint(vars["taskID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid task id : %+v", vars["taskID"]))
		return
	}

	task, exist := tx.getTask(taskID)

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

func getTasks(w http.ResponseWriter, r *http.Request) {
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

		b, err := json.Marshal(dtos)
		if err != nil {
			panic(err)
		}

		logger.Debugf("response : [%s]", string(b))
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

func getKlevrVariables(w http.ResponseWriter, r *http.Request) {
	var variables []KlevrVariable

	variables = append(variables, KlevrVariable{
		Name:        "KLEVR_HOST",
		Type:        "string",
		Length:      "-",
		Description: "klevr host url",
		Example:     "echo {KLEVR_HOST} => echo http://klevr.io",
	})

	b, err := json.Marshal(&variables)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

func addAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	nr := &common.Request{Request: r}

	vars := mux.Vars(r)
	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	auth := &ApiAuthentications{
		ApiKey:  nr.BodyToString(),
		GroupId: groupID,
	}

	tx.addAPIKey(auth)

	w.WriteHeader(200)
}

func updateAPIKey(w http.ResponseWriter, r *http.Request) {
	ctx := CtxGetFromRequest(r)
	tx := GetDBConn(ctx)

	nr := &common.Request{Request: r}

	vars := mux.Vars(r)
	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	auth := &ApiAuthentications{
		ApiKey:  nr.BodyToString(),
		GroupId: groupID,
	}

	tx.updateAPIKey(auth)

	w.WriteHeader(200)
}

func getAPIKey(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", auth.ApiKey)
}

func addGroup(w http.ResponseWriter, r *http.Request) {
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

func getGroups(w http.ResponseWriter, r *http.Request) {
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

func getGroup(w http.ResponseWriter, r *http.Request) {
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

// func AddTask(tx *Tx, taskType common.TaskType, command string, zoneID uint64, agentKey string, params map[string]interface{}) *Tasks {
// 	b, err := json.Marshal(params)
// 	if err != nil {
// 		panic(err)
// 	}

// 	task := &Tasks{
// 		TaskType: taskType,
// 		Command:  command,
// 		ZoneId:   zoneID,
// 		AgentKey: agentKey,
// 		Params: &TaskParams{
// 			Params: string(b),
// 		},
// 		Status: string(common.DELIVERED),
// 	}

// 	return tx.insertTask(task)
// }

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
