package queue

import "errors"

var (
	mustGreaterThanZeroErr = errors.New("the number must grater than 0")
	queueFullErr           = errors.New("queue is full")
	queueEmptyErr          = errors.New("queue is empty")
)
