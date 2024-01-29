package errors

import (
	stderrors "errors"
	"fmt"
	"runtime"
)

var ErrUnsupported = stderrors.ErrUnsupported

func As(err error, target interface{}) bool { return stderrors.As(err, target) }
func Is(err, target error) bool             { return stderrors.Is(err, target) }
func Join(errs ...error) error              { return stderrors.Join(errs...) }
func New(text string) error                 { return stderrors.New(text) }
func Unwrap(err error) error                { return stderrors.Unwrap(err) }

// Wrapped is a wrapped error containing a message and a stack trace.
type Wrapped struct {
	err   error
	msg   string
	stack []uintptr
}

func (w Wrapped) Unwrap() error { return w.err }

func (w Wrapped) Error() string {
	if w.msg == "" {
		return w.err.Error()
	}
	return w.msg + ": " + w.err.Error()
}

// Stack returns a slice of [Frame]
// representing the call stack at the point [Wrap] was called.
func (w Wrapped) Stack() []Frame {
	var (
		rframes = runtime.CallersFrames(w.stack)
		result  []Frame
	)

	for {
		rframe, more := rframes.Next()
		if rframe.PC == 0 {
			break
		}
		result = append(result, Frame{
			Function: rframe.Function,
			File:     rframe.File,
			Line:     rframe.Line,
		})
		if !more {
			break
		}
	}

	return result
}

// Frame is a stack frame.
type Frame struct {
	Function, File string
	Line           int
}

// Wrap creates a [Wrapped] error.
// It wraps the given error and attaches the given message, plus the call stack.
// If err is nil, Wrap returns nil.
func Wrap(err error, msg string) error {
	return dowrap(err, msg)
}

func dowrap(err error, msg string) error {
	if err == nil {
		return nil
	}

	const maxdepth = 32
	var pcs [maxdepth]uintptr
	n := runtime.Callers(3, pcs[:]) // 3 skips the call to runtime.Callers, the call to dowrap (this function), and the call to Wrap or Wrapf that got us here.
	return Wrapped{
		err:   err,
		msg:   msg,
		stack: pcs[:n],
	}
}

// Wrapf creates a [Wrapped] error.
// It is shorthand for Wrap(err, fmt.Sprintf(format, args...)).
func Wrapf(err error, format string, args ...interface{}) error {
	return dowrap(err, fmt.Sprintf(format, args...))
}
