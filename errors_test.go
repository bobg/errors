package errors_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/bobg/errors"
)

func TestAs(t *testing.T) {
	e := errors.New("foo")

	var s *errors.Stacked
	if !errors.As(e, &s) {
		t.Errorf("got %T, want *Stacked", e)
	}

	e2 := fmt.Errorf("bar")
	if errors.As(e2, &s) {
		t.Errorf("fmt.Errorf error should not be *Stacked")
	}
}

func TestAsType(t *testing.T) {
	e := errors.New("foo")

	if _, ok := errors.AsType[*errors.Stacked](e); !ok {
		t.Errorf("got %T, want *Stacked", e)
	}

	e2 := fmt.Errorf("bar")
	if _, ok := errors.AsType[*errors.Stacked](e2); ok {
		t.Errorf("fmt.Errorf error should not be *Stacked")
	}
}

func TestIs(t *testing.T) {
	e := errors.New("foo")
	e2 := errors.Wrap(e, "bar")

	if !errors.Is(e2, e) {
		t.Errorf("wrapped error fails Is test")
	}
	if !errors.Is(e, e) {
		t.Errorf("wrapped error fails Is-self test")
	}
}

func TestJoin(t *testing.T) {
	e := errors.Join(nil, nil)
	if e != nil {
		t.Errorf("got %v, want nil", e)
	}

	var (
		e1 = errors.New("foo")
		e2 = errors.New("bar")
	)

	e = errors.Join(e1, e2)
	if e.Error() != "foo\nbar" {
		t.Errorf("got %q, want %q", e.Error(), "foo\nbar")
	}

	if !errors.Is(e, e1) {
		t.Errorf("got false, want true")
	}
	if !errors.Is(e, e2) {
		t.Errorf("got false, want true")
	}
}

func TestNew(t *testing.T) {
	e := errors.New("foo")
	if e.Error() != "foo" {
		t.Errorf("got %q, want %q", e.Error(), "foo")
	}

	s := errors.Stack(e)
	if len(s) == 0 {
		t.Errorf("got %v, want non-empty", s)
	}
}

func TestNewf(t *testing.T) {
	t.Run("Newf", func(t *testing.T) {
		testNewf(t, errors.Newf)
	})
	t.Run("Errorf", func(t *testing.T) {
		testNewf(t, errors.Errorf)
	})
}

func testNewf(t *testing.T, f func(string, ...any) error) {
	e := f("foo %d", 17)
	if e.Error() != "foo 17" {
		t.Errorf("got %q, want %q", e.Error(), "foo 17")
	}

	s := errors.Stack(e)
	if len(s) == 0 {
		t.Errorf("got %v, want non-empty", s)
	}

	var (
		e1 = errors.New("bar")
		e2 = errors.New("baz")
	)
	e = f("foo %d: %w, %w", 17, e1, e2)
	if e.Error() != "foo 17: bar, baz" {
		t.Errorf("got %q, want %q", e.Error(), "foo 17: bar, baz")
	}
	if !errors.Is(e, e1) {
		t.Errorf("got false, want true")
	}
	if !errors.Is(e, e2) {
		t.Errorf("got false, want true")
	}

	e = f("dup: %[1]w, %[1]w", e1)
	if e.Error() != "dup: bar, bar" {
		t.Errorf("got %q, want %q", e.Error(), "dup: bar, bar")
	}
	if !errors.Is(e, e1) {
		t.Errorf("got false, want true")
	}
}

func TestUnwrap(t *testing.T) {
	var (
		e = errors.New("foo")
		u = errors.Unwrap(e)
	)

	if u == nil {
		t.Fatal("got nil, want non-nil")
	}
	if u.Error() != "foo" {
		t.Errorf("got %q, want %q", u.Error(), "foo")
	}

	u2 := errors.Unwrap(u)
	if u2 != nil {
		t.Errorf("got %v, want nil", u2)
	}
}

