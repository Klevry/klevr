package manager

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/event"
	"github.com/Klevry/klevr/pkg/model"
	"github.com/Klevry/klevr/pkg/queue"
	"github.com/Klevry/klevr/pkg/rabbitmq"
	"github.com/Klevry/klevr/pkg/serialize"
	"github.com/NexClipper/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	concurrent "github.com/orcaman/concurrent-map"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
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
	HasLock    bool
	//EventQueue        queue.Queue
	HandOverTaskQueue queue.Queue
	ShutdownTasks     concurrent.ConcurrentMap
	//Mq                *ManagerMQ
	Event event.EventManager
}

type ManagerMQ struct {
	Connection *rabbitmq.Connection
	Queue      *amqp.Queue
}

// Config klevr manager config struct
type Config struct {
	Server  ServerInfo
	Agent   AgentInfo
	DB      common.DBInfo
	Cache   CacheInfo
	Console ConsoleInfo
}

// ServerInfo klevr manager server info struct
type ServerInfo struct {
	Port              int
	ReadTimeout       int
	WriteTimeout      int
	EncryptionKey     string
	TransEncKey       string
	StatusUpdateCycle int
	EventHandler      string
	Webhook           Webhook
	Mq                Mq
}

type Webhook struct {
	Url       string
	HookTerm  int
	HookCount int
}

type Mq struct {
	Url        []string
	Name       string
	Durable    bool
	AutoDelete bool
}

//AgentInfo klevr agent info struct
type AgentInfo struct {
	LogLevel  string // 로그 레벨
	CallCycle int    // 호출 간격
}

type ConsoleInfo struct {
	Usage bool
}

type CacheInfo struct {
	Type     string
	Address  string
	Port     int
	Password string
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
		RootRouter:        router,
		HandOverTaskQueue: queue.NewMutexQueue(),
		ShutdownTasks:     concurrent.New(),
		HasLock:           false,
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
		logger.Fatal("Database connect failed : ", err)
	}
	defer db.Close()

	ctx = common.BaseContext
	ctx.Put(CtxServer, manager)
	ctx.Put(CtxDbConn, db)
	ctx.Put(CtxPrimary, &sync.Mutex{})

	if manager.Config.DB.Cache == true {
		ctx.Put(CtxCacheLock, &sync.Mutex{})
	}

	cache := NewAgentStorage(manager.Config.Cache.Address, manager.Config.Cache.Port, manager.Config.Cache.Password)
	if cache == nil {
		logger.Fatalf("Cache connect failed: address(%s:%d)", manager.Config.Cache.Address, manager.Config.Cache.Port)
	}

	defer cache.Close()

	ctx.Put(CtxCacheConn, cache)

	if strings.ToLower(serverConfig.EventHandler) == "mq" {
		mqConfig := serverConfig.Mq
		eventOpt := event.KlevrEventOption{URL: mqConfig.Url, MQ_Name: mqConfig.Name, MQ_Durable: mqConfig.Durable, MQ_AutoDelete: mqConfig.AutoDelete}
		manager.Event = event.NewEventMQ(eventOpt)
	} else {
		webConfig := serverConfig.Webhook
		eventOpt := event.KlevrEventOption{URL: []string{webConfig.Url}, Web_HookCount: webConfig.HookCount, Web_HookTerm: webConfig.HookTerm}
		manager.Event = event.NewEventWeb(eventOpt)
	}

	defer manager.Event.Close()

	headerOk := handlers.AllowedHeaders([]string{"*"})
	//originOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	originOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "OPTIONS", "DELETE"})
	credentialsOk := handlers.AllowCredentials()
	corsHandler := handlers.CORS(headerOk, originOk, methodsOk, credentialsOk)(manager.RootRouter)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverConfig.Port),
		Handler:      corsHandler,
		ReadTimeout:  time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(serverConfig.WriteTimeout) * time.Second,
	}

	go manager.getLock(common.FromContext(ctx))
	// klevr 에이전트 상태 체크 및 업데이트 고루틴 시작(서버 lock 획득 시에만 동작)
	go manager.updateAgentStatus(common.FromContext(ctx), manager.Config.Server.StatusUpdateCycle)
	// klevr task 스케쥴 체크 및 업데이트 고루틴 시작(서버 lock 획득 시에만 동작)
	go manager.updateScheduledTask(common.FromContext(ctx), manager.Config.Agent.CallCycle)
	// klevr task hand-over 상태 업데이트
	go manager.startTaskHandoverUpdater(common.FromContext(ctx))

	Init(ctx)

	return s.ListenAndServe()
}

