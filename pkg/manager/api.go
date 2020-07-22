package manager

import (
	"net/http"

	"github.com/NexClipper/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
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
	Root    *mux.Router
	APIRoot *mux.Router

	Legacy  *mux.Router
	Agent   *mux.Router
	Install *mux.Router
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
func Init(db *gorm.DB, baseRouter *mux.Router) *API {
	logger.Debug("API Init")

	api := &API{
		BaseRoutes: &Routes{},
		DB:         db,
	}

	api.BaseRoutes.Root = baseRouter
	api.BaseRoutes.APIRoot = baseRouter.PathPrefix(APIURLPrefix).Subrouter()

	api.BaseRoutes.Legacy = api.BaseRoutes.APIRoot
	api.BaseRoutes.Agent = api.BaseRoutes.APIRoot.PathPrefix("/agents").Subrouter()
	api.BaseRoutes.Install = api.BaseRoutes.APIRoot.PathPrefix("/install").Subrouter()

	api.InitLegacy(api.BaseRoutes.Legacy)
	api.InitAgent(api.BaseRoutes.Agent)
	api.InitInstall(api.BaseRoutes.Install)

	return api
}

func registURI(r *mux.Router, method int, uri string, f func(http.ResponseWriter, *http.Request)) {
	switch method {
	case ANY:
		r.HandleFunc(uri, f)
	case GET:
		r.HandleFunc(uri, f).Methods("GET")
	case POST:
		r.HandleFunc(uri, f).Methods("POST")
	case PUT:
		r.HandleFunc(uri, f).Methods("PUT")
	case DELETE:
		r.HandleFunc(uri, f).Methods("DELETE")
	case PATCH:
		r.HandleFunc(uri, f).Methods("PATCH")
	}
}
