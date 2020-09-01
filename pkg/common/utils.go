package common

import (
	"container/list"
	"sync"
)

// Queue interface for Queue
// 확장성을 위해 인터페이스를 정의하여 사용
type Queue interface {
	AddListener(uint32, func(*Queue, ...interface{}))
	Close()
	IsClosed() bool
	Length() uint64
	Pop() interface{}
	Push(interface{})
	ResetListenerCallCount()
}

// channel을 이용한 queue 구현체
type channelQueue struct {
	buf               chan interface{}             // queue에 삽입하는 데이터를 연결하는 채널
	current           chan *queueItem              // pointer가 가르키는 현재 데이터
	last              *queueItem                   // queue에 삽입된 마지막 데이터
	length            uint64                       // queue에 삽입된 데이터 건수
	alive             bool                         // queue의 종료 여부 (alive 가 false 인 경우 더 이상 queue를 사용할 수 없음)
	listener          func(*Queue, ...interface{}) // listenerCallCount에 해당하는 데이터 건수가 쌓이 때마다 호출 되는 리스너 함수
	listenerCallCount uint32                       // 리스너를 호출하는 기준 데이터 건수
	listenerRunCount  uint32                       // 리스너를 기준 데이터 건수에 호출하기 위해 카운트 하는 변수
}

// queue 데이터 구현체
type queueItem struct {
	next *queueItem  // 다음 데이터
	item interface{} // 실제 큐에 삽입된 데이터
}

// NewChannelQueue creator for queue struct
// Must call Close() method When the use of the queue is compelete.
func NewChannelQueue(chanBufSize uint32) *Queue {
	q := &channelQueue{
		buf:     make(chan interface{}, chanBufSize),
		current: make(chan *queueItem),
		length:  0,
		alive:   true, // alive가 false가 되면 더이상 큐를 사용할 수 없음 (잔여 데이터 갯수만큼 Pop()만 가능)
	}

	var Q Queue = q

	// Queue 생성 시 큐에 삽입되는 데이터를 채널로 받아 처리하기 위한 listener go routine. Close()가 호출되면 종료된다.
	go func() {
		// alive가 false가 될 때까지 반복 처리
		for q.alive {
			// select로 buf 채널을 수신
			select {
			case newItem := <-q.buf:
				// buf에 nil이 들어오면 queue가 종료된다. (Close() 를 통해 nil을 전달 받아 종료시킨다.)
				if newItem != nil {
					nq := &queueItem{
						item: newItem,
						next: nil,
					}

					if q.length == 0 {
						// 빈 큐에 데이터가 삽입될 시 출력 채널과 데이터가 동기화 되어 go routine이 block 되므로 새로운 go routine에서 채널을 전송한다.
						go func() {
							q.current <- nq
						}()
					} else {
						q.last.next = nq
					}

					q.last = nq
					q.length++

					if q.listener != nil {
						q.listenerRunCount++

						// 리스너 함수가 설정되고 호출 건수가 만족되면 리스너 함수를 호출한다.
						if q.listenerRunCount >= q.listenerCallCount {
							// 리스너 함수는 별도의 go routine으로 호출되므로 호출 시점의 누적 건수와 실행 시점의 누적 건수는 차이가 발생할 수 있다.
							go q.listener(&Q)
							q.listenerRunCount = 0
						}
					}
				}
			}
		}
	}()

	return &Q
}

// AddListener callCount 건수의 데이터가 신규로 큐에 삽입되는 시점에 등록한 리스너 함수를 호출한다.
// 리스너 함수는 별도의 go routine에서 수행되며, callCount 개수의 누적 건수가 발생할 시점에 호출되지만 비동기로 호출되므로 실행 시점의 누적 건수는 차이가 발생할 수 있다.
func (q *channelQueue) AddListener(callCount uint32, f func(*Queue, ...interface{})) {
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
func (q *channelQueue) Close() {
	q.alive = false
	q.buf <- nil
	close(q.buf)
}

// IsClosed 큐의 close 상태 여부를 반환한다.
func (q *channelQueue) IsClosed() bool {
	return q.alive
}

func (q *channelQueue) Push(item interface{}) {
	if !q.alive {
		panic("The queue has closed and can no longer be used.")
	}

	if item == nil {
		panic("Cannot insert nil into the Queue.")
	}

	q.buf <- item
}

func (q *channelQueue) Pop() interface{} {
	item := <-q.current

	if item != nil {
		go func() {
			if item.next != nil {
				q.current <- item.next
			}
		}()

		q.length--

		return item.item
	}

	return nil
}

func (q *channelQueue) Length() uint64 {
	return q.length
}

func (q *channelQueue) ResetListenerCallCount() {
	q.listenerCallCount = 0
}

// mutex를 이용한 queue 구현체
type mutexQueue struct {
	sync.RWMutex
	item              *list.List
	listener          func(*Queue, ...interface{})
	listenerCallCount uint32
	listenerRunCount  uint32
	alive             bool
}

// NewMutexQueue creator for queue struct
// Must call Close() method When the use of the queue is compelete.
func NewMutexQueue() *Queue {
	q := &mutexQueue{
		item:  list.New(),
		alive: true,
	}

	var Q Queue = q

	return &Q
}

// AddListener callCount 건수의 데이터가 신규로 큐에 삽입되는 시점에 등록한 리스너 함수를 호출한다.
// 리스너 함수는 별도의 go routine에서 수행되며, callCount 개수의 누적 건수가 발생할 시점에 호출되지만 비동기로 호출되므로 실행 시점의 누적 건수는 차이가 발생할 수 있다.
func (q *mutexQueue) AddListener(callCount uint32, f func(*Queue, ...interface{})) {
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
			go q.listener(&Q)
			q.listenerRunCount = 0
		}
	}
}

func (q *mutexQueue) Pop() interface{} {
	q.Lock()
	defer q.Unlock()

	item := q.item.Back()

	if item != nil {
		q.item.Remove(item)
		return item.Value
	}

	return nil
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
