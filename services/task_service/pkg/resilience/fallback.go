package resilience

import (
	"context"
	"errors"
)

type FallbackFunc[T any] func(ctx context.Context, err error) (T, error)

type Fallback[T any] struct {
	primary   func(ctx context.Context) (T, error)
	fallbacks []FallbackFunc[T]
}

func NewFallback[T any](primary func(ctx context.Context) (T, error)) *Fallback[T] {
	return &Fallback[T]{
		primary:   primary,
		fallbacks: make([]FallbackFunc[T], 0),
	}
}

func (f *Fallback[T]) WithFallback(fallback FallbackFunc[T]) *Fallback[T] {
	f.fallbacks = append(f.fallbacks, fallback)
	return f
}

func (f *Fallback[T]) Execute(ctx context.Context) (T, error) {
	result, err := f.primary(ctx)
	if err == nil {
		return result, nil
	}

	for _, fallback := range f.fallbacks {
		result, fallbackErr := fallback(ctx, err)
		if fallbackErr == nil {
			return result, nil
		}
		err = errors.Join(err, fallbackErr)
	}

	var zero T
	return zero, err
}

type FallbackExecutor struct {
	circuitBreaker *CircuitBreaker
	retry          *Retry
}

func NewFallbackExecutor(cb *CircuitBreaker, retry *Retry) *FallbackExecutor {
	return &FallbackExecutor{
		circuitBreaker: cb,
		retry:          retry,
	}
}

func (fe *FallbackExecutor) Execute(ctx context.Context, fn func() error, fallback func() error) error {
	var execErr error

	cbErr := fe.circuitBreaker.Execute(func() error {
		execErr = fe.retry.Execute(ctx, fn)
		return execErr
	})

	if cbErr != nil {
		if errors.Is(cbErr, ErrCircuitOpen) && fallback != nil {
			return fallback()
		}
		return cbErr
	}

	if execErr != nil && fallback != nil {
		return fallback()
	}

	return execErr
}

func ExecuteWithFallback(fn func() error, fallback func() error) error {
	err := fn()
	if err != nil && fallback != nil {
		return fallback()
	}
	return err
}

type CachedFallback[T any] struct {
	cache    map[string]T
	fallback FallbackFunc[T]
}

func NewCachedFallback[T any]() *CachedFallback[T] {
	return &CachedFallback[T]{
		cache: make(map[string]T),
	}
}

func (cf *CachedFallback[T]) Set(key string, value T) {
	cf.cache[key] = value
}

func (cf *CachedFallback[T]) Get(key string) (T, bool) {
	val, ok := cf.cache[key]
	return val, ok
}

func (cf *CachedFallback[T]) GetOrFallback(ctx context.Context, key string, primary func() (T, error)) (T, error) {
	result, err := primary()
	if err == nil {
		cf.cache[key] = result
		return result, nil
	}

	if cached, ok := cf.cache[key]; ok {
		return cached, nil
	}

	var zero T
	return zero, err
}
