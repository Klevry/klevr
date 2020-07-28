package manager

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gorilla/context"

	"github.com/NexClipper/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"xorm.io/xorm"
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

// DBConnContextName DB connection name for Request context
const DBConnContextName string = "CTX_DB_CONN"

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
	DB         *xorm.Engine
	InstanceID string
}

type apiDef struct {
	method   string
	uri      string
	function func(*gin.Context)
}

// Init initialize API router
func Init(instance *KlevrManager, db *xorm.Engine, baseRouter *mux.Router) *API {
	logger.Debug("API Init")

	api := &API{
		BaseRoutes: &Routes{},
		DB:         db,
		InstanceID: instance.InstanceID,
	}

	api.DB.ShowSQL(true)
	// TODO: ContextLogger interface 구현하여 logger override
	// api.DB.SetLogger(log.NewSimpleLogger(f))

	baseRouter.Use(api.CommonWrappingHandler(api.DB))
	baseRouter.Use(ExecutionInfoLoggerHandler)
	// baseRouter.Use(TestHandler)

	api.BaseRoutes.Root = baseRouter
	api.BaseRoutes.APIRoot = baseRouter.PathPrefix(APIURLPrefix).Subrouter()

	api.BaseRoutes.Legacy = api.BaseRoutes.APIRoot
	api.BaseRoutes.Agent = api.BaseRoutes.APIRoot.PathPrefix("/agents").Subrouter()
	api.BaseRoutes.Install = api.BaseRoutes.APIRoot.PathPrefix("/install").Subrouter()

	api.InitLegacy(api.BaseRoutes.Legacy)
	api.InitAgent(api.BaseRoutes.Agent)
	api.InitInstall(api.BaseRoutes.Install)

	err := baseRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {

		pathTemplate, _ := route.GetPathTemplate()
		// if err == nil {
		// 	logger.Info("ROUTE:", pathTemplate)
		// }
		// pathRegexp, err := route.GetPathRegexp()
		// if err == nil {
		// 	logger.Info("Path regexp:", pathRegexp)
		// }
		queriesTemplates, _ := route.GetQueriesTemplates()
		route.GetQueriesRegexp()
		// if err == nil {
		// 	logger.Info("Queries templates:", strings.Join(queriesTemplates, ","))
		// }
		// queriesRegexps, err := route.GetQueriesRegexp()
		// if err == nil {
		// 	logger.Info("Queries regexps:", strings.Join(queriesRegexps, ","))
		// }
		methods, _ := route.GetMethods()
		// if err == nil {
		// 	logger.Info("Methods:", strings.Join(methods, ","))
		// }

		logger.Info(fmt.Sprintf("[%s] %s [%s]", strings.Join(methods, ","), pathTemplate, strings.Join(queriesTemplates, ",")))

		return nil
	})
	if err != nil {
		logger.Error(err)
	}

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

func registURIWithQuery(r *mux.Router, method int, uri string, f func(http.ResponseWriter, *http.Request), q ...string) {
	switch method {
	case ANY:
		r.Path(uri).Queries(q...).HandlerFunc(f)
	case GET:
		r.Path(uri).Queries(q...).HandlerFunc(f).Methods("GET")
	case POST:
		r.Path(uri).Queries(q...).HandlerFunc(f).Methods("POST")
	case PUT:
		r.Path(uri).Queries(q...).HandlerFunc(f).Methods("PUT")
	case DELETE:
		r.Path(uri).Queries(q...).HandlerFunc(f).Methods("DELETE")
	case PATCH:
		r.Path(uri).Queries(q...).HandlerFunc(f).Methods("PATCH")
	}
}

// GetDBConn return DB connection(session) from Request context
func GetDBConn(r *http.Request) *xorm.Session {
	v := context.Get(r, DBConnContextName)
	if v == nil {
		logger.Warningf("The variable in context is not DB session : %d", debug.Stack())
		panic("DB is not exist")
	}

	db, ok := v.(*xorm.Session)
	if !ok {
		logger.Warningf("DB session not exist : %d", debug.Stack())
		panic("DB is not exist")
	}

	return db
}
