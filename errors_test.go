package errors

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/bobg/go-generics/v3/slices"
)

func TestWrap(t *testing.T) {
	e := func() error {
		return Wrap(ErrUnsupported, "foo")
	}()
	if !Is(e, ErrUnsupported) {
		t.Errorf("got %v, want %v", e, ErrUnsupported)
	}

	const wantstr = "foo: unsupported operation"
	if got := e.Error(); got != wantstr {
		t.Errorf("got %q, want %q", got, wantstr)
	}

	var w Wrapped
	if !errors.As(e, &w) {
		t.Fatalf("e is %T, want Wrapped", e)
	}

	s := w.Stack()
	s = slices.Map(s, func(f Frame) Frame {
		f.File = filepath.Base(f.File)
		f.Line = 0
		return f
	})
	if len(s) < 3 {
		t.Fatalf("got %d frames, want at least 3", len(s))
	}
	s = s[:3]
	want := []Frame{{
		Function: "github.com/bobg/errors.TestWrap.func1",
		File:     "errors_test.go",
	}, {
		Function: "github.com/bobg/errors.TestWrap",
		File:     "errors_test.go",
	}, {
		Function: "testing.tRunner",
		File:     "testing.go",
	}}
	if !slices.Equal(s, want) {
		t.Errorf("got %v, want %v", s, want)
	}
}
