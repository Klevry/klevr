package manager

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gorilla/mux"
)

// InitInstall initialize install API
func (api *API) InitInstall(install *mux.Router) {
	logger.Debug("API InitInstall - init URI")

	// registURIWithQuery(install, POST, "/agents/bootstrap", api.generateBootstrapCommand, "platform", "{platform}", "zoneId", "{zoneId:[0-9]+}")
	registURI(install, POST, "/agents/bootstrap", api.generateBootstrapCommand)
	registURI(install, GET, "/agents/download", api.downloadAgent)
}

// agent setup script 생성
func (api *API) generateBootstrapCommand(w http.ResponseWriter, r *http.Request) {
	var cmd = "curl -sL bit.ly/klevry |bash  && ./klevr -apiKey=\"{apiKey}\" -platform={platform} -manager=\"{managerUrl}\" -zoneId={zoneId}"

	var cr = &common.Request{r}

	cmd = strings.Replace(cmd, "{apiKey}", cr.Param("apiKey"), 1)
	cmd = strings.Replace(cmd, "{platform}", cr.Param("platform"), 1)
	cmd = strings.Replace(cmd, "{managerUrl}", cr.Param("managerUrl"), 1)
	cmd = strings.Replace(cmd, "{zoneId}", cr.Param("zoneId"), 1)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", cmd)
}

// agent 다운로드
func (api *API) downloadAgent(w http.ResponseWriter, r *http.Request) {

}
