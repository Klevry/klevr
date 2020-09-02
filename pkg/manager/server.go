package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	EventQueue *common.Queue
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
	Webhook           Webhook
}

type Webhook struct {
	Url       string
	HookTerm  int
	HookCount int
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
		EventQueue: common.NewMutexQueue(),
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

	go manager.startEventHandler()
	go manager.updateAgentStatus(common.FromContext(ctx), time.Duration(manager.Config.Server.StatusUpdateCycle))

	Init(ctx)

	return s.ListenAndServe()
}
func (manager *KlevrManager) startEventHandler() {
	webhookConf := manager.Config.Server.Webhook
	url := webhookConf.Url

	q := *manager.EventQueue

	if url != "" {
		var nilTime time.Time = time.Time{}
		var cntExecutedTime time.Time

		if webhookConf.HookCount > 0 {
			q.AddListener(uint32(webhookConf.HookCount), func(q *common.Queue, args ...interface{}) {
				var items []KlevrEvent
				var iq = *q

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
}

// AddEvent add klevr event for webhook
func AddEvent(event *KlevrEvent) {
	logger.Debugf("add event : [%+v]", event)

	manager := common.BaseContext.Get(CtxServer).(*KlevrManager)
	hookConfig := manager.Config.Server.Webhook

	logger.Debugf("hookConfig : [%+v]", hookConfig)

	if hookConfig.Url == "" {
		return
	}

	if hookConfig.HookCount <= 1 && hookConfig.HookTerm < 1 {
		go sendSingleEventWebHook(hookConfig.Url, event)
	} else {
		q := *manager.EventQueue
		q.Push(event)
	}
}

func sendSingleEventWebHook(url string, event *KlevrEvent) {
	var arr = []KlevrEvent{*event}

	logger.Debugf("%+v", *event)
	logger.Debugf("%d", len(arr))

	sendBulkEventWebHook(url, &arr)
}

func sendBulkEventWebHook(url string, events *[]KlevrEvent) {
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("sendEvent recorver from - %v", r)
		}
	}()

	b, err := json.Marshal(*events)
	if err != nil {
		retryFailedEvent(events, false)
		panic("kelvr webhook event marshal error.")
	}

	logger.Debugf("%+v", *events)
	logger.Debugf("%d", len(*events))

	res, err := http.Post(url, "application/json", bytes.NewReader(b))

	if err != nil {
		logger.Warningf("Klevr event webhook send failed - %+v", err)
		retryFailedEvent(events, true)
	}

	defer func() {
		if res != nil {
			res.Body.Close()
		}
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Warningf("Klevr event webhook send failed - read response body failed - %+v", err)
		retryFailedEvent(events, true)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		logger.Warningf("Klevr event webhook send failed - status code : [%d], response body : [%s]", res.StatusCode, body)
		retryFailedEvent(events, true)
	}

	logger.Debugf("sendEventWebHook - statusCode : [%d], body : [%s]", res.StatusCode, body)
}

func retryFailedEvent(events *[]KlevrEvent, retryable bool) {

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
