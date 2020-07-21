package manager

import (
	"github.com/NexClipper/logger"
	"github.com/gin-gonic/gin"
)

// InitInstall initialize install API
func (api *API) InitInstall(install *gin.RouterGroup) {
	logger.Debug("API InitInstall - init URI")

	registURI(install, POST, "/agents/bootstrap", api.generateBootstrapCommand)
	registURI(install, GET, "/agents/download", api.ackprimary)
}

// agent setup script 생성
func (api *API) generateBootstrapCommand(c *gin.Context) {

}

// agent 다운로드
func (api *API) downloadAgent(c *gin.Context) {

}
