package skiplist

import "fmt"

type iterator struct {
	sl   *SkipList
	node *skipListNode
}

//跳表迭代器
func NewIterator(sl *SkipList) *iterator {
	return &iterator{
		sl:   sl,
		node: sl.head,
	}
}

//初始化node
func (it *iterator) InitHead() {
	it.SetNode(it.sl.head)
}

//获取level层的下一个
func (it *iterator) Next(level int) *skipListNode {
	if it.node != nil {
		it.node = it.node.level[level].next
	}
	return it.node
}

//设置当前结点
func (it *iterator) SetNode(node *skipListNode) {
	it.node = node
}

//获取当前结点
func (it *iterator) Node() *skipListNode {
	return it.node
}

//获取当前结点当前层的span
func (it *iterator) Span(level int) int {
	return it.node.level[level].span
}

//输出跳表结构，适用于小于3个字符的data
func (it *iterator) PrintGraph() {
	for level := it.sl.currentMaxLevel; level >= 0; level-- {
		fmt.Printf("%d |\t", level)
		for it.InitHead(); it.Node() != nil; it.Next(level) {
			if it.Node() != it.sl.head {
				fmt.Printf("%3v", it.Node().data)
			}
			for span := it.Span(level); span > 0; span-- {
				if span > 1 {
					fmt.Printf("---%3v", "---")
				} else {
					fmt.Printf("-->")
				}
			}
		}
		fmt.Println()
	}
}
