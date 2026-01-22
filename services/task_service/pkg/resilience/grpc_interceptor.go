package resilience

import (
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CircuitBreakerManager struct {
	mu       sync.RWMutex
	breakers map[string]*CircuitBreaker
	config   CircuitBreakerConfig
}

func NewCircuitBreakerManager(config CircuitBreakerConfig) *CircuitBreakerManager {
	return &CircuitBreakerManager{
		breakers: make(map[string]*CircuitBreaker),
		config:   config,
	}
}

func (m *CircuitBreakerManager) GetBreaker(name string) *CircuitBreaker {
	m.mu.RLock()
	cb, exists := m.breakers[name]
	m.mu.RUnlock()

	if exists {
		return cb
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if cb, exists = m.breakers[name]; exists {
		return cb
	}

	cfg := m.config
	cfg.Name = name
	cb = NewCircuitBreaker(cfg)
	m.breakers[name] = cb

	return cb
}

func UnaryClientInterceptorWithResilience(cbManager *CircuitBreakerManager, retryConfig RetryConfig) grpc.UnaryClientInterceptor {
	retry := NewRetry(retryConfig)

	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		cb := cbManager.GetBreaker(cc.Target())

		var lastErr error
		err := cb.Execute(func() error {
			retryErr := retry.Execute(ctx, func() error {
				invokeErr := invoker(ctx, method, req, reply, cc, opts...)
				if invokeErr != nil && isRetryableGRPCError(invokeErr) {
					return invokeErr
				}
				lastErr = invokeErr
				if invokeErr != nil {
					return nil
				}
				return nil
			})
			if retryErr != nil {
				return retryErr
			}
			return lastErr
		})

		if err != nil {
			return err
		}
		return lastErr
	}
}

func isRetryableGRPCError(err error) bool {
	if err == nil {
		return false
	}

	st, ok := status.FromError(err)
	if !ok {
		return true
	}

	switch st.Code() {
	case codes.Unavailable,
		codes.ResourceExhausted,
		codes.Aborted,
		codes.DeadlineExceeded:
		return true
	default:
		return false
	}
}

func isCircuitBreakerTrigger(err error) bool {
	if err == nil {
		return false
	}

	st, ok := status.FromError(err)
	if !ok {
		return true
	}

	switch st.Code() {
	case codes.Unavailable,
		codes.Internal,
		codes.Unknown:
		return true
	default:
		return false
	}
}

type ResilienceInterceptor struct {
	cbManager   *CircuitBreakerManager
	retry       *Retry
	timeout     time.Duration
}

func NewResilienceInterceptor(cbConfig CircuitBreakerConfig, retryConfig RetryConfig, timeout time.Duration) *ResilienceInterceptor {
	return &ResilienceInterceptor{
		cbManager:   NewCircuitBreakerManager(cbConfig),
		retry:       NewRetry(retryConfig),
		timeout:     timeout,
	}
}

func (ri *ResilienceInterceptor) UnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if ri.timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, ri.timeout)
			defer cancel()
		}

		cb := ri.cbManager.GetBreaker(cc.Target())

		return cb.Execute(func() error {
			return ri.retry.Execute(ctx, func() error {
				err := invoker(ctx, method, req, reply, cc, opts...)
				if err != nil && !isRetryableGRPCError(err) {
					return &nonRetryableError{err}
				}
				return err
			})
		})
	}
}

type nonRetryableError struct {
	err error
}

func (e *nonRetryableError) Error() string {
	return e.err.Error()
}

func (e *nonRetryableError) Unwrap() error {
	return e.err
}

func LoggingStateChangeCallback(name string, from, to CircuitState) {
	log.Printf("[CircuitBreaker] %s: state changed from %s to %s", name, from.String(), to.String())
}
