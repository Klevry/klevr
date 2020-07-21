package manager

import (
	"github.com/NexClipper/logger"
	"github.com/gin-gonic/gin"
)

// InitAgent initialize agent API
func (api *API) InitAgent(agent *gin.RouterGroup) {
	logger.Debug("API InitAgent - init URI")

	// registURI(agent, PUT, "/handshake", api.receiveHandshake)
	// registURI(agent, PUT, "/:agentKey", api.receivePolling)
	registURI(agent, PUT, "/:path", api.route1Depth)
	registURI(agent, PUT, "/reports/:agentKey", api.checkPrimaryInfo)
}

func (api *API) route1Depth(c *gin.Context) {
	// path가 handshake인지, :agentKey인지 확인하여 분기처리
	path := c.Param("path")

	if path == "handshake" {
		api.receiveHandshake(c)
	} else {
		api.receivePolling(c)
	}
}

func (api *API) receiveHandshake(c *gin.Context) {

}

func (api *API) receivePolling(c *gin.Context) {

}

func (api *API) checkPrimaryInfo(c *gin.Context) {

}
