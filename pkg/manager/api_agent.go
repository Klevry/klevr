package manager

import (
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
}

func (api *API) receiveHandshake(w http.ResponseWriter, r *http.Request) {

}

func (api *API) receivePolling(w http.ResponseWriter, r *http.Request) {

}

func (api *API) checkPrimaryInfo(w http.ResponseWriter, r *http.Request) {

}
