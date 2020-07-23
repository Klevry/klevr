package manager

import (
	"fmt"
	"net/http"

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
	// zoneId := ch.ZoneID
	// agentKey := ch.AgentKey

	logger.Debug(fmt.Sprintf("CustomHeader : %v", ch))

	api.DB.LogMode(IsDebug)

}

func (api *API) receivePolling(w http.ResponseWriter, r *http.Request) {

}

func (api *API) checkPrimaryInfo(w http.ResponseWriter, r *http.Request) {

}
