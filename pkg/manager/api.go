package manager

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/Klevry/klevr/pkg/common"
	_ "github.com/Klevry/klevr/pkg/manager/docs"
	concurrent "github.com/orcaman/concurrent-map"
	swagger "github.com/swaggo/http-swagger"

	"github.com/NexClipper/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
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
	Inner   *mux.Router
	Page    *mux.Router
}

// API api struct
type API struct {
	BaseRoutes  *Routes
	DB          *common.DB
	Manager     *KlevrManager
	APIKeyMap   concurrent.ConcurrentMap
	BlockKeyMap concurrent.ConcurrentMap
}

type apiDef struct {
	method   string
	uri      string
	function func(*gin.Context)
}

// Init initialize API router
// @title Klevr-Manager API
// @version 1.0
// @description
// @contact.name mrchopa
// @contact.email ys3gods@gmail.com
// @BasePath /
func Init(ctx *common.Context) *API {
	logger.Debug("API Init")

	api := &API{
		BaseRoutes:  &Routes{},
		DB:          CtxGetDbConn(ctx),
		Manager:     CtxGetServer(ctx),
		APIKeyMap:   concurrent.New(),
		BlockKeyMap: concurrent.New(),
	}

	ctx.Put(CtxAPI, api)

	api.DB.ShowSQL(true)
	// TODO: ContextLogger interface 구현하여 logger override
	// api.DB.SetLogger(log.NewSimpleLogger(f))

	api.BaseRoutes.Root = api.Manager.RootRouter

	// swagger 설정
	api.BaseRoutes.Root.PathPrefix("/swagger").Handler(swagger.WrapHandler)

	api.BaseRoutes.APIRoot = api.BaseRoutes.Root

	api.BaseRoutes.Legacy = api.BaseRoutes.APIRoot

	api.BaseRoutes.Agent = api.BaseRoutes.APIRoot.PathPrefix("/agents").Subrouter()
	api.BaseRoutes.Agent.Use(CommonWrappingHandler(ctx))
	api.BaseRoutes.Agent.Use(RequestInfoLoggerHandler)
	api.BaseRoutes.Install = api.BaseRoutes.APIRoot.PathPrefix("/install").Subrouter()
	api.BaseRoutes.Install.Use(CommonWrappingHandler(ctx))
	api.BaseRoutes.Install.Use(RequestInfoLoggerHandler)
	api.BaseRoutes.Inner = api.BaseRoutes.APIRoot.PathPrefix("/inner").Subrouter()
	api.BaseRoutes.Inner.Use(CommonWrappingHandler(ctx))
	api.BaseRoutes.Inner.Use(RequestInfoLoggerHandler)
	api.BaseRoutes.Page = api.BaseRoutes.APIRoot.PathPrefix("/page").Subrouter()
	api.BaseRoutes.Page.Use(CommonWrappingHandler(ctx))
	api.BaseRoutes.Page.Use(RequestInfoLoggerHandler)

	// api.InitLegacy(api.BaseRoutes.Legacy)
	api.InitAgent(api.BaseRoutes.Agent)
	api.InitInstall(api.BaseRoutes.Install)
	api.InitInner(api.BaseRoutes.Inner)
	if api.Manager.Config.Page.Secret != "" {
		api.InitPage(api.BaseRoutes.Page)
	}

	// health check handler(~/health)
	api.BaseRoutes.Root.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	err := api.BaseRoutes.Root.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {

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
func GetDBConn(ctx *common.Context) *Tx {
	tx := CtxGetDbSession(ctx)
	if tx == nil {
		logger.Warningf("The variable in context is not DB session : %d", debug.Stack())
		panic("DB session is not exist")
	}

	return tx
}
