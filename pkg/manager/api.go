package manager

import (
	"github.com/NexClipper/logger"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// APIURLPrefix default API URL prefix
const APIURLPrefix string = "/"

// HTTP method constants
const (
	ANY int = iota
	GET
	POST
	PUT
	DELETE
	PATCH
)

// Routes router struct
type Routes struct {
	Root    *gin.Engine
	APIRoot *gin.RouterGroup

	Legacy  *gin.RouterGroup
	Agent   *gin.RouterGroup
	Install *gin.RouterGroup
}

// API api struct
type API struct {
	BaseRoutes *Routes
	DB         *gorm.DB
}

type apiDef struct {
	method   string
	uri      string
	function func(*gin.Context)
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
	api.BaseRoutes.Agent = api.BaseRoutes.APIRoot.Group("/agents")
	api.BaseRoutes.Install = api.BaseRoutes.APIRoot.Group("/install")

	api.InitLegacy(api.BaseRoutes.Legacy)
	api.InitAgent(api.BaseRoutes.Agent)
	api.InitInstall(api.BaseRoutes.Install)

	return api
}

func registURI(g *gin.RouterGroup, method int, uri string, f func(c *gin.Context)) {
	switch method {
	case ANY:
		g.Any(uri, f)
	case GET:
		g.GET(uri, f)
	case POST:
		g.POST(uri, f)
	case PUT:
		g.PUT(uri, f)
	case DELETE:
		g.DELETE(uri, f)
	case PATCH:
		g.PATCH(uri, f)
	}
}
