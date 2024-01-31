# Errors - drop-in replacement for the stdlib errors package that adds wrapping and call stacks

[![Go Reference](https://pkg.go.dev/badge/github.com/bobg/errors.svg)](https://pkg.go.dev/github.com/bobg/errors)
[![Go Report Card](https://goreportcard.com/badge/github.com/bobg/errors)](https://goreportcard.com/report/github.com/bobg/errors)
[![Tests](https://github.com/bobg/errors/actions/workflows/go.yml/badge.svg)](https://github.com/bobg/errors/actions/workflows/go.yml)
[![Coverage Status](https://coveralls.io/repos/github/bobg/errors/badge.svg?branch=master)](https://coveralls.io/github/bobg/errors?branch=master)

This is errors,
a drop-in replacement for [the standard Go errors package](https://pkg.go.dev/errors).
It adds an API for wrapping an error with a caller-specified message and a snapshot of the call stack.

This module began as a fork of the venerable [github.com/pkg/errors](https://github.com/pkg/errors).
That module’s GitHub repository has been frozen for a long time and has not kept up with changes to the standard library.
(For example, it does not export [the Join function](https://pkg.go.dev/errors#Join) added in Go 1.20.)

This is now a brand-new implementation with a different API.
