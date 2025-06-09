package limit

import (
	"sync"
	"testing"
	"time"
)

// TestAllowOverLimit 测试超限时的拒绝请求
func TestAllowOverLimit(t *testing.T) {
	limiter := NewRingWindowLimiter(time.Second, 2)
	for i := 0; i < 2; i++ {
		if err := limiter.Allow(); err != nil {
			t.Fatalf("前两次请求应允许，但被拒绝: %v", err)
		}
	}

	err := limiter.Allow()
	if err == nil || err.Error() != "request rate limit exceeded" {
		t.Fatalf("第三次请求应被拒绝，但错误不匹配。got: %v", err)
	}
}

// TestSameTimestampRequests 测试同一时间戳的多个请求
func TestSameTimestampRequests(t *testing.T) {
	limiter := NewRingWindowLimiter(time.Second, 3)
	now := time.Now()

	// 直接设置状态模拟已有3个请求
	limiter.mu.Lock()
	for i := 0; i < 3; i++ {
		limiter.timestamps[i] = now
	}
	limiter.count = 3
	limiter.mu.Unlock()

	err := limiter.Allow()
	if err == nil || err.Error() != "request rate limit exceeded" {
		t.Fatalf("应拒绝第4个请求，但错误不匹配。got: %v", err)
	}
}

// TestConcurrentAccess 测试并发安全性
func TestConcurrentAccess(t *testing.T) {
	limiter := NewRingWindowLimiter(100*time.Millisecond, 100)
	var wg sync.WaitGroup
	wg.Add(100)

	// 并发发送100个请求
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			if err := limiter.Allow(); err != nil {
				t.Errorf("并发请求被错误拒绝: %v", err)
			}
		}()
	}
	wg.Wait()

	// 检查第101个请求应被拒绝
	err := limiter.Allow()
	if err == nil || err.Error() != "request rate limit exceeded" {
		t.Fatalf("应拒绝第101个请求，但错误不匹配。got: %v", err)
	}
}

func TestRingWindowLimiter_BasicAllow(t *testing.T) {
	limiter := NewRingWindowLimiter(1*time.Second, 3)

	// 前3次应该都允许
	for i := 0; i < 3; i++ {
		if err := limiter.Allow(); err != nil {
			t.Errorf("expected allow, got error: %v", err)
		}
	}

	// 第4次应该被拒绝
	if err := limiter.Allow(); err == nil {
		t.Errorf("expected error, got allow")
	}
}

func TestRingWindowLimiter_ExpireRequests(t *testing.T) {
	limiter := NewRingWindowLimiter(100*time.Millisecond, 2)

	if err := limiter.Allow(); err != nil {
		t.Errorf("expected allow, got error: %v", err)
	}

	time.Sleep(150 * time.Millisecond) // 使第一个请求过期

	if err := limiter.Allow(); err != nil {
		t.Errorf("expected allow after expiration, got error: %v", err)
	}

	if err := limiter.Allow(); err != nil {
		t.Errorf("expected allow, got error: %v", err)
	}

	// 第三个请求应该被拒绝
	if err := limiter.Allow(); err == nil {
		t.Errorf("expected error, got allow")
	}
}

func TestRingWindowLimiter_ExactWindow(t *testing.T) {
	limiter := NewRingWindowLimiter(200*time.Millisecond, 2)

	_ = limiter.Allow()
	time.Sleep(100 * time.Millisecond)
	_ = limiter.Allow()
	time.Sleep(100 * time.Millisecond)

	// 此时第一个请求应过期，可以接受新请求
	if err := limiter.Allow(); err != nil {
		t.Errorf("expected allow after window shift, got error: %v", err)
	}
}

func TestRingWindowLimiter_WrapAround(t *testing.T) {
	limiter := NewRingWindowLimiter(500*time.Millisecond, 3)

	_ = limiter.Allow()
	time.Sleep(200 * time.Millisecond)
	_ = limiter.Allow()
	time.Sleep(200 * time.Millisecond)
	_ = limiter.Allow()

	time.Sleep(400 * time.Millisecond) // 使第一个过期

	if err := limiter.Allow(); err != nil {
		t.Errorf("expected allow after wrap-around, got error: %v", err)
	}
}

func TestRingWindowLimiter_Concurrency(t *testing.T) {
	limiter := NewRingWindowLimiter(1*time.Second, 100)
	var wg sync.WaitGroup
	var allowCount int
	var mu sync.Mutex

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := limiter.Allow(); err == nil {
				mu.Lock()
				allowCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	if allowCount > 100 {
		t.Errorf("allowed too many requests: %d", allowCount)
	} else {
		t.Logf("allowed %d requests as expected", allowCount)
	}
}
