# Errors - drop-in replacement for the stdlib errors package that adds wrapping and call stacks

[![Go Reference](https://pkg.go.dev/badge/github.com/bobg/errors.svg)](https://pkg.go.dev/github.com/bobg/errors)
[![Go Report Card](https://goreportcard.com/badge/github.com/bobg/errors)](https://goreportcard.com/report/github.com/bobg/errors)
[![Tests](https://github.com/bobg/errors/actions/workflows/go.yml/badge.svg)](https://github.com/bobg/errors/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/bobg/errors/badge.svg?branch=master)](https://coveralls.io/github/bobg/errors?branch=master)

This is errors,
a drop-in replacement for [the standard Go errors package](https://pkg.go.dev/errors).
It adds an API for wrapping an error with a caller-specified message and a snapshot of the call stack.
It also adds a function for traversing an error’s tree of wrapped sub-errors.

This module began as a fork of the venerable [github.com/pkg/errors](https://github.com/pkg/errors).
That module’s GitHub repository has been frozen for a long time and has not kept up with changes to the standard library.
(For example, it does not export [the Join function](https://pkg.go.dev/errors#Join) added in Go 1.20,
or [the AsType function](https://pkg.go.dev/errors#AsType) added in Go 1.26.)

This is now a brand-new implementation with a different API.

## Usage

When you write code that handles an error from a function call by returning the error to your caller,
it’s a good practice to “wrap” the error with context from your callsite first.
The Go standard library lets you do this like so:

```go
if err := doThing(...); err != nil {
  return fmt.Errorf("doing thing: %w", err)
}
```

The resulting error still “is” `err`
(in the sense of [errors.Is](https://pkg.go.dev/errors#Is))
but is now wrapped with context about how `err` was produced −
to wit, “doing thing.”

This library lets you do the same thing like this:

```go
if err := doThing(...); err != nil {
  return errors.Wrap(err, "doing thing")
}
```

This approach has a couple of benefits over wrapping with `fmt.Errorf` and `%w`.
First, as a convenience, `Wrap` returns `nil` if its error argument is `nil`.
So if `doThing` is the last thing that happens in your function,
you can safely write:

```go
err := doThing(...)
return errors.Wrap(err, "doing thing")
```

This will correctly return `nil` when `doThing` succeeds.

Second, errors produced by this package contain a snapshot of the call stack from the moment the error was created.
This can be retrieved with the `Stack` function.

Note that diligent error wrapping often makes the stack trace superfluous.
As errors work their way up the stack
it’s possible to decorate them with enough extra context
to understand exactly where and why the error occurred.
But you can’t always rely on diligent error wrapping,
and you can’t always wait for an error to propagate all the way up the call stack before reporting it.
