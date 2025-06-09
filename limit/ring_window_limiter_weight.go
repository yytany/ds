package limit

import (
	"errors"
	"sync"
	"time"
)

// entry 用于在环形数组中存储单次请求的时间戳和权重
type entry struct {
	timestamp time.Time
	weight    int
}

// 滑动窗口权重限流
type RingWindowLimiterWeight struct {
	mu          sync.Mutex
	windowSize  time.Duration // 窗口大小
	maxWeight   int           // 窗口内允许的最大权重和
	entries     []entry       // 环形数组，长度 = maxWeight
	start       int           // 当前最早（还没过期）请求在 entries 中的索引
	count       int           // 当前环形队列中有效的元素个数（<= len(entries)）
	totalWeight int           // 窗口内所有未过期请求的权重和
}

// NewRingWindowLimiterWeight 构造函数：windowSize 表示滑动窗口时长，maxWeight 表示窗口内允许的最大权重和。
// 内部会把 entries 初始化为长度 maxWeight 的切片，这样在最糟糕的情况下（所有请求权重都为 1）也能存得下所有元素。
func NewRingWindowLimiterWeight(windowSize time.Duration, maxWeight int) *RingWindowLimiterWeight {
	return &RingWindowLimiterWeight{
		windowSize: windowSize,
		maxWeight:  maxWeight,
		entries:    make([]entry, maxWeight),
	}
}

// Allow 带权重地尝试放行一次请求。如果当前窗口内已存在的所有未过期请求权重和 + 新请求的 weight > maxWeight，则拒绝。
// 否则，将新请求的信息插入环形队列，并累加权重。
// weight 必须是 >=1 的整数。
func (l *RingWindowLimiterWeight) Allow(weight int) error {
	if weight <= 0 {
		return errors.New("weight must be positive")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-l.windowSize)

	// 1. 先清理过期请求：只要队列不空，且队首（索引为 l.start）的 timestamp <= windowStart，就把它移出
	for l.count > 0 {
		e := l.entries[l.start]
		if e.timestamp.After(windowStart) {
			// 队首尚未过期，跳出
			break
		}
		// 队首已过期，需要移出
		l.totalWeight -= e.weight
		l.start = (l.start + 1) % l.maxWeight
		l.count--
	}

	// 2. 检查是否能够放行：当前窗口内累加权重 + 新请求权重大于 maxWeight，就拒绝
	if l.totalWeight+weight > l.maxWeight {
		return errors.New("request rate limit exceeded (weight too large)")
	}

	// 3. 放行：把新请求插入环形数组中
	insertIndex := (l.start + l.count) % l.maxWeight
	l.entries[insertIndex] = entry{
		timestamp: now,
		weight:    weight,
	}
	l.count++
	l.totalWeight += weight

	return nil
}
