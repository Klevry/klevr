package api

import (
	"github.com/NexClipper/logger"
	"github.com/gin-gonic/gin"
)

// APIHandler API handler
func (api *API) APIHandler(h func(*gin.Context)) gin.HandlerFunc {
	logger.Debug("pre handle")

	return h
}

// Test test handler
func Test(c *gin.Context) {
	logger.Debug("호출전")

	c.Next()

	logger.Debug("호출후")
}
