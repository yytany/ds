package skiplist

import (
	"errors"
	"math/rand"
)

var (
	levelErr       = errors.New("level must grater than 1")
	probabilityErr = errors.New("probability must grater than 0 and less than 1")
	randErr        = errors.New("*rand.Rand is nil")
	cacheErr       = errors.New("level cache must grater than 0")
)

type Option func(*skipList) error

//设置最大层数
func WithMaxLevel(level int) Option {
	return func(sl *skipList) error {
		if level < 1 {
			return levelErr
		}
		sl.maxLevel = level
		return nil
	}
}

//设置层数生成概率
func WithProbability(probability float64) Option {
	return func(sl *skipList) error {
		if probability <= 0 || probability >= 1 {
			return probabilityErr
		}
		sl.probability = probability
		return nil
	}
}

//设置随机数
func WithLevelRandSource(rd *rand.Rand) Option {
	return func(sl *skipList) error {
		if rd == nil {
			return randErr
		}
		sl.rd = rd
		return nil
	}
}

//设置level缓冲区大小
func WithLevelCacheSize(size int) Option {
	return func(sl *skipList) error {
		if size < 1 {
			return cacheErr
		}
		sl.levelCh = make(chan int, size)
		return nil
	}
}
