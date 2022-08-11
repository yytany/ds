package skiplist

import (
	"errors"
	"math/rand"
)

var (
	levelErr       = errors.New("level must grater than 1")
	probabilityErr = errors.New("probability must grater than 0 and less than 1")
	randErr        = errors.New("*rand.Rand is nil")
	cacheErr       = errors.New("cache size  must grater than 0")
	cacheParamsErr = errors.New("cache params can not greater than cache size")
)

type Option func(*SkipList) error

//设置最大层数
func WithMaxLevel(level int) Option {
	return func(sl *SkipList) error {
		if level < 1 {
			return levelErr
		}
		sl.constMaxLevel = level
		return nil
	}
}

//设置层数生成概率
func WithProbability(probability float64) Option {
	return func(sl *SkipList) error {
		if probability <= 0 || probability >= 1 {
			return probabilityErr
		}
		sl.probability = probability
		return nil
	}
}

//设置随机数
func WithLevelRandSource(rd *rand.Rand) Option {
	return func(sl *SkipList) error {
		if rd == nil {
			return randErr
		}
		sl.rd = rd
		return nil
	}
}

//设置level缓冲区大小  缓冲区可以预先给定默认值，以达到初始定制层数
func WithLevelCacheSize(size int, params ...int) Option {
	return func(sl *SkipList) error {
		if size < 1 {
			return cacheErr
		}
		sl.levelCh = make(chan int, size)
		if len(params) > size {
			return cacheParamsErr
		}
		for k := range params {
			sl.levelCh <- params[k]
		}
		return nil
	}
}

//设置允许相同的key   如果为false，在插入相同key的值时将不生效
func WithAllowTheSameKey(allow bool) Option {
	return func(sl *SkipList) error {
		sl.allowSameKey = allow
		return nil
	}
}