func (manager *KlevrManager) startTaskHandoverUpdater(ctx *common.Context) {
	db := CtxGetDbConn(ctx)
	q := manager.HandOverTaskQueue

	q.AddListener(1, func(q queue.Queue, args ...interface{}) {
		var ids []uint64
		var iq = q

		logger.Debugf("hand-over task queue count : %d", iq.Length())

		tx := &Tx{db.NewSession()}

		common.Block{
			Try: func() {
				err := tx.Begin()
				if err != nil {
					logger.Errorf("DB session begin error : %v", err)
					common.Throw(err)
				}

				for iq.Length() > 0 {
					t := iq.Pop().(model.Tasks)

					if t.Status != model.HandOver {
						ids = append(ids, t.Id)
					}
				}

				if len(ids) > 0 {
					tx.updateHandoverTasks(ids)
				}

				tx.Commit()
			},
			Catch: func(e error) {
				if !tx.IsClosed() {
					tx.Rollback()
				}

				logger.Errorf("update task status to hand-over failed : %+v", e)
			},
			Finally: func() {
				if !tx.IsClosed() {
					tx.Close()
				}
			},
		}.Do()
	})
}

func (manager *KlevrManager) getLock(ctx *common.Context) {
	for {
		db := CtxGetDbConn(ctx)

		st := time.Duration(manager.Config.Server.StatusUpdateCycle) / 2
		if st == 0 {
			st = 1
		}

		time.Sleep(st * time.Second)
		logger.Debugf("getLock sleep duration : %+v", st*time.Second)

		tx := &Tx{db.NewSession()}
		duration := time.Duration(manager.Config.Server.StatusUpdateCycle) * time.Second

		common.Block{
			Try: func() {
				err := tx.Begin()
				if err != nil {
					logger.Errorf("DB session begin error : %v", err)
					common.Throw(err)
				}

				manager.HasLock = checkLock(tx, manager.InstanceID, duration)
				tx.Commit()
			},
			Catch: func(e error) {
				if !tx.IsClosed() {
					tx.Rollback()
				}

				logger.Errorf("getLock failed : %+v", e)
			},
			Finally: func() {
				if !tx.IsClosed() {
					tx.Close()
				}
			},
		}.Do()
	}
}

func (manager *KlevrManager) updateScheduledTask(ctx *common.Context, cycle int) {
	st := cycle / 2
	if st < 1 {
		st = 1
	}

	sleep := time.Duration(st) * time.Second

	for {
		if manager.HasLock {
			db := CtxGetDbConn(ctx)

			time.Sleep(sleep)
			logger.Debugf("sleep duration : %+v", sleep)

			tx := &Tx{db.NewSession()}

			common.Block{
				Try: func() {
					err := tx.Begin()
					if err != nil {
						logger.Errorf("DB session begin error : %v", err)
						common.Throw(err)
					}

					cnt := tx.updateScheduledTask()

					logger.Debugf("Scheduled tasks status updated to wait-polling : %d", cnt)

					tx.Commit()
				},
				Catch: func(e error) {
					if !tx.IsClosed() {
						tx.Rollback()
					}

					logger.Errorf("update scheduled task failed : %+v", e)
				},
				Finally: func() {
					if !tx.IsClosed() {
						tx.Close()
					}
				},
			}.Do()
		}
	}
}

/*func (manager *KlevrManager) startEventHandler() {
	webhookConf := manager.Config.Server.Webhook
	url := webhookConf.Url

	q := manager.EventQueue

	if url != "" {
		var nilTime time.Time = time.Time{}
		var cntExecutedTime time.Time

		if webhookConf.HookCount > 0 {
			q.AddListener(uint32(webhookConf.HookCount), func(q queue.Queue, args ...interface{}) {
				var items []KlevrEvent
				var iq = q

				logger.Debugf("event queue count : %d", iq.Length())

				for iq.Length() > 0 {
					items = append(items, *(iq.Pop().(*KlevrEvent)))
				}

				logger.Debugf("%+v", items)

				sendBulkEventWebHook(url, &items)

				cntExecutedTime = time.Now().UTC()
			})
		}

		if webhookConf.HookTerm > 0 {
			baseTime := time.Duration(webhookConf.HookTerm) * time.Second
			sleepTime := baseTime

			for {
				logger.Debugf("Webhook sleep time : %+v", sleepTime)
				time.Sleep(sleepTime)

				if cntExecutedTime != nilTime {
					sleepTime = baseTime - (time.Duration(int(time.Now().UTC().Sub(cntExecutedTime))) * time.Second)
					logger.Debugf("Webhook new sleep time : %+v", sleepTime)
					cntExecutedTime = nilTime
					continue
				}

				var items []KlevrEvent

				logger.Debugf("Webhook scheduler event count : %d", q.Length())

				for q.Length() > 0 {
					items = append(items, *(q.Pop().(*KlevrEvent)))
				}

				logger.Debugf("%+v", items)

				if len(items) > 0 {
					sendBulkEventWebHook(url, &items)
				}

				sleepTime = baseTime
			}
		}
	}
}*/

