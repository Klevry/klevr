package manager

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gorilla/mux"
)

//  InitInstall initialize install API

func (api *API) InitInstall(install *mux.Router) {
	logger.Debug("API InitInstall - init URI")

	// registURIWithQuery(install, POST, "/agents/bootstrap", api.generateBootstrapCommand, "platform", "{platform}", "zoneId", "{zoneId:[0-9]+}")
	registURI(install, POST, "/agents/bootstrap", api.generateBootstrapCommand)
	registURI(install, GET, "/agents/download", api.downloadAgent)
}

type BootstrapCommand struct {
	APIKey     string `json:"apiKey"`
	Platform   string `json:"platform"`
	ManagerURL string `json:"managerUrl"`
	ZoneID     uint64 `json:"zoneId"`
}

//
// generateBootstrapCommand godoc
// @Summary agent setup script 생성
// @Description 에이전트가 설치되기 위한 최소 정보를 이용해 설치를 할 수 있는 스크립트를 생성한다
// @Tags Install
// @Accept json
// @Produce plain
// @Router /install/agents/bootstrap [post]
// @Param b body manager.BootstrapCommand true "Bootstrap Info"
// @Success 200 {string} string
func (api *API) generateBootstrapCommand(w http.ResponseWriter, r *http.Request) {
	var cmd = "curl -sL gg.gg/klevr |bash  && ./klevr -apiKey=\"{apiKey}\" -platform={platform} -manager=\"{managerUrl}\" -zoneId={zoneId}"

	var bc BootstrapCommand
	err := json.NewDecoder(r.Body).Decode(&bc)
	if err != nil {
		common.WriteHTTPError(500, w, err, "JSON parsing error")
		return
	}

	zone := strconv.FormatUint(bc.ZoneID, 10)

	cmd = strings.Replace(cmd, "{apiKey}", bc.APIKey, 1)
	cmd = strings.Replace(cmd, "{platform}", bc.Platform, 1)
	cmd = strings.Replace(cmd, "{managerUrl}", bc.ManagerURL, 1)
	cmd = strings.Replace(cmd, "{zoneId}", zone, 1)

	//w.WriteHeader(http.StatusUnauthorized)
	// w.Write([]byte(fmt.Sprintf("%s", cmd)))
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "%s", cmd)
}

// agent 다운로드
func (api *API) downloadAgent(w http.ResponseWriter, r *http.Request) {

}
