package queue

import (
	"fmt"
	"testing"
)

func Test_LoopQueue(t *testing.T) {
	lq, _ := NewLoopQueue(3)
	fmt.Println(lq.Push(1))
	fmt.Println(lq.Push(2))
	fmt.Println(lq.Push(3))
	fmt.Println(lq.Front())
	fmt.Println(lq.Tail())
	fmt.Println(lq.Pop())
	fmt.Println(lq.Pop())
	fmt.Println(lq.Pop())
	fmt.Println(lq.Pop())
	fmt.Println(lq.Front())
	fmt.Println(lq.Tail())
}

func Test_LoopDeque(t *testing.T) {
	lq, _ := NewLoopDeque(3)
	fmt.Println(lq.PushFront(7))
	fmt.Println(lq.PopTail())
	fmt.Println(lq.IsFull())
	fmt.Println(lq.IsEmpty())
	fmt.Println(lq.PushFront(7))
	fmt.Println(lq.PopFront())
	fmt.Println(lq.IsFull())
	fmt.Println(lq.IsEmpty())
	fmt.Println(lq.PushFront(7))
	fmt.Println(lq.PushFront(1))
	fmt.Println(lq.PushTail(2))
	fmt.Println(lq.PushFront(3))
	fmt.Println(lq.GetFront())
	fmt.Println(lq.GetTail())
	fmt.Println(lq.IsFull())
	fmt.Println(lq.IsEmpty())
	fmt.Println(lq.PopFront())
	fmt.Println(lq.PopTail())
	fmt.Println(lq.PopFront())
	fmt.Println(lq.PopTail())
	fmt.Println(lq.PopFront())
}
