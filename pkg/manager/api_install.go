package manager

import (
	"net/http"

	"github.com/NexClipper/logger"
	"github.com/gorilla/mux"
)

// InitInstall initialize install API
func (api *API) InitInstall(install *mux.Router) {
	logger.Debug("API InitInstall - init URI")

	registURI(install, POST, "/agents/bootstrap", api.generateBootstrapCommand)
	registURI(install, GET, "/agents/download", api.downloadAgent)
}

// agent setup script 생성
func (api *API) generateBootstrapCommand(w http.ResponseWriter, r *http.Request) {

}

// agent 다운로드
func (api *API) downloadAgent(w http.ResponseWriter, r *http.Request) {

}
