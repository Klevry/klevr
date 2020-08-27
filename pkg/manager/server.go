package manager

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gorilla/mux"
)

// IsDebug debugabble for all
var IsDebug = false
var ctx *common.Context

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
	DB     common.DBInfo
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

	ctx = common.BaseContext
	ctx.Put(CtxServer, manager)
	ctx.Put(CtxDbConn, db)

	s := &http.Server{
		Addr:         ":8090",
		Handler:      manager.RootRouter,
		ReadTimeout:  time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(serverConfig.WriteTimeout) * time.Second,
	}

	go manager.updateAgentStatus(common.FromContext(ctx), time.Duration(manager.Config.Server.StatusUpdateCycle))

	Init(ctx)

	return s.ListenAndServe()
}

func (manager *KlevrManager) updateAgentStatus(ctx *common.Context, cycle time.Duration) {
	for {
		db := CtxGetDbConn(ctx)

		time.Sleep(cycle * time.Second)
		logger.Debugf("sleep duration : %+v", cycle*time.Second)

		tx := &Tx{db.NewSession()}
		duration := time.Duration(manager.Config.Server.StatusUpdateCycle) * time.Second

		common.Block{
			Try: func() {
				err := tx.Begin()
				if err != nil {
					logger.Errorf("DB session begin error : %v", err)
					common.Throw(err)
				}

				if checkLock(tx, manager.InstanceID, duration) {
					current := time.Now().UTC()
					before := current.Add(-time.Duration(manager.Config.Server.StatusUpdateCycle) * time.Second)

					cnt, agents := tx.getAgentsForInactive(before)

					if cnt > 0 {
						len := len(*agents)
						ids := make([]uint64, len)

						for i := 0; i < len; i++ {
							ids[i] = (*agents)[i].Id
						}

						tx.updateAgentStatus(ids)
					}

					tx.Commit()
				}
			},
			Catch: func(e common.Exception) {
				if !tx.IsClosed() {
					tx.Rollback()
				}

				logger.Errorf("update agent status failed : %+v", e)
			},
			Finally: func() {
				if !tx.IsClosed() {
					tx.Close()
				}
			},
		}.Do()
	}
}

func checkLock(tx *Tx, instanceID string, d time.Duration) bool {
	var hasLock = false

	lock, exist := tx.getLock(AgentStatusUpdateTask)

	if !exist {
		lock.Task = AgentStatusUpdateTask
		lock.InstanceId = instanceID
		lock.LockDate = time.Now().UTC()

		tx.insertLock(lock)

		hasLock = true
	} else if expired(lock.LockDate, d) || lock.InstanceId == instanceID {
		lock.InstanceId = instanceID
		lock.LockDate = time.Now().UTC()

		tx.updateLock(lock)

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
