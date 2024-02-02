// Package errors is a drop-in replacement for the stdlib errors package.
// It adds an error-wrapping API and stack traces.
package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// Newf creates a new error with the given format string and arguments,
// formatted as with [fmt.Errorf].
// The result is a [Wrapped] error containing the message and a stack trace.
// The format specified %w works as in fmt.Errorf.
func Newf(format string, args ...any) error {
	return dowrap(fmt.Errorf(format, args...))
}

// Wrapped is a wrapped error containing a message and a stack trace.
type Wrapped struct {
	err   error // This can be nil, in which case msg is the entire error message.
	stack []uintptr
}

func (w *Wrapped) Unwrap() error {
	return w.err
}

func (w *Wrapped) Error() string {
	return w.err.Error()
}

// Stack returns a slice of [Frame]
// representing the call stack at the point [Wrap] was called.
func (w *Wrapped) Stack() Frames {
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

// Stack returns the stack trace from an error.
// If the error is nil or is not a [Wrapped] error,
// Stack returns nil.
func Stack(e error) Frames {
	if e == nil {
		return nil
	}

	var w *Wrapped
	if !As(e, &w) {
		return nil
	}
	return w.Stack()
}

// Frames is a slice of [Frame].
type Frames []Frame

func (fs Frames) String() string {
	var sb strings.Builder
	for _, f := range fs {
		fmt.Fprintln(&sb, f)
	}
	return sb.String()
}

// Frame is a stack frame.
type Frame struct {
	Function, File string
	Line           int
}

func (f Frame) String() string {
	if f.File != "" && f.Function != "" {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%s:", f.File)
		if f.Line != 0 {
			fmt.Fprintf(&sb, "%d:", f.Line)
		}
		fmt.Fprintf(&sb, " %s", f.Function)
		return sb.String()
	}

	if f.Function != "" {
		return f.Function
	}

	if f.File != "" && f.Line != 0 {
		return fmt.Sprintf("%s:%d: (unknown function)", f.File, f.Line)
	}

	return "(unknown)"
}

// Wrap creates a [Wrapped] error.
// It wraps the given error and attaches the given message, plus the call stack.
// If err is nil, Wrap returns nil.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return dowrap(fmt.Errorf("%s: %w", msg, err))
}

func dowrap(err error) error {
	var w *Wrapped
	if As(err, &w) {
		// err already contains a stack trace, no need to further decorate it
		return err
	}

	const maxdepth = 32
	var pcs [maxdepth]uintptr
	n := runtime.Callers(3, pcs[:]) // skip runtime.Callers, dowrap, and the Wrap/Wrapf/New/Newf call that got us here.
	return &Wrapped{
		err:   err,
		stack: pcs[:n],
	}
}

// Wrapf creates a [Wrapped] error.
// It is shorthand for Wrap(err, fmt.Sprintf(format, args...)).
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	format += ": %w"
	args = append(args, err)
	return dowrap(fmt.Errorf(format, args...))
}
