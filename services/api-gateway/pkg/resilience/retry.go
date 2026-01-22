package resilience

import (
	"context"
	"math"
	"math/rand"
	"time"
)

type RetryConfig struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	Jitter          float64
	RetryableErrors []error
	IsRetryable     func(error) bool
}

func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     10 * time.Second,
		Multiplier:      2.0,
		Jitter:          0.1,
	}
}

type Retry struct {
	config RetryConfig
}

func NewRetry(config RetryConfig) *Retry {
	if config.MaxAttempts == 0 {
		config.MaxAttempts = 3
	}
	if config.InitialInterval == 0 {
		config.InitialInterval = 100 * time.Millisecond
	}
	if config.MaxInterval == 0 {
		config.MaxInterval = 10 * time.Second
	}
	if config.Multiplier == 0 {
		config.Multiplier = 2.0
	}

	return &Retry{config: config}
}

func (r *Retry) Execute(ctx context.Context, fn func() error) error {
	var lastErr error
	interval := r.config.InitialInterval

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if !r.isRetryable(lastErr) {
			return lastErr
		}

		if attempt == r.config.MaxAttempts {
			break
		}

		waitTime := r.calculateWaitTime(interval)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
		}

		interval = time.Duration(float64(interval) * r.config.Multiplier)
		if interval > r.config.MaxInterval {
			interval = r.config.MaxInterval
		}
	}

	return lastErr
}

func (r *Retry) isRetryable(err error) bool {
	if r.config.IsRetryable != nil {
		return r.config.IsRetryable(err)
	}

	for _, retryableErr := range r.config.RetryableErrors {
		if err == retryableErr {
			return true
		}
	}

	return len(r.config.RetryableErrors) == 0
}

func (r *Retry) calculateWaitTime(interval time.Duration) time.Duration {
	if r.config.Jitter == 0 {
		return interval
	}

	jitter := r.config.Jitter * float64(interval)
	return interval + time.Duration(jitter*(rand.Float64()*2-1))
}

func RetryWithExponentialBackoff(ctx context.Context, maxAttempts int, fn func() error) error {
	retry := NewRetry(RetryConfig{
		MaxAttempts:     maxAttempts,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
		Jitter:          0.1,
	})
	return retry.Execute(ctx, fn)
}

func CalculateBackoff(attempt int, initialInterval, maxInterval time.Duration, multiplier float64) time.Duration {
	backoff := float64(initialInterval) * math.Pow(multiplier, float64(attempt-1))
	if time.Duration(backoff) > maxInterval {
		return maxInterval
	}
	return time.Duration(backoff)
}