func TestWrap(t *testing.T) {
	e := errors.Wrap(nil, "foo")
	if e != nil {
		t.Errorf("got %v, want nil", e)
	}

	e = errors.New("foo")

	e2 := errors.Wrap(e, "bar")
	if e2.Error() != "bar: foo" {
		t.Errorf("got %q, want %q", e2.Error(), "bar: foo")
	}
	if !errors.Is(e2, e) {
		t.Errorf("got false, want true")
	}
	if _, ok := e2.(*errors.Stacked); ok {
		t.Errorf("double-stacked error should not be *Stacked")
	}
}

func TestWrapf(t *testing.T) {
	e := errors.Wrapf(nil, "foo")
	if e != nil {
		t.Errorf("got %v, want nil", e)
	}

	e = errors.New("foo")
	if got := errors.Wrapf(e, "bar %d", 17); got.Error() != "bar 17: foo" {
		t.Errorf("got %q, want %q", got.Error(), "bar 17: foo")
	}

	e2 := errors.New("bar")
	e3 := errors.Wrapf(e, "baz %d: %w", 17, e2)
	if e3.Error() != "baz 17: bar: foo" {
		t.Errorf("got %q, want %q", e3.Error(), "baz 17: bar: foo")
	}

	if !errors.Is(e3, e2) {
		t.Errorf("got false, want true")
	}
	if !errors.Is(e3, e) {
		t.Errorf("got false, want true")
	}
}

type emptyStacker struct{}

func (emptyStacker) Error() string    { return "empty stack" }
func (emptyStacker) Stack() []uintptr { return nil }

func TestStack(t *testing.T) {
	if s := errors.Stack(nil); s != nil {
		t.Errorf("got %v, want nil", s)
	}

	if s := errors.Stack(fmt.Errorf("foo")); s != nil {
		t.Errorf("got %v, want nil", s)
	}

	if s := errors.Stack(emptyStacker{}); len(s) != 0 {
		t.Errorf("got %v, want empty", s)
	}

	e := errors.New("foo")
	s := errors.Stack(e)
	if len(s) == 0 {
		t.Errorf("got %v, want non-empty", s)
	}

	if len(s) < 2 {
		t.Fatalf("got %d frames, want at least 2", len(s))
	}

	f0 := s[0]
	if f0.Function != "github.com/bobg/errors_test.TestStack" {
		t.Errorf("got %q, want %q", f0, "github.com/bobg/errors_test.TestStack")
	}
	if !strings.HasSuffix(f0.File, "/errors_test.go") {
		t.Errorf("got %q, want suffix %q", f0.File, "/errors_test.go")
	}
	if f0.Line == 0 {
		t.Error("got 0, want non-zero")
	}

	f1 := s[1]
	if f1.Function != "testing.tRunner" {
		t.Errorf("got %q, want %q", f1, "testing.tRunner")
	}
	if !strings.HasSuffix(f1.File, "/testing.go") {
		t.Errorf("got %q, want suffix %q", f1.File, "/testing.go")
	}

	ss := s.String()
	if ok, _ := regexp.MatchString(`/errors_test\.go:\d+: github\.com/bobg/errors_test\.TestStack`, ss); !ok {
		t.Error("got no match, want match")
	}
}

func TestFrameString(t *testing.T) {
	cases := []struct {
		frame errors.Frame
		want  string
	}{
		{
			frame: errors.Frame{File: "foo.go", Line: 10, Function: "pkg.Fn"},
			want:  "foo.go:10: pkg.Fn",
		},
		{
			frame: errors.Frame{File: "foo.go", Function: "pkg.Fn"},
			want:  "foo.go: pkg.Fn",
		},
		{
			frame: errors.Frame{Function: "pkg.Fn"},
			want:  "pkg.Fn",
		},
		{
			frame: errors.Frame{File: "foo.go", Line: 10},
			want:  "foo.go:10: (unknown function)",
		},
		{
			frame: errors.Frame{},
			want:  "(unknown)",
		},
	}

	for _, tc := range cases {
		if got := tc.frame.String(); got != tc.want {
			t.Errorf("Frame(%+v).String() = %q, want %q", tc.frame, got, tc.want)
		}
	}
}
