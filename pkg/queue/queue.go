package queue

// Queue interface for Queue
// 확장성을 위해 인터페이스를 정의하여 사용
type Queue interface {
	AddListener(uint32, func(Queue, ...interface{}))
	Close()
	IsClosed() bool
	Length() uint64
	Pop() interface{}
	Push(interface{})
	BulkPush(interface{})
	ResetListenerCallCount()
	PopAll() []interface{}
}