func AddHandOverTasks(tasks *[]model.Tasks) {
	manager := common.BaseContext.Get(CtxServer).(*KlevrManager)

	q := manager.HandOverTaskQueue
	q.BulkPush(*tasks)
}

func AddShutdownTask(task *model.Tasks) bool {
	manager := common.BaseContext.Get(CtxServer).(*KlevrManager)

	if !manager.ShutdownTasks.Has(task.AgentKey) {
		manager.ShutdownTasks.Set(task.AgentKey, task)
		return true
	}

	return false
}

func CheckShutdownTask(agentKey string) (uint64, bool) {
	manager := common.BaseContext.Get(CtxServer).(*KlevrManager)

	v, ok := manager.ShutdownTasks.Get(agentKey)
	if ok {
		t := v.(*model.Tasks)
		return t.Id, ok
	}
	return 0, ok
}

func RemoveShutdownTask(agentKeys []string) {
	manager := common.BaseContext.Get(CtxServer).(*KlevrManager)

	for _, k := range agentKeys {
		manager.ShutdownTasks.Remove(k)
	}
}

func (manager *KlevrManager) updateAgentStatus(ctx *common.Context, cycle int) {
	st := cycle / 2
	if st < 1 {
		st = 1
	}

	sleep := time.Duration(st) * time.Second

	for {
		if manager.HasLock {
			db := CtxGetDbConn(ctx)

			time.Sleep(sleep)
			logger.Debugf("sleep duration : %+v", sleep)

			tx := &Tx{db.NewSession()}

			common.Block{
				Try: func() {
					err := tx.Begin()
					if err != nil {
						logger.Errorf("DB session begin error : %v", err)
						common.Throw(err)
					}

					current := time.Now().UTC()
					before := current.Add(-time.Duration(manager.Config.Server.StatusUpdateCycle) * time.Second)

					txManager := CtxGetCacheConn(ctx)
					cnt, agents := txManager.GetAgentsForInactive(ctx, tx, before)

					if cnt > 0 {
						len := len(*agents)
						inactiveIDs := make([]uint64, len)
						inactiveAgentKeys := make([]string, len)
						forceShutdownAgentKeys := make([]string, 0)
						taskIDs := make([]uint64, 0)

						var events = make([]event.KlevrEvent, len)
						var eventTime = &serialize.JSONTime{Time: time.Now().UTC()}

						for i := 0; i < len; i++ {
							agent := (*agents)[i]

							inactiveIDs[i] = agent.Id
							inactiveAgentKeys[i] = agent.AgentKey

							if tid, ok := CheckShutdownTask(agent.AgentKey); ok {
								forceShutdownAgentKeys = append(forceShutdownAgentKeys, agent.AgentKey)
								taskIDs = append(taskIDs, tid)
							}

							events[i] = event.KlevrEvent{
								EventType: event.AgentDisconnect,
								AgentKey:  agent.AgentKey,
								GroupID:   agent.GroupId,
								Result:    "",
								EventTime: eventTime,
							}

							logger.Debugf("disconnected event : [%+v]", events[i])
						}

						//tx.updateAgentStatus(ids)
						txManager.UpdateAgentStatus(ctx, tx, agents, inactiveIDs)
						tx.updateInitIterationTasks(inactiveAgentKeys)
						tx.updateShutdownTasks(taskIDs)

						RemoveShutdownTask(forceShutdownAgentKeys)

						//AddEvents(&events)
						manager.Event.AddEvents(&events)
					}

					tx.Commit()
				},
				Catch: func(e error) {
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

func (manager *KlevrManager) encrypt(msg string) string {
	if msg == "" {
		return ""
	}

	encKey := manager.Config.Server.EncryptionKey

	enc, err := common.Encrypt(encKey, msg)
	if err != nil {
		logger.Errorf("An error occurred during encryption.\n%+v", errors.WithStack(err))
		panic("Internal server error")
	}

	return enc
}

func (manager *KlevrManager) decrypt(encrypted string) string {
	if encrypted == "" {
		return ""
	}

	encKey := manager.Config.Server.EncryptionKey

	dec, err := common.Decrypt(encKey, encrypted)
	if err != nil {
		logger.Errorf("An error occurred during decryption.\n%+v", errors.WithStack(err))
		panic("Internal server error")
	}

	return dec
}
