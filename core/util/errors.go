package util

import (
	"errors"
	"fmt"
)

// triggers exit code 1 without any error messaging
var ErrSilent = errors.New("silent error")

// user-initiated cancellation
var ErrCancel = errors.New("cancel error")

// nothing failed but something is pending
var ErrPending = errors.New("pending error")

// FlagErrorf returns a new FlagError that wraps an error produced by
// fmt.Errorf(format, args...).
func FlagErrorf(format string, args ...interface{}) error {
	return FlagErrorWrap(fmt.Errorf(format, args...))
}

// FlagError returns a new FlagError that wraps the specified error.
func FlagErrorWrap(err error) error { return &FlagError{err} }

// A *FlagError indicates an error processing command-line flags or other arguments.
// Such errors cause the application to display the usage message.
type FlagError struct {
	// Note: not struct{error}: only *FlagError should satisfy error.
	err error
}

func (fe *FlagError) Error() string {
	return fe.err.Error()
}

func (fe *FlagError) Unwrap() error {
	return fe.err
}

func MutuallyExclusive(message string, conditions ...bool) error {
	numTrue := 0
	for _, ok := range conditions {
		if ok {
			numTrue++
		}
	}
	if numTrue > 1 {
		return FlagErrorf("%s", message)
	}
	return nil
}

type NoResultsError struct {
	message string
}

func (e NoResultsError) Error() string {
	return e.message
}

func NewNoResultsError(message string) NoResultsError {
	return NoResultsError{message: message}
}