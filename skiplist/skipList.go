package skiplist

import (
	"math/rand"
	"time"
)

/*
	rank   顺位 从 1 ~ n
*/

const (
	defaultMaxLevel       = 32  //默认最大层数
	defaultProbability    = 0.5 //默认层数生成概率
	defaultLevelCacheSize = 8   //默认层数生成缓冲区大小
)

//跳表
type SkipList struct {
	rd              *rand.Rand
	levelCh         chan int      //创建结点时获取已经创建好的层数序列
	length          int           //结点数量，不包含头结点
	constMaxLevel   int           //能生成的最大层数
	currentMaxLevel int           //当前的最大层数
	probability     float64       //层数生成概率
	compareAble     CompareAble   //需要实现比较接口
	head, tail      *skipListNode //头尾结点
}

//跳表结点
type skipListNode struct {
	prev  *skipListNode //前置结点
	level []levelNode   //层数
	key   interface{}   //比较条件
	data  interface{}   //数据
}

//跳表层结点
type levelNode struct {
	next *skipListNode //下一个结点
	span int           //到下一个结点的跨度
}

//生成层数
func (sl *SkipList) levelGenerate() {
	for {
		level := 1
		for level < defaultMaxLevel &&
			defaultProbability <= sl.rd.Float64() {
			level++
		}
		sl.levelCh <- level
	}
}

//生成新结点
func (sl *SkipList) nodeGenerate(key, data interface{}) *skipListNode {
	level := <-sl.levelCh
	if level > sl.currentMaxLevel {
		sl.currentMaxLevel = level
	}
	sl.length++
	return &skipListNode{
		prev:  nil,
		level: make([]levelNode, level),
		key:   key,
		data:  data,
	}
}

//初始化一个跳表。需要实现 score 的比较接口,以便实现升序或者降序跳表
func New(compareAble CompareAble, options ...Option) (*SkipList, error) {
	sl := &SkipList{
		rd:              rand.New(rand.NewSource(time.Now().UnixNano())),
		levelCh:         make(chan int, defaultLevelCacheSize),
		length:          0,
		constMaxLevel:   defaultMaxLevel,
		currentMaxLevel: 0,
		probability:     defaultProbability,
		compareAble:     compareAble,
		head:            nil,
		tail:            nil,
	}
	for k := range options {
		if err := options[k](sl); err != nil {
			return sl, err
		}
	}
	sl.headNodeInit()
	go sl.levelGenerate()
	return sl, nil
}

//初始化头结点   (头结点仅映射层数，不存储数据)
func (sl *SkipList) headNodeInit() {
	sl.head = &skipListNode{
		prev:  nil,
		level: make([]levelNode, sl.constMaxLevel),
		key:   nil,
		data:  nil,
	}
}

//更新当前最大层数
func (sl *SkipList) updateCurrentMaxLevel() {
	if sl.length > 0 {
		for level := sl.constMaxLevel; level >= 0; level-- {
			if sl.head.level[level].next != nil {
				sl.currentMaxLevel = level
				return
			}
		}
	}
	sl.currentMaxLevel = 0
}

//获取所有相等结点
func (sl *SkipList) searchAllByKey(key interface{}) []*skipListNode {
	list := []*skipListNode{}
	if node := sl.searchRandOneByKey(key); node != nil {
		list = append(list, node)
		for preNode := node.prev; preNode != nil && sl.equals(key, preNode.key); preNode = preNode.prev {
			list = append(list, preNode)
		}
		sl.reverse(list)
		for nextNode := node.level[0].next; nextNode != nil && sl.equals(key, nextNode.key); nextNode = nextNode.level[0].next {
			list = append(list, nextNode)
		}
	}

	return list
}

//获取相等的第一个
func (sl *SkipList) searchFirstOneByKey(key interface{}) *skipListNode {
	node := sl.searchRandOneByKey(key)
	if node != nil {
		for node.prev != nil && sl.equals(key, node.prev.key) {
			node = node.prev
		}
	}
	return node
}

//获取相等的末尾一个
func (sl *SkipList) searchTailOneByKey(key interface{}) *skipListNode {
	node := sl.searchRandOneByKey(key)
	if node != nil {
		for node.level[0].next != nil && sl.equals(key, node.level[0].next) {
			node = node.level[0].next
		}
	}
	return node
}

//获取任意一个,只要找到相等的就返回
func (sl *SkipList) searchRandOneByKey(key interface{}) *skipListNode {
	if sl.length > 0 {
		var node *skipListNode
		for level := sl.currentMaxLevel; level >= 0; level-- {
			if node = sl.head.level[level].next; node != nil && sl.lessOrEquals(node.key, key) {
				for ; level >= 0; level-- {
					for {
						if sl.equals(node.key, key) {
							return node
						} else if node.level[level].next != nil && sl.lessOrEquals(node.level[level].next.key, key) {
							node = node.level[level].next
						} else {
							break
						}
					}
				}
			}
		}
	}
	return nil
}

//通过顺位排序搜索   顺位 1~n
func (sl *SkipList) searchByRankRange(start, end int) []*skipListNode {
	list := []*skipListNode{}
	if start > end || start < 1 || start > sl.length {
		return list
	}
	if start == sl.length {
		list = append(list, sl.tail)
	} else if end == sl.length {
		for node := sl.tail; start <= end && end >= 1; end, node = end-1, node.prev {
			list = append(list, node)
		}
		sl.reverse(list)
	} else if start < sl.length {
		node := sl.head.level[0].next
		if start > 1 {
			node = sl.searchByRank(start)
		}
		for ; start <= end && start <= sl.length && node != nil; start, node = start+1, node.level[0].next {
			list = append(list, node)
		}
	}
	return list
}

