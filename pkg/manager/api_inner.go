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

// InitInner initialize inner server API
func (api *API) InitInner(inner *mux.Router) {
	logger.Debug("API InitAgent - init URI")

	// registURI(agent, PUT, "/handshake", api.receiveHandshake)
	// registURI(agent, PUT, "/:agentKey", api.receivePolling)

	registURI(inner, POST, "/groups", api.addGroup)
	registURI(inner, GET, "/groups", api.getGroups)
	registURI(inner, GET, "/groups/{groupID}", api.getGroup)
	registURI(inner, POST, "/groups/{groupID}/apikey", api.addAPIKey)
	registURI(inner, PUT, "/groups/{groupID}/apikey", api.updateAPIKey)
	registURI(inner, GET, "/groups/{groupID}/apikey", api.getAPIKey)
}

func (api *API) addAPIKey(w http.ResponseWriter, r *http.Request) {
	var conn = GetDBConn(r)

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

	addAPIKey(conn, auth)

	w.WriteHeader(200)
}

func (api *API) updateAPIKey(w http.ResponseWriter, r *http.Request) {
	var conn = GetDBConn(r)

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

	updateAPIKey(conn, auth)

	w.WriteHeader(200)
}

func (api *API) getAPIKey(w http.ResponseWriter, r *http.Request) {
	var conn = GetDBConn(r)

	vars := mux.Vars(r)
	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	auth, exist := getAPIKey(conn, groupID)
	if !exist {
		common.WriteHTTPError(400, w, nil, fmt.Sprintf("Does not exist APIKey for groupId : %d", groupID))
		return
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", auth.ApiKey)
}

func (api *API) addGroup(w http.ResponseWriter, r *http.Request) {
	var conn = GetDBConn(r)
	var ag AgentGroups

	err := json.NewDecoder(r.Body).Decode(&ag)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	logger.Debugf("request AgentGroup : %+v", &ag)
	// logger.Debug("%v", time.Now().UTC())

	addAgentGroup(conn, &ag)

	logger.Debugf("response AgentGroup : %+v", &ag)

	b, err := json.Marshal(&ag)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

func (api *API) getGroups(w http.ResponseWriter, r *http.Request) {
	var conn = GetDBConn(r)

	ags := getAgentGroups(conn)

	b, err := json.Marshal(&ags)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", b)
}

func (api *API) getGroup(w http.ResponseWriter, r *http.Request) {
	var conn = GetDBConn(r)

	vars := mux.Vars(r)
	groupID, err := strconv.ParseUint(vars["groupID"], 10, 64)
	if err != nil {
		common.WriteHTTPError(500, w, err, fmt.Sprintf("Invalid group id : %+v", vars["groupID"]))
		return
	}

	ag, exist := getAgentGroup(conn, groupID)
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
