# Errors - drop-in replacement for the stdlib errors package that adds wrapping, call stacks, and more

[![Go Reference](https://pkg.go.dev/badge/github.com/bobg/errors.svg)](https://pkg.go.dev/github.com/bobg/errors)
[![Go Report Card](https://goreportcard.com/badge/github.com/bobg/errors)](https://goreportcard.com/report/github.com/bobg/errors)
[![Tests](https://github.com/bobg/errors/actions/workflows/go.yml/badge.svg)](https://github.com/bobg/errors/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/bobg/errors/badge.svg?branch=master)](https://coveralls.io/github/bobg/errors?branch=master)

This is errors,
a fork of [github.com/pkg/errors](https://github.com/pkg/errors).
That module was a drop-in replacement for [the standard Go errors package](https://pkg.go.dev/errors).
It added convenient wrapping,
call stacks,
and other features to error objects.
That repo has been archived for some time.
This fork tracks further developments in the standard library —
for example, [the Join function](https://pkg.go.dev/errors#Join) adding in Go 1.20 —
so that it can remain a drop-in replacement.

## License

BSD-2-Clause
