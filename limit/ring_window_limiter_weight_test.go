package limit

import (
	"sync"
	"testing"
	"time"
)

func TestNewRingWindowLimiterWeight(t *testing.T) {
	windowSize := time.Second
	maxWeight := 10
	limiter := NewRingWindowLimiterWeight(windowSize, maxWeight)

	if limiter.windowSize != windowSize {
		t.Errorf("Expected windowSize %v, got %v", windowSize, limiter.windowSize)
	}
	if limiter.maxWeight != maxWeight {
		t.Errorf("Expected maxWeight %d, got %d", maxWeight, limiter.maxWeight)
	}
	if len(limiter.entries) != maxWeight {
		t.Errorf("Expected entries length %d, got %d", maxWeight, len(limiter.entries))
	}
}

func TestAllow_InvalidWeight(t *testing.T) {
	limiter := NewRingWindowLimiterWeight(time.Second, 10)

	err := limiter.Allow(0)
	if err == nil || err.Error() != "weight must be positive" {
		t.Errorf("Expected error for zero weight, got %v", err)
	}

	err = limiter.Allow(-1)
	if err == nil || err.Error() != "weight must be positive" {
		t.Errorf("Expected error for negative weight, got %v", err)
	}
}

func TestAllow_WithinLimit(t *testing.T) {
	limiter := NewRingWindowLimiterWeight(time.Second, 10)

	// Test single request
	err := limiter.Allow(5)
	if err != nil {
		t.Errorf("Expected request to be allowed, got error: %v", err)
	}

	// Test multiple requests within limit
	err = limiter.Allow(4)
	if err != nil {
		t.Errorf("Expected request to be allowed, got error: %v", err)
	}

	// Should be at 9/10 now
	err = limiter.Allow(2)
	if err == nil || err.Error() != "request rate limit exceeded (weight too large)" {
		t.Errorf("Expected rate limit exceeded, got %v", err)
	}
}

func TestAllow_ExactLimit(t *testing.T) {
	limiter := NewRingWindowLimiterWeight(time.Second, 10)

	// Fill exactly to the limit
	err := limiter.Allow(10)
	if err != nil {
		t.Errorf("Expected request to be allowed, got error: %v", err)
	}

	// Next request should be rejected
	err = limiter.Allow(1)
	if err == nil || err.Error() != "request rate limit exceeded (weight too large)" {
		t.Errorf("Expected rate limit exceeded, got %v", err)
	}
}

func TestAllow_ExpiredEntries(t *testing.T) {
	limiter := NewRingWindowLimiterWeight(100*time.Millisecond, 10)

	// First request
	err := limiter.Allow(5)
	if err != nil {
		t.Errorf("Expected request to be allowed, got error: %v", err)
	}

	// Wait for first request to expire
	time.Sleep(150 * time.Millisecond)

	// Second request should be allowed since first expired
	err = limiter.Allow(6)
	if err != nil {
		t.Errorf("Expected request to be allowed, got error: %v", err)
	}

	// Third request should be rejected (6 + 5 > 10)
	err = limiter.Allow(5)
	if err == nil || err.Error() != "request rate limit exceeded (weight too large)" {
		t.Errorf("Expected rate limit exceeded, got %v", err)
	}
}

func TestAllow_WrapAround(t *testing.T) {
	limiter := NewRingWindowLimiterWeight(time.Second, 3) // Small window for testing

	// Fill the ring buffer
	err := limiter.Allow(1) // index 0
	if err != nil {
		t.Fatal(err)
	}
	err = limiter.Allow(1) // index 1
	if err != nil {
		t.Fatal(err)
	}
	err = limiter.Allow(1) // index 2
	if err != nil {
		t.Fatal(err)
	}

	// Wait for first entry to expire
	time.Sleep(1100 * time.Millisecond)

	// Next request should wrap around to index 0
	err = limiter.Allow(1)
	if err != nil {
		t.Errorf("Expected request to be allowed after wrap around, got error: %v", err)
	}
	// Wait for first entry to expire
	time.Sleep(1000 * time.Millisecond)
	// Next request should wrap around to index 0
	err = limiter.Allow(1)
	if err != nil {
		t.Errorf("Expected request to be allowed after wrap around, got error: %v", err)
	}
}

func TestAllow_Concurrent(t *testing.T) {
	limiter := NewRingWindowLimiterWeight(time.Second, 100)
	var wg sync.WaitGroup
	successCount := 0
	errorCount := 0
	lock := &sync.Mutex{}

	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := limiter.Allow(1)
			lock.Lock()
			if err == nil {
				successCount++
			} else {
				errorCount++
			}
			lock.Unlock()
		}()
	}

	wg.Wait()
	if successCount != 100 {
		t.Errorf("Expected 100 successful requests, got %d", successCount)
	}
	if errorCount != 100 {
		t.Errorf("Expected 100 failed requests, got %d", errorCount)
	}
}
