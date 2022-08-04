package queue

import (
	"fmt"
	"testing"
)

func Test_Loop(t *testing.T) {
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
