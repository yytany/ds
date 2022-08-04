package queue

//循环队列
type LoopQueue struct {
	queue            []interface{} //数据
	headPtr, tailPtr int           //头尾指针 头指针指向旧数据，尾指针指向新数据
	length           int           //当前存储长度
}

func NewLoopQueue(k int) (*LoopQueue, error) {
	if k < 1 {
		return nil, mustGreaterThanZeroErr
	}
	return &LoopQueue{
		queue:   make([]interface{}, k),
		headPtr: 0,
		tailPtr: -1,
		length:  0,
	}, nil
}

//入队列
func (q *LoopQueue) Push(value interface{}) error {
	if q.IsFull() {
		return queueFullErr
	}
	q.tailPtr++
	if q.tailPtr == len(q.queue) {
		q.tailPtr = 0
	}
	q.queue[q.tailPtr] = value
	q.length++
	return nil
}

//出队列
func (q *LoopQueue) Pop() (interface{}, error) {
	if q.IsEmpty() {
		return nil, queueEmptyErr
	}
	data := q.queue[q.headPtr]
	q.headPtr++
	if q.headPtr == len(q.queue) {
		q.headPtr = 0
	}
	q.length--
	return data, nil
}

//取队头数据而不出队列
func (q *LoopQueue) Front() (interface{}, error) {
	if q.IsEmpty() {
		return nil, queueEmptyErr
	}
	return q.queue[q.headPtr], nil
}

//取队尾数据而不出队列
func (q *LoopQueue) Tail() (interface{}, error) {
	if q.IsEmpty() {
		return nil, queueEmptyErr
	}
	return q.queue[q.tailPtr], nil
}

//是否为空
func (q *LoopQueue) IsEmpty() bool {
	return q.length == 0
}

//是否满了
func (q *LoopQueue) IsFull() bool {
	return q.length == len(q.queue)
}
