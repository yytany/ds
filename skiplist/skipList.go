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

type skipList struct {
	rd          *rand.Rand
	levelCh     chan int                   //创建结点时获取已经创建好的层数序列
	length      int                        //结点数量，不包含头结点
	maxLevel    int                        //最大层数
	probability float64                    //层数生成概率
	compare     func(a, b interface{}) int //实现比较接口的函数
	head, tail  *skipList                  //头尾结点
}

type skipListNode struct {
	prev      *skipListNode //前置结点
	level     []levelNode   //层数
	condition interface{}   //比较条件
	data      interface{}   //数据
}

type levelNode struct {
	next *skipListNode //下一个结点
	span int           //到下一个结点的跨度
}

//生成层数
func (sl *skipList) levelGenerate() {
	for {
		level := 1
		for level < defaultMaxLevel &&
			defaultProbability <= sl.rd.Float64() {
			level++
		}
		sl.levelCh <- level
	}
}

//获取生成的层数
func (sl *skipList) getLevelGenerate() int {
	return <-sl.levelCh
}

//初始化一个跳表。需要实现 score 的比较接口,以便实现升序或者降序跳表
func New(compare func(a, b interface{}) int, options ...Option) (*skipList, error) {
	sl := &skipList{
		rd:          rand.New(rand.NewSource(time.Now().UnixNano())),
		levelCh:     make(chan int, defaultLevelCacheSize),
		length:      0,
		maxLevel:    defaultMaxLevel,
		probability: defaultProbability,
		compare:     compare,
	}
	for k := range options {
		if err := options[k](sl); err != nil {
			return sl, err
		}
	}
	go sl.levelGenerate()
	return sl, nil
}
