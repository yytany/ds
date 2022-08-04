package skiplist

import (
	"math/rand"
	"time"
)

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
	sln := &skipListNode{
		prev:  nil,
		level: make([]levelNode, sl.constMaxLevel),
		key:   nil,
		data:  nil,
	}
	sl.head = sln
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

//通过条件搜索
func (sl *SkipList) searchByKey(key interface{}) {
	if sl.length == 0 {
		return
	}

}

//通过顺位搜索
func (sl *SkipList) searchByRank(start, end int)

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
