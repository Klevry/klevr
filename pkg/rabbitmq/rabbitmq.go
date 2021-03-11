package rabbitmq

import (
	"sync"
	"time"

	"github.com/pkg/errors"

	"sync/atomic"

	"github.com/NexClipper/logger"
	"github.com/streadway/amqp"
)

const delay = 3 // reconnect after delay seconds

// Connection amqp.Connection wrapper
type Connection struct {
	*amqp.Connection
	wg *sync.WaitGroup
}

// Channel wrap amqp.Connection.Channel, get a auto reconnect channel
func (c *Connection) Channel() (*Channel, error) {
	c.wg.Wait()
	ch, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	channel := &Channel{
		Channel: ch,
	}

	go func() {
		for {
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok || channel.IsClosed() {
				logger.Debug("channel closed")
				channel.Close() // close again, ensure closed flag set when connection closed
				break
			}
			logger.Debugf("channel closed, reason: %v", reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for connection reconnect
				time.Sleep(delay * time.Second)

				ch, err := c.Connection.Channel()
				if err == nil {
					logger.Debug("channel recreate success")
					channel.Channel = ch
					break
				}

				logger.Debugf("channel recreate failed, err: %v", err)
			}
		}

	}()

	return channel, nil
}

// Dial wrap amqp.Dial, dial and get a reconnect connection
func Dial(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Error("Connection failed.")
		return nil, err
	}

	connection := &Connection{
		Connection: conn,
		wg:         &sync.WaitGroup{},
	}

	go func() {
		defer connection.wg.Done()

		for {
			reason, ok := <-connection.Connection.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok {
				logger.Debug("connection closed")
				break
			}
			logger.Infof("connection closed, reason: %v", reason)

			connection.wg.Add(1)

			// reconnect if not closed by developer
			for {
				// wait 1s for reconnect
				time.Sleep(delay * time.Second)

				conn, err := amqp.Dial(url)
				if err == nil {
					connection.Connection = conn
					logger.Infof("reconnect success")
					break
				}

				logger.Infof("reconnect failed, err: %v", err)
			}

			connection.wg.Done()
		}
	}()

	return connection, nil
}

// DialCluster with reconnect
func DialCluster(urls []string) (*Connection, error) {
	var connection *Connection
	nodeSequence := 0
	count := len(urls)

	for i, node := range urls {
		conn, err := amqp.Dial(node)

		if err != nil {
			if i < count-1 {
				err = nil
				continue
			} else {
				return nil, errors.Wrap(err, "all connection failed.")
			}
		}

		connection = &Connection{
			Connection: conn,
			wg:         &sync.WaitGroup{},
		}

		logger.Infof("ampq connection opened - url : %s", node)

		nodeSequence = i

		break
	}

	go func(urls []string, seq *int) {
		defer connection.wg.Done()

		var oldseq = *seq

		for {
			reason, ok := <-connection.Connection.NotifyClose(make(chan *amqp.Error))
			if !ok {
				logger.Debug("connection closed")
				break
			}
			logger.Debugf("connection closed, reason: %v", reason)

			connection.wg.Add(1)

			// reconnect with another node of cluster
			for {
				newSeq := next(urls, *seq)
				*seq = newSeq

				if oldseq == newSeq {
					time.Sleep(delay * time.Second)
				}

				conn, err := amqp.Dial(urls[newSeq])
				if err == nil {
					connection.Connection = conn
					logger.Infof("reconnect success - url : %s", urls[newSeq])
					break
				}

				logger.Warningf("reconnect failed - url : %s, err: %v", urls[newSeq], err)
			}

			connection.wg.Done()
		}
	}(urls, &nodeSequence)

	return connection, nil
}

// Next element index of slice
func next(s []string, lastSeq int) int {
	length := len(s)
	if length == 0 || lastSeq == length-1 {
		return 0
	} else if lastSeq < length-1 {
		return lastSeq + 1
	} else {
		return -1
	}
}

// Channel amqp.Channel wapper
type Channel struct {
	*amqp.Channel
	closed int32
}

// IsClosed indicate closed by developer
func (ch *Channel) IsClosed() bool {
	return (atomic.LoadInt32(&ch.closed) == 1)
}

// Close ensure closed flag set
func (ch *Channel) Close() error {
	if ch.IsClosed() {
		return amqp.ErrClosed
	}

	atomic.StoreInt32(&ch.closed, 1)

	return ch.Channel.Close()
}

// Consume warp amqp.Channel.Consume, the returned delivery will end only when channel closed by developer
func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			d, err := ch.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
			if err != nil {
				logger.Debugf("consume failed, err: %v", err)
				time.Sleep(delay * time.Second)
				continue
			}

			for msg := range d {
				deliveries <- msg
			}

			// sleep before IsClose call. closed flag may not set before sleep.
			time.Sleep(delay * time.Second)

			if ch.IsClosed() {
				break
			}
		}
	}()

	return deliveries, nil
}
