package errors_test

import (
	"fmt"
	"testing"

	"github.com/bobg/errors"
)

func BenchmarkNew(b *testing.B) {
	for b.Loop() {
		_ = errors.New("foo")
	}
}

func BenchmarkNewf(b *testing.B) {
	b.Run("simple", func(b *testing.B) {
		for b.Loop() {
			_ = errors.Newf("foo %d", 17)
		}
	})
	b.Run("wrap_plain", func(b *testing.B) {
		err := fmt.Errorf("bar")
		for b.Loop() {
			_ = errors.Newf("foo %d: %w", 17, err)
		}
	})
	b.Run("wrap_stack", func(b *testing.B) {
		err := errors.New("bar")
		for b.Loop() {
			_ = errors.Newf("foo %d: %w", 17, err)
		}
	})
}

func BenchmarkWrap(b *testing.B) {
	b.Run("plain", func(b *testing.B) {
		err := fmt.Errorf("foo")
		for b.Loop() {
			_ = errors.Wrap(err, "bar")
		}
	})
	b.Run("already_wrapped", func(b *testing.B) {
		err := errors.New("foo")
		for b.Loop() {
			_ = errors.Wrap(err, "bar")
		}
	})
}

func BenchmarkWrapf(b *testing.B) {
	b.Run("plain", func(b *testing.B) {
		err := fmt.Errorf("foo")
		for b.Loop() {
			_ = errors.Wrapf(err, "bar %d", 17)
		}
	})
	b.Run("already_wrapped", func(b *testing.B) {
		err := errors.New("foo")
		for b.Loop() {
			_ = errors.Wrapf(err, "bar %d", 17)
		}
	})
}

func BenchmarkJoin(b *testing.B) {
	e1 := errors.New("foo")
	e2 := errors.New("bar")
	for b.Loop() {
		_ = errors.Join(e1, e2)
	}
}
