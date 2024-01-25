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

// 跳表
type SkipList struct {
	rd              *rand.Rand
	allowSameKey    bool          //是否允许存在相同的key  默认允许
	levelCh         chan int      //创建结点时获取已经创建好的层数序列
	length          int           //结点数量，不包含头结点
	constMaxLevel   int           //能生成的最大层数
	currentMaxLevel int           //当前的最大层数
	probability     float64       //层数生成概率
	compareAble     CompareAble   //需要实现比较接口
	head, tail      *skipListNode //头尾结点
}

// 跳表结点
type skipListNode struct {
	prev  *skipListNode //前置结点
	level []levelNode   //层数
	key   interface{}   //比较条件
	data  interface{}   //数据
}

// 跳表层结点
type levelNode struct {
	next *skipListNode //下一个结点
	span int           //到下一个结点的跨度
}

// 生成层数
func (sl *SkipList) levelGenerate() {
	for {
		level := 1
		for level < sl.constMaxLevel &&
			sl.probability <= sl.rd.Float64() {
			level++
		}
		sl.levelCh <- level
	}
}

// 生成新结点
func (sl *SkipList) nodeGenerate(key, data interface{}) *skipListNode {
	level := <-sl.levelCh
	if level-1 > sl.currentMaxLevel {
		sl.currentMaxLevel = level - 1
	}
	sl.length++
	return &skipListNode{
		prev:  nil,
		level: make([]levelNode, level),
		key:   key,
		data:  data,
	}
}

