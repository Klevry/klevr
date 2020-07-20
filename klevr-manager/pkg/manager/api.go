package manager

import (
	"github.com/NexClipper/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// APIURLPrefix default API URL prefix
const APIURLPrefix string = "/"

// Routes router struct
type Routes struct {
	Root    *gin.Engine
	APIRoot *gin.RouterGroup

	Legacy *gin.RouterGroup
}

// API api struct
type API struct {
	BaseRoutes *Routes
	DB         *gorm.DB
}

// Init initialize API router
func Init(db *gorm.DB, baseRouter *gin.Engine) *API {
	logger.Debug("API Init")

	baseRouter.Use(Test)

	api := &API{
		BaseRoutes: &Routes{},
		DB:         db,
	}

	api.BaseRoutes.Root = baseRouter
	api.BaseRoutes.APIRoot = baseRouter.Group(APIURLPrefix)

	api.BaseRoutes.Legacy = api.BaseRoutes.APIRoot

	api.InitLegacy(api.BaseRoutes.Legacy)

	return api
}
