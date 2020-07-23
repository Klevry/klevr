package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gorilla/mux"
)

// InitAgent initialize agent API
func (api *API) InitAgent(agent *mux.Router) {
	logger.Debug("API InitAgent - init URI")

	// registURI(agent, PUT, "/handshake", api.receiveHandshake)
	// registURI(agent, PUT, "/:agentKey", api.receivePolling)

	registURI(agent, PUT, "/handshake", api.receiveHandshake)
	registURI(agent, PUT, "/{agentKey}", api.receivePolling)
	registURI(agent, PUT, "/reports/{agentKey}", api.checkPrimaryInfo)

	// agent API 핸들러 추가
	agent.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: CustomHeader parsing 로직 추가

			// TODO: Support agent version 입력 추가

			// TODO: APIKey 인증 로직 추가

			// TODO: hash 검증 로직 추가

			next.ServeHTTP(w, r)

			// TODO: 전송구간 암호화 로직 추가

			// TODO: hash 생성 로직 추가
		})
	})
}

func (api *API) receiveHandshake(w http.ResponseWriter, r *http.Request) {
	ch := getCustomHeader(r)
	var cr = &common.Request{r}
	var pa Agent

	err := json.NewDecoder(r.Body).Decode(&pa)
	if err != nil {
		common.HTTPError(500, w, err, "JSON parsing error")
		return
	}

	logger.Debug(fmt.Sprintf("CustomHeader : %v", ch))

	var group = getAgentGroup(ch.ZoneID)

	logger.Debugf("%v", group)

	if group.ID == 0 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "Does not exist zone for zoneId : %d", ch.ZoneID)
		return
	}

	var agent = api.getAgent(ch.AgentKey)

	if agent.AgentKey == "" {
		agent.AgentKey = ch.AgentKey
		agent.GroupID = ch.ZoneID
		agent.IsActive = true
		agent.LastAccessTime = time.Now().UTC()
		agent.Ip = cr.Param("ip")

		api.addAgent(agent)
	}

	logger.Debug(agent)

	api.DB.Model(&group).Related(&agent)

	logger.Debug(group)
	logger.Debug(agent)
}

func (api *API) receivePolling(w http.ResponseWriter, r *http.Request) {

}

func (api *API) checkPrimaryInfo(w http.ResponseWriter, r *http.Request) {

}
