// Package errors is a drop-in replacement for the stdlib errors package.
// It adds error-wrapping functions, stack traces, error-tree traversal.
package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// Newf creates a new error with the given format string and arguments,
// formatted as with [fmt.Errorf].
// The result is a wrapped error containing the message and a stack trace.
// The format specifier %w works as in [fmt.Errorf].
func Newf(format string, args ...any) error {
	return dowrap(fmt.Errorf(format, args...))
}

// Errorf is a synonym for [Newf].
func Errorf(format string, args ...any) error {
	return Newf(format, args...)
}

type wrapped struct {
	err   error // This can be nil, in which case msg is the entire error message.
	stack []uintptr
}

func (w *wrapped) Unwrap() error {
	return w.err
}

func (w *wrapped) Error() string {
	return w.err.Error()
}

// Stack implements the [Stacker] interface.
func (w *wrapped) Stack() []uintptr {
	return w.stack
}

func dowrap(err error) error {
	var s Stacker
	if As(err, &s) {
		// err already contains a stack trace, no need to further decorate it
		return err
	}

	const maxdepth = 32
	var pcs [maxdepth]uintptr
	n := runtime.Callers(3, pcs[:]) // 3 skips runtime.Callers, dowrap (this function), and the Wrap/Wrapf/New/Newf call that got us here.
	return &wrapped{
		err:   err,
		stack: pcs[:n],
	}
}

// Stacker is the interface implemented by errors with a stack trace.
// Types wishing to implement Stacker can use [runtime.Callers] to get the call stack.
type Stacker interface {
	Stack() []uintptr
}

// Stack returns the stack trace from an error as a [Frames].
// If the error is nil or does not contain a stack trace,
// Stack returns nil.
func Stack(e error) Frames {
	if e == nil {
		return nil
	}

	var s Stacker
	if !As(e, &s) {
		return nil
	}

	var (
		pcs     = s.Stack()
		rframes = runtime.CallersFrames(pcs)
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

// Wrap creates a wrapped error.
// It wraps the given error and attaches the given message.
// It may also attach a stack trace,
// but only if err doesn't already contain one.
// If err is nil, Wrap returns nil.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return dowrap(fmt.Errorf("%s: %w", msg, err))
}

// Wrapf creates a wrapped error.
// It wraps the given error and attaches the message that results from formatting format with args.
// It may also attach a stack trace,
// but only if err doesn't already contain one.
// The format string may include %w specifiers, which work as in [fmt.Errorf].
// If any errors included via %w include a stack trace,
// that will also prevent Wrapf from including a new one.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}

	format += ": %w"
	args = append(args, err)
	return dowrap(fmt.Errorf(format, args...))
}
