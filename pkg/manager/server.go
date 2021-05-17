package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/Klevry/klevr/pkg/common"
	"github.com/Klevry/klevr/pkg/rabbitmq"
	"github.com/NexClipper/logger"
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
	ServerName        string
	Config            Config
	RootRouter        *mux.Router
	InstanceID        string
	HasLock           bool
	EventQueue        *common.Queue
	HandOverTaskQueue *common.Queue
	ShutdownTasks     concurrent.ConcurrentMap
	Mq                *ManagerMQ
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
		EventQueue:        common.NewMutexQueue(),
		HandOverTaskQueue: common.NewMutexQueue(),
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

	if serverConfig.EventHandler == "mq" {
		mqConfig := serverConfig.Mq

		mqConn, err := rabbitmq.DialCluster(mqConfig.Url)
		if err != nil {
			logger.Errorf("Failed to connect to MQ - %+v", errors.Cause(err))
			panic(err)
		}

		defer mqConn.Close()

		mqChannel, err := mqConn.Channel()
		if err != nil {
			logger.Errorf("Failed to open a channel to MQ - %+v", errors.Cause(err))
			panic(err)
		}

		queue, err := mqChannel.QueueDeclare(mqConfig.Name, mqConfig.Durable, mqConfig.AutoDelete, false, false, nil)
		if err != nil {
			logger.Errorf("Failed to declare queue from MQ - %+v", errors.Cause(err))
			panic(err)
		}

		mqChannel.Close()

		manager.Mq = &ManagerMQ{
			Connection: mqConn,
			Queue:      &queue,
		}
	}

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverConfig.Port),
		Handler:      manager.RootRouter,
		ReadTimeout:  time.Duration(serverConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(serverConfig.WriteTimeout) * time.Second,
	}

	go manager.getLock(common.FromContext(ctx))
	// klevr 이벤트 hook(발송) 고루틴 핸들러 시작(항상 동작 - 강제 종료시 이벤트 메모리 소실)
	if serverConfig.EventHandler != "mq" {
		go manager.startEventHandler()
	}
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
	q := *manager.HandOverTaskQueue

	q.AddListener(1, func(q *common.Queue, args ...interface{}) {
		var ids []uint64
		var iq = *q

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
					t := iq.Pop().(Tasks)

					if t.Status != common.HandOver {
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

func AddHandOverTasks(tasks *[]Tasks) {
	manager := common.BaseContext.Get(CtxServer).(*KlevrManager)

	q := *manager.HandOverTaskQueue
	q.BulkPush(*tasks)
}

func AddShutdownTask(task *Tasks) bool {
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
		t := v.(*Tasks)
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

// AddEvent add klevr event for webhook
func AddEvent(event *KlevrEvent) {
	logger.Debugf("add event : [%+v]", *event)

	manager := common.BaseContext.Get(CtxServer).(*KlevrManager)

	if manager.Mq != nil {
		arr := []KlevrEvent{*event}

		go sendBulkEventMQ(&arr)
	} else {
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
}

func AddEvents(events *[]KlevrEvent) {
	manager := common.BaseContext.Get(CtxServer).(*KlevrManager)

	if manager.Mq != nil {
		go sendBulkEventMQ(events)
	} else {
		hookConfig := manager.Config.Server.Webhook

		logger.Debugf("hookConfig : [%+v]", hookConfig)

		if hookConfig.Url == "" {
			return
		}

		go sendBulkEventWebHook(hookConfig.Url, events)
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
			logger.Errorf("sendEvent recover from - %v", r)
		}
	}()

	if events == nil {
		logger.Debug("Klevr events is nil")
		return
	}

	b, err := json.Marshal(*events)
	if err != nil {
		retryFailedEvent(events, false)
		panic("klevr webhook event marshal error.")
	}

	logger.Debugf("%+v", *events)
	logger.Debugf("%d", len(*events))
	logger.Debugf("%s", string(b))

	res, err := http.Post(url, "application/json", bytes.NewReader(b))

	if err != nil {
		logger.Warningf("Klevr event webhook send failed - %+v", err)
		retryFailedEvent(events, true)
	}

	if res == nil {
		return
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

func sendBulkEventMQ(events *[]KlevrEvent) {
	manager := common.BaseContext.Get(CtxServer).(*KlevrManager)
	mq := manager.Mq

	b, err := json.Marshal(*events)
	if err != nil {
		logger.Errorf("klevr event MQ publish marshal error - %+v", errors.Cause(err))
		retryFailedEvent(events, false)
	}

	channel, err := manager.Mq.Connection.Channel()
	if err != nil {
		logger.Errorf("Failed to open a channel to MQ - %+v", errors.Cause(err))
		retryFailedEvent(events, true)
	}

	defer channel.Close()

	err = channel.Publish("", mq.Queue.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: 2,
		Body:         b,
	})

	if err != nil {
		logger.Errorf("Failed to publish to MQ - %+v", errors.Cause(err))
		retryFailedEvent(events, true)
	}
}

// TODO: event 발송 실패 재처리 구현
func retryFailedEvent(events *[]KlevrEvent, retryable bool) {

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

					//cnt, agents := tx.getAgentsForInactive(before)
					txManager := NewAgentStorage()
					cnt, agents := txManager.GetAgentsForInactive(tx, before)

					if cnt > 0 {
						len := len(*agents)
						ids := make([]uint64, len)
						agentKeys := make([]string, 0)
						taskIDs := make([]uint64, 0)

						var events = make([]KlevrEvent, len)
						var eventTime = &common.JSONTime{Time: time.Now().UTC()}

						for i := 0; i < len; i++ {
							agent := (*agents)[i]

							ids[i] = agent.Id
							if tid, ok := CheckShutdownTask(agent.AgentKey); ok {
								agentKeys = append(agentKeys, agent.AgentKey)
								taskIDs = append(taskIDs, tid)
							}

							events[i] = KlevrEvent{
								EventType: AgentDisconnect,
								AgentKey:  agent.AgentKey,
								GroupID:   agent.GroupId,
								Result:    "",
								EventTime: eventTime,
							}

							logger.Debugf("disconnected event : [%+v]", events[i])
						}

						//tx.updateAgentStatus(ids)
						txManager := NewAgentStorage()
						txManager.UpdateAgentStatus(tx, agents, ids)
						tx.updateShutdownTasks(taskIDs)

						RemoveShutdownTask(agentKeys)

						AddEvents(&events)
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
