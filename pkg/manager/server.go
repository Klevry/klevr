package manager

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	klevr "github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gorilla/mux"
	"xorm.io/xorm"
)

// IsDebug debugabble for all
var IsDebug = false

// KlevrManager klevr manager struct
type KlevrManager struct {
	ServerName string
	Config     Config
	RootRouter *mux.Router
	InstanceID string
}

// Config klevr manager config struct
type Config struct {
	Server ServerInfo
	Db     klevr.DBInfo
}

// ServerInfo klevr manager server info struct
type ServerInfo struct {
	Port          int
	EncryptionKey string
}

func init() {
	var level = logger.GetLevel()

	if level == 0 {
		IsDebug = true
	}
}

// SetConfig setter for Config struct
func (manager *KlevrManager) SetConfig(config *Config) {
	manager.Config = *config
}

// NewKlevrManager constructor for KlevrManager
func NewKlevrManager() (*KlevrManager, error) {
	router := mux.NewRouter()

	instance := &KlevrManager{
		RootRouter: router,
	}

	instance.InstanceID = fmt.Sprintf("%v_%v", &instance, time.Now().UTC().Unix())

	return instance, nil
}

// Run run klevr manager
func (manager *KlevrManager) Run() error {
	logger.Info(manager)

	db, err := manager.Config.Db.Connect()
	if err != nil {
		logger.Debug("gggg")
		logger.Fatal("Database connect failed : ", err)
	}
	defer db.Close()

	s := &http.Server{
		Addr:         ":8090",
		Handler:      manager.RootRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}

	go manager.updateAgentStatus(db, 10)

	Init(manager, db, manager.RootRouter)

	return s.ListenAndServe()
}

func (manager *KlevrManager) updateAgentStatus(db *xorm.Engine, cycle time.Duration) {
	for {
		time.Sleep(cycle * time.Second)

		conn := db.NewSession()

		defer func() {
			defer func() {
				if !conn.IsClosed() {
					conn.Close()
				}
			}()

			r := recover()
			if r != nil {
				logger.Errorf("recovered : %v", r)
			}

			if !conn.IsClosed() {
				conn.Rollback()
			}
		}()

		common.ErrorWithPanic(conn.Begin(),
			"updateAgentStatus() connection open error")

		getLock(conn, "UPDATE_AGENT_STATUS")

		common.ErrorWithPanic(conn.Commit(),
			"updateAgentStatus() commit error")
	}
}