// 初始化一个跳表。需要实现 key 的比较接口,以便实现升序或者降序跳表
func New(compareAble CompareAble, options ...Option) (*SkipList, error) {
	sl := &SkipList{
		rd:              rand.New(rand.NewSource(time.Now().UnixNano())),
		allowSameKey:    true,
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

// 初始化头结点   (头结点仅映射层数，不存储数据)
func (sl *SkipList) headNodeInit() {
	sl.head = &skipListNode{
		prev:  nil,
		level: make([]levelNode, sl.constMaxLevel),
		key:   nil,
		data:  nil,
	}
}

// 更新当前最大层数
func (sl *SkipList) updateCurrentMaxLevel(currentLevel int) {
	if sl.length > 0 {
		for level := currentLevel; level >= 0; level-- {
			if sl.head.level[level].next != nil {
				sl.currentMaxLevel = level
				return
			}
		}
	}
	sl.currentMaxLevel = 0
}

// 获取所有相等结点
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

// 获取相等的第一个
func (sl *SkipList) searchFirstOneByKey(key interface{}) *skipListNode {
	node := sl.searchRandOneByKey(key)
	if node != nil {
		for node.prev != nil && sl.equals(key, node.prev.key) {
			node = node.prev
		}
	}
	return node
}

// 获取相等的末尾一个
func (sl *SkipList) searchTailOneByKey(key interface{}) *skipListNode {
	node := sl.searchRandOneByKey(key)
	if node != nil {
		for node.level[0].next != nil && sl.equals(key, node.level[0].next.key) {
			node = node.level[0].next
		}
	}
	return node
}

// 获取任意一个,只要找到相等的就返回
func (sl *SkipList) searchRandOneByKey(key interface{}) *skipListNode {
	if sl.length > 0 {
		preNode := sl.head
		for level := sl.currentMaxLevel; level >= 0; level-- {
			for ; ; preNode = preNode.level[level].next {
				if preNode.level[level].next == nil || sl.greaterThan(preNode.level[level].next.key, key) {
					if preNode != sl.head && sl.equals(preNode.key, key) {
						return preNode
					}
					break
				}
			}
		}
	}
	return nil
}

// 获取相同key的rank值，
func (sl *SkipList) searchRandNodeAndRankByKey(key interface{}) (*skipListNode, int) {
	if sl.length > 0 {
		currentRank := 0
		preNode := sl.head
		for level := sl.currentMaxLevel; level >= 0; level-- {
			for ; ; preNode = preNode.level[level].next {
				if preNode.level[level].next == nil || sl.greaterThan(preNode.level[level].next.key, key) {
					if preNode != sl.head && sl.equals(preNode.key, key) {
						return preNode, currentRank
					}
					break
				}
				currentRank += preNode.level[level].span
			}
		}
	}
	return nil, -1
}

// 获取相等的第一个
func (sl *SkipList) searchFirstNodeAndRankByKey(key interface{}) (*skipListNode, int) {
	node, rank := sl.searchRandNodeAndRankByKey(key)
	if node != nil {
		for node.prev != nil && sl.equals(key, node.prev.key) {
			node = node.prev
			rank--
		}
	}
	return node, rank
}

// 获取相等的第一个
func (sl *SkipList) searchTailNodeAndRankByKey(key interface{}) (*skipListNode, int) {
	node, rank := sl.searchRandNodeAndRankByKey(key)
	if node != nil {
		for node.level[0].next != nil && sl.equals(key, node.level[0].next.key) {
			node = node.level[0].next
			rank++
		}
	}
	return node, rank
}

// 通过顺位排序搜索   顺位 1~n
func (sl *SkipList) searchByRankRange(start, end int) []*skipListNode {
	list := []*skipListNode{}
	if start < 1 {
		start = 1
	}
	if start > end || start > sl.length || sl.length == 0 {
		return list
	}
	if start == sl.length {
		list = append(list, sl.tail)
	} else if start == 1 {
		for node := sl.head.level[0].next; node != nil && start <= end && start <= sl.length; start, node = start+1, node.level[0].next {
			list = append(list, node)
		}
	} else if end == sl.length {
		for node := sl.tail; node != nil && start <= end && end >= 1; end, node = end-1, node.prev {
			list = append(list, node)
		}
		sl.reverse(list)
	} else if start < sl.length {
		node := sl.head.level[0].next
		if start > 1 {
			node = sl.searchByRank(start)
		}
		for ; node != nil && start <= end && start <= sl.length; start, node = start+1, node.level[0].next {
			list = append(list, node)
		}
	}
	return list
}

// 通过精确rank搜索
func (sl *SkipList) searchByRank(rk int) *skipListNode {
	if rk > 0 && rk <= sl.length {
		if rk == 1 {
			return sl.head.level[0].next
		} else if rk == sl.length {
			return sl.tail
		}
		currentRank := 0
		preNode := sl.head
		for level := sl.currentMaxLevel; level >= 0; level-- {
			for ; ; preNode = preNode.level[level].next {
				if preNode.level[level].next == nil || preNode.level[level].span+currentRank > rk {
					if currentRank == rk {
						return preNode
					}
					break
				}
				currentRank += preNode.level[level].span
			}
		}
	}
	return nil
}

// 通过key批量更新
func (sl *SkipList) updateBatchByKey(key, data interface{}) bool {
	list := sl.searchAllByKey(key)
	for k := range list {
		sl.updateByNode(list[k], data)
	}
	return len(list) > 0
}

// 通过key更新  无重复key下更新成功
func (sl *SkipList) updateByKey(key, data interface{}) bool {
	if node := sl.searchRandOneByKey(key); node != nil {
		if (node.prev == nil || !sl.equals(node.key, node.prev.key)) &&
			(node.level[0].next == nil || !sl.equals(node.key, node.level[0].next.key)) {
			sl.updateByNode(node, data)
			return true
		}
	}
	return false
}

// 通过key删除  无重复key时删除成功
func (sl *SkipList) deleteByKey(key interface{}) bool {
	if node := sl.searchRandOneByKey(key); node != nil {
		if (node.prev == nil || !sl.equals(node.key, node.prev.key)) &&
			(node.level[0].next == nil || !sl.equals(node.key, node.level[0].next.key)) {
			sl.delNode(node)
			return true
		}
	}
	return false
}

// 通过结点更新
func (sl *SkipList) updateByNode(node *skipListNode, data interface{}) {
	node.data = data
}

// 添加结点   如果不允许有相同结点的话，重复添加时会失败
func (sl *SkipList) addNode(key, data interface{}) (int, bool) {
	if !sl.allowSameKey && sl.searchRandOneByKey(key) != nil {
		return 0, false
	}
	addNode := sl.nodeGenerate(key, data)
	if sl.length == 1 { //generate +1 了
		for level := sl.currentMaxLevel; level >= 0; level-- {
			sl.head.level[level].next = addNode
			sl.head.level[level].span = 1
		}
		sl.tail = addNode
		return 1, true
	}
	prevL := make([]*skipListNode, len(addNode.level)) // [层数]前置结点
	nextL := make([]*skipListNode, len(addNode.level)) // [层数]后置结点
	nrm := map[*skipListNode]int{}                     // [结点:rank]
	nodeRank := 0                                      //当前结点rank
	var preNode *skipListNode = sl.head
	//找前置与后置结点，并记录rank
	for level := sl.currentMaxLevel; level >= 0; level-- {
		for {
			if preNode.level[level].next == nil || sl.greaterThan(preNode.level[level].next.key, key) {
				if len(addNode.level) <= level {
					if preNode.level[level].next != nil {
						preNode.level[level].span++
					}
				} else {
					prevL[level] = preNode
					nrm[preNode] = nodeRank
					if preNode.level[level].next != nil {
						nextL[level] = preNode.level[level].next
						nrm[preNode.level[level].next] = preNode.level[level].span + nodeRank + 1
					}
				}
				break
			} else {
				nodeRank += preNode.level[level].span
				preNode = preNode.level[level].next
			}
		}
	}
	//当前结点的rank
	nodeRank++
	//更新前后置指向结点及本结点span
	for level := len(addNode.level) - 1; level >= 0; level-- {
		prevL[level].level[level].span = nodeRank - nrm[prevL[level]]
		addNode.level[level].next = prevL[level].level[level].next
		prevL[level].level[level].next = addNode
		addNode.prev = prevL[level]
		if nextL[level] != nil {
			addNode.level[level].span = nrm[nextL[level]] - nodeRank
			nextL[level].prev = addNode
		}
	}
	//更新tail
	if sl.tail == nil || sl.tail.level[0].next != nil {
		sl.tail = addNode
	}
	return nodeRank, true
}

// 通过key删除结点 所有key相等的结点
func (sl *SkipList) delByKey(key interface{}) bool {
	if sl.allowSameKey {
		list := sl.searchAllByKey(key)
		if len(list) > 0 {
			for k := range list {
				sl.delNode(list[k])
			}
			return true
		}
	} else {
		node := sl.searchRandOneByKey(key)
		if node != nil {
			sl.delNode(node)
			return true
		}
	}
	return false
}

// 通过node删除结点
func (sl *SkipList) delNode(delNode *skipListNode) {
	defer func(sl *SkipList) {
		sl.length--
		if len(delNode.level)-1 >= sl.currentMaxLevel {
			sl.updateCurrentMaxLevel(sl.currentMaxLevel)
		}
	}(sl)
	if sl.tail == delNode {
		sl.tail = delNode.prev
	}

	//与当前key相等，但处于delNode的后面
	equalsNextKeyMap := map[*skipListNode]bool{}
	for delNodeNext := delNode.level[0].next; delNodeNext != nil && sl.equals(delNode.key, delNodeNext.key); delNodeNext = delNodeNext.level[0].next {
		equalsNextKeyMap[delNodeNext] = true
	}

	preNode := sl.head
	for level := sl.currentMaxLevel; level >= 0; level-- {
		for ; ; preNode = preNode.level[level].next {
			if preNode.level[level].next == nil || sl.greaterThan(preNode.level[level].next.key, delNode.key) ||
				(equalsNextKeyMap[preNode.level[level].next] && sl.equals(preNode.level[level].next.key, delNode.key)) {
				if preNode.level[level].next != nil {
					preNode.level[level].span--
				}
				break
			} else if preNode.level[level].next == delNode {
				preNode.level[level].next = delNode.level[level].next
				preNode.level[level].span += delNode.level[level].span - 1
				if delNode.level[level].next == nil {
					preNode.level[level].span = 0
				}
				if level == 0 && delNode.level[level].next != nil {
					delNode.level[level].next.prev = delNode.prev
				}
				break
			}
		}
	}
}

// a,b相同
func (sl *SkipList) equals(a, b interface{}) bool {
	return sl.compareAble.Compare(a, b) == 0
}

// a小于b
func (sl *SkipList) lessThan(a, b interface{}) bool {
	return sl.compareAble.Compare(a, b) == -1
}

// a大于b
func (sl *SkipList) greaterThan(a, b interface{}) bool {
	return sl.compareAble.Compare(a, b) == 1
}

// a小于等于b
func (sl *SkipList) lessOrEquals(a, b interface{}) bool {
	return sl.compareAble.Compare(a, b) != 1
}

// a大于等于b
func (sl *SkipList) greaterOrEquals(a, b interface{}) bool {
	return sl.compareAble.Compare(a, b) != -1
}

// 翻转node
func (sl *SkipList) reverse(list []*skipListNode) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 { //让前序相等结点保持原顺序
		list[i], list[j] = list[j], list[i]
	}
}

/*
   对外实现
*/

// 获取结点数量
func (sl *SkipList) GetLength() int {
	return sl.length
}

// 获取第一个结点数据
func (sl *SkipList) GetFirst() interface{} {
	if sl.length > 0 {
		return sl.head.level[0].next.data
	}
	return nil
}

// 获取最后一个节点数据
func (sl *SkipList) GetTail() interface{} {
	if sl.length > 0 {
		return sl.tail.data
	}
	return nil
}

// 通过key搜索相等的第一个结点数据
func (sl *SkipList) GetFirstByKey(key interface{}) interface{} {
	node := sl.searchFirstOneByKey(key)
	if node != nil {
		return node.data
	}
	return nil
}

// 通过key搜索相等的最后一个结点数据
func (sl *SkipList) GetTailByKey(key interface{}) interface{} {
	node := sl.searchTailOneByKey(key)
	if node != nil {
		return node.data
	}
	return nil
}

// 通过key搜索相等的某一个结点数据 有重复key的结点则返回任意一个结点数据
func (sl *SkipList) GetRandByKey(key interface{}) interface{} {
	node := sl.searchRandOneByKey(key)
	if node != nil {
		return node.data
	}
	return nil
}

// 通过key搜索所有结点数据  返回所有结点数据
func (sl *SkipList) GetAllByKey(key interface{}) []interface{} {
	list := sl.searchAllByKey(key)
	data := make([]interface{}, len(list))
	for k := range list {
		data[k] = list[k].data
	}
	return data
}

// 获取指定key的任意相等结点数据及所在的排位  重复结点key则返回任意一个结点数据
func (sl *SkipList) GetRandWithRankByKey(key interface{}) (interface{}, int) {
	node, rk := sl.searchRandNodeAndRankByKey(key)
	if node != nil {
		return node.data, rk
	}
	return nil, rk
}

// 获取指定key的第一个相等结点数据及所在的排位
func (sl *SkipList) GetFirstWithRankByKey(key interface{}) (interface{}, int) {
	node, rk := sl.searchFirstNodeAndRankByKey(key)
	if node != nil {
		return node.data, rk
	}
	return nil, rk
}

// 获取指定key的最后一个相等结点数据及所在的排位
func (sl *SkipList) GetTailWithRankByKey(key interface{}) (interface{}, int) {
	node, rk := sl.searchTailNodeAndRankByKey(key)
	if node != nil {
		return node.data, rk
	}
	return nil, rk
}

// 获取指定排位的数据
func (sl *SkipList) GetByRank(rk int) interface{} {
	node := sl.searchByRank(rk)
	if node != nil {
		return node.data
	}
	return nil
}

// 获取指定排位区间的数据
func (sl *SkipList) GetByRankRange(start, end int) []interface{} {
	list := sl.searchByRankRange(start, end)
	data := make([]interface{}, len(list))
	for k := range list {
		data[k] = list[k].data
	}
	return data
}

// 更新所有和key相同的数据 所有相同的都会被更新 (更新结点数大于0时返回true)
func (sl *SkipList) UpdateBatchByKey(key, data interface{}) bool {
	return sl.updateBatchByKey(key, data)
}

// 更新和key相同的数据  当只有一个相同key的结点数据时能更新成功
func (sl *SkipList) UpdateByKey(key, data interface{}) bool {
	return sl.updateByKey(key, data)
}

// 更新指定排名的数据
func (sl *SkipList) UpdateByRank(rank int, data interface{}) bool {
	node := sl.searchByRank(rank)
	if node != nil {
		sl.updateByNode(node, data)
		return true
	}
	return false
}

// 删除所有和key相同的数据
func (sl *SkipList) DeleteBatchByKey(key interface{}) bool {
	return sl.delByKey(key)
}

// 删除和key相同的数据  当只有一个相同key的结点数据时能删除成功
func (sl *SkipList) DeleteByKey(key interface{}) bool {
	return sl.deleteByKey(key)
}

// 删除指定排位的结点
func (sl *SkipList) DeleteByRank(rank int) bool {
	node := sl.searchByRank(rank)
	if node != nil {
		sl.delNode(node)
		return true
	}
	return false
}

/*
插入数据
在 设置了  WithAllowTheSameKey(false)
即 allowSameKey == false 时,不允许有重复key时，重复的key添加将会返回false
返回当前排名和插入结果
*/
func (sl *SkipList) Insert(key, data interface{}) (int, bool) {
	return sl.addNode(key, data)
}
