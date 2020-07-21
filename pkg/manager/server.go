package manager

import (
	"net/http"
	"time"

	klevr "github.com/Klevry/klevr/pkg/common"
	"github.com/NexClipper/logger"
	"github.com/gin-gonic/gin"
)

// KlevrManager klevr manager struct
type KlevrManager struct {
	ServerName string
	Config     Config
	RootRouter *gin.Engine
}

// Config klevr manager config struct
type Config struct {
	Server ServerInfo
	Db     klevr.DBInfo
}

// ServerInfo klevr manager server info struct
type ServerInfo struct {
	Port int
}

// SetConfig setter for Config struct
func (manager *KlevrManager) SetConfig(config *Config) {
	manager.Config = *config
}

// NewKlevrManager constructor for KlevrManager
func NewKlevrManager() (*KlevrManager, error) {
	router := gin.Default()

	instance := &KlevrManager{
		RootRouter: router,
	}

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

	Init(db, manager.RootRouter)

	return s.ListenAndServe()
}
