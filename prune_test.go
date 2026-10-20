package errors_test

import (
	"fmt"
	"testing"

	"github.com/bobg/errors"
)

func TestPrune(t *testing.T) {
	var (
		a = errors.New("a")
		b = errors.New("b")
		c = errors.New("c")
	)

	if got := errors.Prune(nil, a); got != nil {
		t.Errorf("got %v, want nil", got)
	}

	if got := errors.Prune(a); got != a {
		t.Errorf("got %v, want %v", got, a)
	}

	// Prune(Join(a, b, c), b) -> Join(a, c)
	j := errors.Join(a, b, c)
	got := errors.Prune(j, b)
	if got == nil || got.Error() != "a\nc" {
		t.Errorf("got %v, want %q", got, "a\nc")
	}
	if !errors.Is(got, a) || errors.Is(got, b) || !errors.Is(got, c) {
		t.Errorf("unexpected Is results on %v", got)
	}

	// Prune(Wrap(a, "x"), a) -> nil
	if got := errors.Prune(errors.Wrap(a, "x"), a); got != nil {
		t.Errorf("got %v, want nil", got)
	}

	// Prune(Wrap(a, "x"), b) preserves wrapping info via Rewrapper
	got = errors.Prune(errors.Wrap(a, "x"), b)
	if got == nil || got.Error() != "x: a" {
		t.Errorf("got %v, want %q", got, "x: a")
	}
	if !errors.Is(got, a) {
		t.Errorf("expected Is(got, a) to be true")
	}

	// Non-Rewrapper (fmt.Errorf) loses wrapping info when pruned
	got = errors.Prune(fmt.Errorf("x: %w", a), b)
	if got == nil || got.Error() != "a" {
		t.Errorf("got %v, want %q", got, "a")
	}

	// Wrapf with %w preserves wrapping info when a subtree of the wrapped error is pruned
	w := errors.Wrapf(errors.Join(a, b), "ctx %d (%w)", 42, c)
	got = errors.Prune(w, b)
	if got == nil || got.Error() != "ctx 42 (c): a" {
		t.Errorf("got %v, want %q", got, "ctx 42 (c): a")
	}
	if !errors.Is(got, a) || !errors.Is(got, c) || errors.Is(got, b) {
		t.Errorf("unexpected Is results on %v", got)
	}

	// Newf with %w preserves wrapping info when pruned
	nf := errors.Newf("ctx %d: %w", 42, errors.Join(a, b))
	got = errors.Prune(nf, a)
	if got == nil || got.Error() != "ctx 42: b" {
		t.Errorf("got %v, want %q", got, "ctx 42: b")
	}
	if !errors.Is(got, b) || errors.Is(got, a) {
		t.Errorf("unexpected Is results on %v", got)
	}

	// singleFormatted.Rewrap when child is unchanged
	sf := errors.Newf("ctx: %w", a)
	if got := errors.Prune(sf, b); got != sf {
		t.Errorf("got %v, want %v", got, sf)
	}

	// multiFormatted.Rewrap when number of children changes, and when unchanged
	mf := errors.Newf("ctx: %w, %w", a, b)
	if got := errors.Prune(mf, c); got != mf {
		t.Errorf("got %v, want %v", got, mf)
	}
	got = errors.Prune(mf, a)
	if got == nil || got.Error() != "b" {
		t.Errorf("got %v, want %q", got, "b")
	}
	if !errors.Is(got, b) || errors.Is(got, a) {
		t.Errorf("unexpected Is results on %v", got)
	}

	// wrapped.Rewrap and stacked.Rewrap when child changes, plus leaf match in default branch
	var (
		plainA = fmt.Errorf("a")
		plainB = fmt.Errorf("b")
	)
	ws := errors.Wrap(errors.Join(plainA, plainB), "x")
	got = errors.Prune(ws, plainB)
	if got == nil || got.Error() != "x: a" {
		t.Errorf("got %v, want %q", got, "x: a")
	}
	if !errors.Is(got, plainA) || errors.Is(got, plainB) {
		t.Errorf("unexpected Is results on %v", got)
	}

	// MultiUnwrapper matching target directly
	if got := errors.Prune(j, j); got != nil {
		t.Errorf("got %v, want nil", got)
	}

	// MultiUnwrapper with all children pruned
	if got := errors.Prune(errors.Join(a, b), a, b); got != nil {
		t.Errorf("got %v, want nil", got)
	}
}
