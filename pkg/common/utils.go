package common

// Queue interface for Queue
// 확장성을 위해 인터페이스를 정의하여 사용
type Queue interface {
	AddListener(uint32, func(*Queue))
	Close()
	Length() uint64
	Pop() interface{}
	Push(interface{})
}

// queue 구현체
type queue struct {
	buf               chan interface{} // queue에 삽입하는 데이터를 연결하는 채널
	current           chan *queueItem  // pointer가 가르키는 현재 데이터
	last              *queueItem       // queue에 삽입된 마지막 데이터
	length            uint64           // queue에 삽입된 데이터 건수
	alive             bool             // queue의 종료 여부 (alive 가 false 인 경우 더 이상 queue를 사용할 수 없음)
	listener          func(*Queue)     // listenerCallCount에 해당하는 데이터 건수가 쌓이 때마다 호출 되는 리스너 함수
	listenerCallCount uint32           // 리스너를 호출하는 기준 데이터 건수
	listenerRunCount  uint32           // 리스너를 기준 데이터 건수에 호출하기 위해 카운트 하는 변수
}

// queue 데이터 구현체
type queueItem struct {
	next *queueItem  // 다음 데이터
	item interface{} // 실제 큐에 삽입된 데이터
}

// NewQueue creator for queue struct
// Must call Close() method When the use of the queue is compelete.
func NewQueue(chanBufSize uint32) *Queue {
	q := &queue{
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
						if q.listenerRunCount != 0 && q.listenerRunCount >= q.listenerCallCount {
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

//
func (q *queue) AddListener(callCount uint32, f func(*Queue)) {
	if callCount == 0 {
		panic("AddListener's parameter callCount must bigger than 0")
	}

	if f == nil {
		panic("AddListener's parameter listener f() can not be nil")
	}

	q.listenerCallCount = callCount
	q.listener = f
}

func (q *queue) Close() {
	q.alive = false
	q.buf <- nil
}

func (q *queue) Push(item interface{}) {
	if !q.alive {
		panic("The queue has closed and can no longer be used.")
	}

	if item == nil {
		panic("Cannot insert nil into the Queue.")
	}

	q.buf <- item
}

func (q *queue) Pop() interface{} {
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

func (q *queue) Length() uint64 {
	return q.length
}
