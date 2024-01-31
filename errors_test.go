package errors

import (
	"errors"
	"path/filepath"
	"slices"
	"testing"
)

func TestWrap(t *testing.T) {
	e := func() error {
		return Wrap(ErrUnsupported, "foo")
	}()

	t.Run("Is", func(t *testing.T) {
		if !Is(e, ErrUnsupported) {
			t.Errorf("got %v, want %v", e, ErrUnsupported)
		}
	})

	t.Run("ErrStr", func(t *testing.T) {
		const wantstr = "foo: unsupported operation"
		if got := e.Error(); got != wantstr {
			t.Errorf("got %q, want %q", got, wantstr)
		}
	})

	var w Wrapped
	if !errors.As(e, &w) {
		t.Fatalf("e is %T, want Wrapped", e)
	}

	s := w.Stack()
	for i := range s {
		s[i].File = filepath.Base(s[i].File)
		s[i].Line = 0
	}

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

	str := s.String()
	const wantstr = "errors_test.go: github.com/bobg/errors.TestWrap.func1\nerrors_test.go: github.com/bobg/errors.TestWrap\ntesting.go: testing.tRunner\n"
	if str != wantstr {
		t.Errorf("got %q, want %q", str, wantstr)
	}
}
