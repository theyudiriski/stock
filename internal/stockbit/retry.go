package stockbit

import (
	"time"

	retry "github.com/avast/retry-go/v4"
)

// RetryOption configures WithRetry behavior.
type RetryOption func(*retryConfig)

type retryConfig struct {
	attempts uint
	delay    time.Duration
	onRetry  func(n uint, err error)
}

// Default retry settings for stockbit API calls.
const (
	DefaultAttempts = 5
	DefaultDelay    = 500 * time.Millisecond
)

// Attempts sets the number of retry attempts (default: 5).
func Attempts(n uint) RetryOption {
	return func(c *retryConfig) {
		c.attempts = n
	}
}

// Delay sets the delay between retries (default: 500ms).
func Delay(d time.Duration) RetryOption {
	return func(c *retryConfig) {
		c.delay = d
	}
}

// OnRetry sets a callback invoked on each retry (e.g. for logging).
func OnRetry(fn func(n uint, err error)) RetryOption {
	return func(c *retryConfig) {
		c.onRetry = fn
	}
}

// retryable wraps a function so it can be run with .WithRetry().
type retryable[T any] struct {
	fn func() (T, error)
}

// Retryable wraps fn so you can call .WithRetry() on it.
func Retryable[T any](fn func() (T, error)) retryable[T] {
	return retryable[T]{fn: fn}
}

// WithRetry runs the wrapped function with retries.
// Options default to DefaultAttempts and DefaultDelay if not set.
func (r retryable[T]) WithRetry(opts ...RetryOption) (T, error) {
	return WithRetry(r.fn, opts...)
}

// WithRetry runs fn with retries. Used as the parent for all stockbit API methods.
// Options default to DefaultAttempts and DefaultDelay if not set.
func WithRetry[T any](fn func() (T, error), opts ...RetryOption) (T, error) {
	cfg := &retryConfig{
		attempts: DefaultAttempts,
		delay:    DefaultDelay,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	retryOpts := []retry.Option{
		retry.Attempts(cfg.attempts),
		retry.Delay(cfg.delay),
	}
	if cfg.onRetry != nil {
		retryOpts = append(retryOpts, retry.OnRetry(cfg.onRetry))
	}

	return retry.DoWithData(fn, retryOpts...)
}
