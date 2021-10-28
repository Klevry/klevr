package event

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Klevry/klevr/pkg/queue"
	"github.com/NexClipper/logger"
)

type EventWeb struct {
	eventQueue queue.Queue
	url        []string
	hookCount  int
	hookTerm   int
}

func NewEventWeb(opt KlevrEventOption) EventManager {
	event := &EventWeb{
		eventQueue: queue.NewMutexQueue(),
		url:        opt.URL,
		hookCount:  opt.Web_HookCount,
		hookTerm:   opt.Web_HookTerm,
	}

	event.registHandler()

	return event
}

func (e *EventWeb) registHandler() {
	if e.url[0] != "" {
		var nilTime time.Time = time.Time{}
		var cntExecutedTime time.Time

		if e.hookCount > 0 {
			e.eventQueue.AddListener(uint32(e.hookCount), func(q queue.Queue, args ...interface{}) {
				var items []KlevrEvent
				var iq = q

				logger.Debugf("event queue count : %d", iq.Length())

				for iq.Length() > 0 {
					items = append(items, *(iq.Pop().(*KlevrEvent)))
				}

				logger.Debugf("%+v", items)
				option := KlevrEventOption{URL: e.url, Web_HookCount: e.hookCount, Web_HookTerm: e.hookTerm}
				e.sendBulkEvent(&items, option)

				cntExecutedTime = time.Now().UTC()
			})
		}

		if e.hookTerm > 0 {
			baseTime := time.Duration(e.hookTerm) * time.Second
			sleepTime := baseTime

			go func() {
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

					logger.Debugf("Webhook scheduler event count : %d", e.eventQueue.Length())

					for e.eventQueue.Length() > 0 {
						items = append(items, *(e.eventQueue.Pop().(*KlevrEvent)))
					}

					logger.Debugf("%+v", items)

					if len(items) > 0 {
						option := KlevrEventOption{URL: e.url, Web_HookCount: e.hookCount, Web_HookTerm: e.hookTerm}
						e.sendBulkEvent(&items, option)
					}

					sleepTime = baseTime
				}
			}()
		}
	}
}

func (e *EventWeb) Close() {}

// AddEvent add klevr event for webhook
func (e *EventWeb) AddEvent(event *KlevrEvent) {
	logger.Debugf("add event : [%+v]", *event)

	//manager := common.BaseContext.Get(CtxServer).(*KlevrManager)
	//hookConfig := manager.Config.Server.Webhook

	//logger.Debugf("option(hookConfig) : [%+v]", option)

	if e.url[0] == "" {
		return
	}

	if e.hookCount <= 1 && e.hookTerm < 1 {
		option := KlevrEventOption{URL: e.url, Web_HookCount: e.hookCount, Web_HookTerm: e.hookTerm}
		go e.sendSingleEvent(event, option)
	} else {
		e.eventQueue.Push(event)
	}
}

func (e *EventWeb) AddEvents(events *[]KlevrEvent) {
	//manager := common.BaseContext.Get(CtxServer).(*KlevrManager)
	//hookConfig := manager.Config.Server.Webhook

	//logger.Debugf("option(hookConfig) : [%+v]", option)

	if e.url[0] == "" {
		return
	}

	option := KlevrEventOption{URL: e.url, Web_HookCount: e.hookCount, Web_HookTerm: e.hookTerm}
	go e.sendBulkEvent(events, option)
}

func (e *EventWeb) sendSingleEvent(event *KlevrEvent, option KlevrEventOption) {
	var arr = []KlevrEvent{*event}

	logger.Debugf("%+v", *event)
	logger.Debugf("%d", len(arr))

	e.sendBulkEvent(&arr, option)
}

func (e *EventWeb) sendBulkEvent(events *[]KlevrEvent, option KlevrEventOption) {
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

	res, err := http.Post(option.URL[0], "application/json", bytes.NewReader(b))

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
