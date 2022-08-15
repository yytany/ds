package queue

//循环双端队列
type LoopDeque struct {
	queue            []interface{} //数据
	headPtr, tailPtr int           //头尾指针， [head][data][tail] 分别指向数据的前后
	length           int           //当前存储长度
	isFirst          bool          //是否为第一个数据添加
}

func NewLoopDeque(k int) (*LoopDeque, error) {
	if k < 1 {
		return nil, mustGreaterThanZeroErr
	}
	return &LoopDeque{
		queue:   make([]interface{}, k),
		headPtr: 0,
		tailPtr: 0,
		length:  0,
		isFirst: true,
	}, nil
}

//头部入队列
func (q *LoopDeque) PushFront(value interface{}) error {
	if q.IsFull() {
		return queueFullErr
	}
	q.queue[q.headPtr] = value
	q.length++
	q.headPtr--
	if q.headPtr == -1 {
		q.headPtr = len(q.queue) - 1
	}
	if q.isFirst {
		q.tailPtr++
		q.isFirst = false
	}
	return nil
}

//尾部入队列
func (q *LoopDeque) PushTail(value interface{}) error {
	if q.IsFull() {
		return queueEmptyErr
	}
	q.queue[q.tailPtr] = value
	q.length++
	q.tailPtr++
	if q.tailPtr == len(q.queue) {
		q.tailPtr = 0
	}
	if q.isFirst {
		q.headPtr = len(q.queue) - 1
		q.isFirst = false
	}
	return nil
}

//头部出队列
func (q *LoopDeque) PopFront() (interface{}, error) {
	if q.IsEmpty() {
		return nil, queueEmptyErr
	}
	q.headPtr++
	q.length--
	if q.headPtr == len(q.queue) {
		q.headPtr = 0
	}
	return q.queue[q.headPtr], nil
}

//尾部出队列
func (q *LoopDeque) PopTail() (interface{}, error) {
	if q.IsEmpty() {
		return nil, queueEmptyErr
	}
	q.tailPtr--
	q.length--
	if q.tailPtr == -1 {
		q.tailPtr = len(q.queue) - 1
	}
	return q.queue[q.tailPtr], nil
}

//获取头部数据而不出队列
func (q *LoopDeque) GetFront() (interface{}, error) {
	if q.IsEmpty() {
		return nil, queueEmptyErr
	}
	if q.headPtr+1 == len(q.queue) {
		return q.queue[0], nil
	}
	return q.queue[q.headPtr+1], nil
}

//获取尾部数据而不出队列
func (q *LoopDeque) GetTail() (interface{}, error) {
	if q.IsEmpty() {
		return nil, queueEmptyErr
	}
	if q.tailPtr-1 == -1 {
		return q.queue[len(q.queue)-1], nil
	}
	return q.queue[q.tailPtr-1], nil
}

//队列是否空
func (q *LoopDeque) IsEmpty() bool {
	return q.length == 0
}

//队列是否满
func (q *LoopDeque) IsFull() bool {
	return q.length == len(q.queue)
}
