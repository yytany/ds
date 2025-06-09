package limit

import (
	"errors"
	"sync"
	"time"
)

// 滑动窗口限流
type RingWindowLimiter struct {
	mu         sync.Mutex
	windowSize time.Duration //	窗口大小
	maxCount   int           // 窗口最大请求数
	timestamps []time.Time   // 固定大小环形数组
	start      int           // 当前窗口内最早请求的位置
	count      int           // 当前窗口内有效请求数
}

func NewRingWindowLimiter(windowSize time.Duration, maxCount int) *RingWindowLimiter {
	return &RingWindowLimiter{
		windowSize: windowSize,
		maxCount:   maxCount,
		timestamps: make([]time.Time, maxCount),
	}
}

func (l *RingWindowLimiter) Allow() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-l.windowSize)

	// 清理过期请求：时间戳 <= windowStart
	for l.count > 0 && !l.timestamps[l.start].After(windowStart) {
		l.start = (l.start + 1) % l.maxCount
		l.count--
	}

	if l.count >= l.maxCount {
		return errors.New("request rate limit exceeded")
	}

	// 插入新请求
	insertIndex := (l.start + l.count) % l.maxCount
	l.timestamps[insertIndex] = now
	l.count++
	return nil
}
