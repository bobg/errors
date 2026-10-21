// Package errors is a drop-in replacement for the stdlib errors package.
// It adds error-wrapping functions, stack traces, error-tree traversal.
package errors

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
)

// Newf creates a new error with the given format string and arguments,
// formatted as with [fmt.Errorf].
// The result is a wrapped error containing the message and a stack trace.
// The format specifier %w works as in [fmt.Errorf].
func Newf(format string, args ...any) error {
	return addStack(errorf(format, args...))
}

// Errorf is a synonym for [Newf].
func Errorf(format string, args ...any) error {
	return addStack(errorf(format, args...))
}

// Errors formatted by fmt.Errorf with one or more %w verbs implement Unwrap() error or Unwrap() []error,
// depending on whether one or multiple errors are wrapped.
// To support Rewrap (for Prune) on errors created by Newf, Errorf, and Wrapf,
// errorf wraps the result of fmt.Errorf in either singleFormatted or multiFormatted,
// retaining the original format string and arguments.
//
// An argument list may contain multiple error values where only a subset are wrapped with %w
// (for example, "%s: %w", err1, err2), or where positional verbs reorder them (for example, "%[2]w: %[1]v").
// Rather than parsing the format string to figure out which arguments correspond to %w verbs,
// Rewrap lazily calls taggedErrorf when a wrapped error actually changes.
// taggedErrorf replaces every error in args with an indexedError recording its slice index,
// passes those tagged args to fmt.Errorf, and unwraps the result to discover which arg indices
// were captured by %w. It then substitutes the new error(s) at those indices and re-runs errorf.

type indexedError struct {
	idx int
	err error
}

func (e indexedError) Error() string {
	return e.err.Error()
}

func errorf(format string, args ...any) error {
	raw := fmt.Errorf(format, args...)

	switch raw := raw.(type) {
	case Unwrapper:
		return &singleFormatted{
			raw:    raw,
			format: format,
			args:   slices.Clone(args),
		}

	case MultiUnwrapper:
		return &multiFormatted{
			raw:    raw,
			format: format,
			args:   slices.Clone(args),
		}
	}

	return raw
}

func hasErrorArg(args []any) bool {
	for _, arg := range args {
		if _, ok := arg.(error); ok {
			return true
		}
	}
	return false
}

// taggedErrorf runs fmt.Errorf with all error arguments wrapped in indexedError,
// allowing the caller to unwrap the resulting error and determine which positions
// in args were wrapped via %w.
func taggedErrorf(format string, args []any) error {
	taggedArgs := make([]any, len(args))
	for i, arg := range args {
		if err, ok := arg.(error); ok {
			taggedArgs[i] = indexedError{idx: i, err: err}
		} else {
			taggedArgs[i] = arg
		}
	}
	return fmt.Errorf(format, taggedArgs...)
}

// singleFormatted wraps an error produced by fmt.Errorf with a single %w verb,
// retaining the format string and arguments so it can implement [Rewrapper].
type singleFormatted struct {
	raw    Unwrapper
	format string
	args   []any
}

func (f *singleFormatted) Error() string {
	return f.raw.Error()
}

func (f *singleFormatted) Unwrap() error {
	return f.raw.Unwrap()
}

func (f *singleFormatted) Rewrap(errs ...error) error {
	err := maybeJoin(errs...)

	if err == f.raw.Unwrap() {
		return f
	}

	var (
		probe = taggedErrorf(f.format, f.args).(Unwrapper)
		idx   = probe.Unwrap().(indexedError).idx
		args  = slices.Clone(f.args)
	)

	args[idx] = err
	return errorf(f.format, args...)
}

// multiFormatted wraps an error produced by fmt.Errorf with multiple %w verbs,
// retaining the format string and arguments so it can implement [Rewrapper].
type multiFormatted struct {
	raw    MultiUnwrapper
	format string
	args   []any
}

func (f *multiFormatted) Error() string {
	return f.raw.Error()
}

func (f *multiFormatted) Unwrap() []error {
	return f.raw.Unwrap()
}

func (f *multiFormatted) Rewrap(errs ...error) error {
	uu := f.raw.Unwrap()
	if len(errs) != len(uu) {
		// If the number of wrapped errors changed (e.g., one was pruned away entirely),
		// the original format string can no longer be used with the same %w verbs,
		// so fall back to Join.
		return Join(errs...)
	}
	changed := false
	for i, u := range uu {
		if errs[i] != u {
			changed = true
			break
		}
	}
	if !changed {
		return f
	}
	probe := taggedErrorf(f.format, f.args).(MultiUnwrapper)
	args := slices.Clone(f.args)
	for i, u := range probe.Unwrap() {
		args[u.(indexedError).idx] = errs[i]
	}
	return errorf(f.format, args...)
}

type wrapped struct {
	err error
	msg string
}

func (w *wrapped) Error() string {
	return w.msg + ": " + w.err.Error()
}

func (w *wrapped) Unwrap() error {
	return w.err
}

func (w *wrapped) Rewrap(errs ...error) error {
	err := maybeJoin(errs...)

	if err == w.err {
		return w
	}
	return &wrapped{
		err: err,
		msg: w.msg,
	}
}

type stacked struct {
	err   error
	stack []uintptr
}

func (s *stacked) Error() string {
	return s.err.Error()
}

func (s *stacked) Unwrap() error {
	return s.err
}

func (s *stacked) Stack() []uintptr {
	return s.stack
}

func (s *stacked) Rewrap(errs ...error) error {
	err := maybeJoin(errs...)

	if err == s.err {
		return s
	}

	return &stacked{
		err:   err,
		stack: s.stack,
	}
}

// If errs is empty or contains only nils,
// return nil.
// If it contains a single non-nil element,
// return that.
// Otherwise use [Join] to combine the errors.
func maybeJoin(errs ...error) error {
	var (
		n   int
		err error
	)

	for _, e := range errs {
		if e == nil {
			continue
		}
		if n > 0 {
			return Join(errs...)
		}
		err = e
		n++
	}

	return err
}

type stackerError interface {
	error
	Stacker
}

func addStack(err error) error {
	if _, ok := AsType[stackerError](err); ok {
		// err already contains a stack trace, no need to further decorate it
		return err
	}

	const maxdepth = 32
	var pcs [maxdepth]uintptr

	// 3 skips runtime.Callers, addStack (this function),
	// and the Wrap/Wrapf/New/Newf/Errorf call that got us here.
	n := runtime.Callers(3, pcs[:])

	return &stacked{
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

	s, ok := AsType[stackerError](e)
	if !ok {
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
	return addStack(&wrapped{err: err, msg: msg})
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
	if !hasErrorArg(args) {
		f := format
		return addStack(&wrapped{
			err: err,
			msg: fmt.Sprintf(f, args...),
		})
	}
	return addStack(errorf(format+": %w", append(args, err)...))
}
