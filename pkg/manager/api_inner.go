package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

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

type KlevrTask struct {
	ID          uint64                 `json:"id"`
	ZoneID      uint64                 `json:"zoneId"`
	Type        common.TaskType        `json:"taskType"`
	Command     string                 `json:"command"`
	Params      map[string]interface{} `json:"params"`
	CallbackURL string                 `json:"callbackUrl"`
	Log         string                 `json:"log"`
	Result      string                 `json:"result"`
	Status      string                 `json:"status"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	AgentKey    string                 `json:"agentKey"`
	ExeAgentKey string                 `json:"exeAgentKey"`
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
	registURI(inner, POST, "/tasks", addTask)
	registURI(inner, GET, "/groups/{groupID}/agents", getAgents)
	registURI(inner, GET, "/groups/{groupID}/primary", getPrimaryAgent)
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
	// ctx := CtxGetFromRequest(r)
	// var tx = GetDBConn(ctx)
	var t KlevrTask

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

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

func AddTask(tx *Tx, taskType common.TaskType, command string, zoneID uint64, agentKey string, params map[string]interface{}) *Tasks {
	b, err := json.Marshal(params)
	if err != nil {
		panic(err)
	}

	task := &Tasks{
		Type:     string(taskType),
		Command:  command,
		ZoneId:   zoneID,
		AgentKey: agentKey,
		Params: &TaskParams{
			Params: string(b),
		},
		Status: string(common.DELIVERED),
	}

	return tx.insertTask(task)
}