//通过精确rank搜索
func (sl *SkipList) searchByRank(rk int) *skipListNode {
	if rk > sl.length || rk < 1 {
		return nil
	}
	if rk == 1 {
		return sl.head.level[0].next
	} else if rk == sl.length {
		return sl.tail
	} else {
		for level := sl.currentMaxLevel; level >= 0; level-- {
			if node := sl.head.level[level].next; node != nil && sl.head.level[level].span <= rk {
				currentRank := sl.head.level[level].span
				for ; level >= 0; level-- {
					for {
						if currentRank == rk {
							return node
						} else if node.level[level].next != nil && node.level[level].span+currentRank <= rk {
							currentRank += node.level[level].span
							node = node.level[level].next
						} else {
							break
						}
					}
				}
			}
		}
	}
	return nil
}

//TODO 添加结点
func (sl *SkipList) addNode(key, data interface{}) {
	addNode := sl.nodeGenerate(key, data)
	if sl.length == 1 { //generate +1 了
		for level := sl.currentMaxLevel; level >= 0; level-- {
			sl.head.level[level].next = addNode
			sl.head.level[level].span = 1
		}
		sl.tail = addNode
		return
	}
	prevL := make([]*skipListNode, len(addNode.level)) // [层数]前置结点
	nextL := make([]*skipListNode, len(addNode.level)) // [层数]后置结点
	nrm := map[*skipListNode]int{}                     // [结点:rank]
	nodeRank := 0                                      //当前结点rank
	var preNode *skipListNode
	for level := sl.currentMaxLevel; level >= 0; level-- {
		if preNode = sl.head.level[level].next; preNode != nil {
			if sl.lessOrEquals(preNode.key, addNode.key) {
				nodeRank = sl.head.level[level].span
				for ; level >= 0; level-- {
					if sl.lessOrEquals(preNode.key, key) &&
						(preNode.level[level].next == nil || sl.greaterThan(preNode.level[level].next, key)) {
						addNode.prev = preNode
						addNode.level[level].next = preNode.level[level].next
						preNode.level[level].next = addNode
						addNode.level[level].next.prev = addNode
						nextL[level] = addNode.level[level].next
						nrm[preNode] = nodeRank
						nrm[addNode.level[level].next] = nodeRank + 1 + preNode.level[level].span
					} else if sl.lessOrEquals(preNode.key, key) {
						nodeRank += preNode.level[level].span
						preNode = preNode.level[level].next
					} else {
						prevL[level] = preNode
						nrm[preNode] = nodeRank
						preNode.level[level].next = addNode
					}
				}
			} else if level < len(addNode.level) {
				nextL[level] = preNode
				prevL[level] = sl.head
				nextL[level] = preNode
			}
		} else if level < len(addNode.level) {
			prevL[level] = sl.head
		}
	}

	//更新span

}

//更新结点span (从该结点向两边扩散)
func (sl *skipListNode) updateSpan(node *skipListNode) {

}

// 删除结点
func (sl *SkipList) delNode(delNode *skipListNode) {
	if sl.length == 0 || delNode == nil {
		return
	}
	defer func(sl *SkipList) {
		sl.length--
		if len(delNode.level) == sl.currentMaxLevel {
			sl.updateCurrentMaxLevel()
		}
	}(sl)
	if sl.tail == delNode {
		sl.tail = delNode.prev
	}
	var preNode *skipListNode
	for level := sl.currentMaxLevel; level >= 0; level-- {
		if preNode = sl.head.level[level].next; preNode != nil {
			if preNode == delNode {
				sl.head.level[level].next = delNode.level[level].next
				delNode.level[level].next.prev = nil
				sl.head.level[level].span += delNode.level[level].span - 1
				if delNode.level[level].next == nil {
					sl.head.level[level].span = 0
				}
			} else if sl.lessOrEquals(preNode.key, delNode.key) {
				for ; level >= 0; level-- {
				currentLevelNext:
					if preNode.level[level].next != nil {
						if preNode.level[level].next == delNode {
							preNode.level[level].next = delNode.level[level].next
							delNode.level[level].next.prev = preNode
							preNode.level[level].span += delNode.level[level].span - 1
							if delNode.level[level].next == nil {
								preNode.level[level].span = 0
							}
						} else {
							if sl.lessOrEquals(preNode.level[level].next.key, delNode.key) {
								preNode = preNode.level[level].next
								goto currentLevelNext
							} else {
								preNode.level[level].span--
							}
						}
					}
				}
			} else {
				sl.head.level[level].span--
			}
		}
	}
}

//a,b相同
func (sl *SkipList) equals(a, b interface{}) bool {
	return sl.compareAble.Compare(a, b) == 0
}

//a小于b
func (sl *SkipList) lessThan(a, b interface{}) bool {
	return sl.compareAble.Compare(a, b) == -1
}

//a大于b
func (sl *SkipList) greaterThan(a, b interface{}) bool {
	return sl.compareAble.Compare(a, b) == 1
}

//a小于等于b
func (sl *SkipList) lessOrEquals(a, b interface{}) bool {
	return sl.compareAble.Compare(a, b) != 1
}

//a大于等于b
func (sl *SkipList) greaterOrEquals(a, b interface{}) bool {
	return sl.compareAble.Compare(a, b) != -1
}

//翻转node
func (sl *SkipList) reverse(list []*skipListNode) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 { //让前序相等结点保持原顺序
		list[i], list[j] = list[j], list[i]
	}
}
