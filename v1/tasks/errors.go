package tasks

import (
	"fmt"
	"time"
)

// ErrRetryTaskLater ...
type ErrRetryTaskLater struct {
	name, msg string
	retryIn   time.Duration
}

// RetryIn returns time.Duration from now when task should be retried
func (e ErrRetryTaskLater) RetryIn() time.Duration {
	return e.retryIn
}

// Error implements the error interface
func (e ErrRetryTaskLater) Error() string {
	return fmt.Sprintf("Task error: %s Will retry in: %s", e.msg, e.retryIn)
}

// NewErrRetryTaskLater returns new ErrRetryTaskLater instance
func NewErrRetryTaskLater(msg string, retryIn time.Duration) ErrRetryTaskLater {
	return ErrRetryTaskLater{msg: msg, retryIn: retryIn}
}

// Retriable is interface that retriable errors should implement
type Retriable interface {
	RetryIn() time.Duration
}

// ErrUnrecoverable is an unrecoverable error, ie. the task will not be
// retried.
type ErrUnrecoverable struct {
	error
}

// Unrecoverable wraps an error in `ErrUnrecoverable` struct
func Unrecoverable(err error) error {
	return ErrUnrecoverable{err}
}

// IsUnrecoverable checks if error is an instance of `ErrUnrecoverable`
func IsUnrecoverable(err error) bool {
	_, isUnrecoverable := err.(ErrUnrecoverable)
	return isUnrecoverable
}

// UnpackUnrecoverable checks if err is unrecoverable and returns the wrapped
// error.
func UnpackUnrecoverable(err error) (isUnrecoverable bool, werr error) {
	if unrecoverable, isUnrecoverable := err.(ErrUnrecoverable); isUnrecoverable {
		return true, unrecoverable.error
	}

	return false, err
}
