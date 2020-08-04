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

// AgentStatusUpdateTask lock task name for agent status update
const AgentStatusUpdateTask = "AGENT_STATUS_UPDATE"

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
	Agent  AgentInfo
	DB     klevr.DBInfo
}

// ServerInfo klevr manager server info struct
type ServerInfo struct {
	Port              int
	ReadTimeout       int
	WriteTimeout      int
	EncryptionKey     string
	TransEncKey       string
	StatusUpdateCycle int
}

//AgentInfo klevr agent info struct
type AgentInfo struct {
	LogLevel  string // 로그 레벨
	CallCycle int    // 호출 간격
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

	serverConfig := manager.Config.Server

	db, err := manager.Config.DB.Connect()
	if err != nil {
		logger.Debug("gggg")
		logger.Fatal("Database connect failed : ", err)
	}
	defer db.Close()

	s := &http.Server{
		Addr:         ":8090",
		Handler:      manager.RootRouter,
		ReadTimeout:  time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(serverConfig.WriteTimeout) * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}

	go manager.updateAgentStatus(db, time.Duration(manager.Config.Server.StatusUpdateCycle))

	Init(manager, db)

	return s.ListenAndServe()
}

func (manager *KlevrManager) updateAgentStatus(db *xorm.Engine, cycle time.Duration) {
	for {
		time.Sleep(cycle * time.Second)
		logger.Debugf("sleep duration : %+v", cycle*time.Second)

		conn := db.NewSession()
		duration := time.Duration(manager.Config.Server.StatusUpdateCycle) * time.Second

		common.Block{
			Try: func() {
				if checkLock(conn, manager.InstanceID, duration) {
					current := time.Now().UTC()
					before := current.Add(-time.Duration(manager.Config.Server.StatusUpdateCycle) * time.Second)

					cnt, agents := getAgentsForInactive(conn, before)

					if cnt > 0 {
						len := len(*agents)
						ids := make([]uint64, len)

						for i := 0; i < len; i++ {
							ids[i] = (*agents)[i].Id
						}

						updateAgentStatus(conn, ids)
					}

					conn.Commit()
				}
			},
			Catch: func(e common.Exception) {
				if !conn.IsClosed() {
					conn.Rollback()
				}

				logger.Errorf("update agent status failed : %+v", e)
			},
			Finally: func() {
				if !conn.IsClosed() {
					conn.Close()
				}
			},
		}.Do()
	}
}

func checkLock(conn *xorm.Session, instanceID string, d time.Duration) bool {
	var hasLock = false

	lock, exist := getLock(conn, AgentStatusUpdateTask)

	if !exist {
		lock.Task = AgentStatusUpdateTask
		lock.InstanceId = instanceID
		lock.LockDate = time.Now().UTC()

		insertLock(conn, lock)

		hasLock = true
	} else if expired(lock.LockDate, d) || lock.InstanceId == instanceID {
		lock.InstanceId = instanceID
		lock.LockDate = time.Now().UTC()

		updateLock(conn, lock)

		hasLock = true
	}

	return hasLock
}

func expired(lockDate time.Time, d time.Duration) bool {
	current := time.Now().UTC()
	compare := lockDate.Add(d)

	if current.After(compare) {
		return true
	}

	return false
}
