package queue

import (
	"container/list"
	"reflect"
	"sync"

	"github.com/NexClipper/logger"
)

// mutex를 이용한 queue 구현체
type mutexQueue struct {
	sync.RWMutex
	item              *list.List
	listener          func(Queue, ...interface{})
	listenerCallCount uint32
	listenerRunCount  uint32
	alive             bool
}

// NewMutexQueue creator for queue struct
// Must call Close() method When the use of the queue is complete.
func NewMutexQueue() Queue {
	q := &mutexQueue{
		item:  list.New(),
		alive: true,
	}

	var Q Queue = q

	return Q
}

// AddListener callCount 건수의 데이터가 신규로 큐에 삽입되는 시점에 등록한 리스너 함수를 호출한다.
// 리스너 함수는 별도의 go routine에서 수행되며, callCount 개수의 누적 건수가 발생할 시점에 호출되지만 비동기로 호출되므로 실행 시점의 누적 건수는 차이가 발생할 수 있다.
func (q *mutexQueue) AddListener(callCount uint32, f func(Queue, ...interface{})) {
	if callCount == 0 {
		panic("AddListener's parameter callCount must bigger than 0")
	}

	if f == nil {
		panic("AddListener's parameter listener f() can not be nil")
	}

	q.listenerCallCount = callCount
	q.listener = f
}

// Close 큐를 종료하여 더이상 큐에 데이터를 적재하지 못하게 한다.
// Close 후에도 이미 적재된 데이터는 Pop()을 통해 취득할 수 있다.
func (q *mutexQueue) Close() {
	q.alive = false
}

// IsClosed 큐의 close 상태 여부를 반환한다.
func (q *mutexQueue) IsClosed() bool {
	return q.alive
}

func (q *mutexQueue) BulkPush(items interface{}) {
	switch reflect.TypeOf(items).Kind() {
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		q.Lock()
		defer q.Unlock()

		if !q.alive {
			panic("The queue has closed and can no longer be used.")
		}

		if items == nil {
			panic("Cannot insert nil into the Queue.")
		}

		logger.Debugf("bulk push - [%+v]", items)

		arr := reflect.ValueOf(items)

		for i := 0; i < arr.Len(); i++ {
			q.item.PushBack(arr.Index(i).Interface())
		}

		if q.listener != nil {
			q.listenerRunCount++

			// 리스너 함수가 설정되고 호출 건수가 만족되면 리스너 함수를 호출한다.
			if q.listenerRunCount != 0 && q.listenerRunCount >= q.listenerCallCount {
				// 리스너 함수는 별도의 go routine으로 호출되므로 호출 시점의 누적 건수와 실행 시점의 누적 건수는 차이가 발생할 수 있다.
				var Q Queue = q
				go q.listener(Q)
				q.listenerRunCount = 0
			}
		}
	default:
		logger.Debugf("single push - [%+v]", items)
		q.Push(items)
	}
}

func (q *mutexQueue) Push(item interface{}) {
	q.Lock()
	defer q.Unlock()

	if !q.alive {
		panic("The queue has closed and can no longer be used.")
	}

	if item == nil {
		panic("Cannot insert nil into the Queue.")
	}

	q.item.PushBack(item)

	if q.listener != nil {
		q.listenerRunCount++

		// 리스너 함수가 설정되고 호출 건수가 만족되면 리스너 함수를 호출한다.
		if q.listenerRunCount != 0 && q.listenerRunCount >= q.listenerCallCount {
			// 리스너 함수는 별도의 go routine으로 호출되므로 호출 시점의 누적 건수와 실행 시점의 누적 건수는 차이가 발생할 수 있다.
			var Q Queue = q
			go q.listener(Q)
			q.listenerRunCount = 0
		}
	}
}

func (q *mutexQueue) Pop() interface{} {
	q.Lock()
	defer q.Unlock()

	item := q.item.Front()

	if item != nil {
		q.item.Remove(item)
		return item.Value
	}

	return nil
}

func (q *mutexQueue) PopAll() []interface{} {
	q.Lock()
	defer q.Unlock()

	length := uint64(q.item.Len())

	values := make([]interface{}, 0, length)

	for i := uint64(0); i < length; i++ {
		item := q.item.Front()

		if item == nil {
			break
		}

		q.item.Remove(item)
		values = append(values, item.Value)
	}

	return values
}

func (q *mutexQueue) Length() uint64 {
	q.Lock()
	defer q.Unlock()

	return uint64(q.item.Len())
}

func (q *mutexQueue) ResetListenerCallCount() {
	q.Lock()
	defer q.Unlock()

	q.listenerCallCount = 0
}
