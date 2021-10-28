package event

import (
	"encoding/json"

	"github.com/Klevry/klevr/pkg/rabbitmq"
	"github.com/NexClipper/logger"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

type MQHandle struct {
	Connection *rabbitmq.Connection
	Queue      *amqp.Queue
}

type EventMQ struct {
	mqHandle *MQHandle
}

func NewEventMQ(opt KlevrEventOption) EventManager {
	mqConn, err := rabbitmq.DialCluster(opt.URL)
	if err != nil {
		logger.Errorf("Failed to connect to MQ - %+v", errors.Cause(err))
		panic(err)
	}

	mqChannel, err := mqConn.Channel()
	if err != nil {
		logger.Errorf("Failed to open a channel to MQ - %+v", errors.Cause(err))
		panic(err)
	}

	queue, err := mqChannel.QueueDeclare(opt.MQ_Name, opt.MQ_Durable, opt.MQ_AutoDelete, false, false, nil)
	if err != nil {
		logger.Errorf("Failed to declare queue from MQ - %+v", errors.Cause(err))
		panic(err)
	}

	mqChannel.Close()

	return &EventMQ{mqHandle: &MQHandle{Connection: mqConn, Queue: &queue}}
}

func (e *EventMQ) Close() {
	e.mqHandle.Connection.Close()
}

// AddEvent add klevr event for webhook
func (e *EventMQ) AddEvent(event *KlevrEvent) {
	logger.Debugf("add event : [%+v]", *event)

	arr := []KlevrEvent{*event}
	go e.sendBulkEvent(&arr, KlevrEventOption{})
}

func (e *EventMQ) AddEvents(events *[]KlevrEvent) {
	go e.sendBulkEvent(events, KlevrEventOption{})

}

func (e *EventMQ) sendSingleEvent(event *KlevrEvent, option KlevrEventOption) {}

func (e *EventMQ) sendBulkEvent(events *[]KlevrEvent, option KlevrEventOption) {
	b, err := json.Marshal(*events)
	if err != nil {
		logger.Errorf("klevr event MQ publish marshal error - %+v", errors.Cause(err))
		retryFailedEvent(events, false)
	}

	channel, err := e.mqHandle.Connection.Channel()
	if err != nil {
		logger.Errorf("Failed to open a channel to MQ - %+v", errors.Cause(err))
		retryFailedEvent(events, true)
	}

	defer channel.Close()
	err = channel.Publish("", e.mqHandle.Queue.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: 2,
		Body:         b,
	})

	if err != nil {
		logger.Errorf("Failed to publish to MQ - %+v", errors.Cause(err))
		retryFailedEvent(events, true)
	}
}
