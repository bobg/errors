//go:build go1.20
// +build go1.20

package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestJoin(t *testing.T) {
	var (
		e1 = fmt.Errorf("error 1")
		e2 = fmt.Errorf("error 2")
	)

	cases := []struct {
		errs              []error
		isErr, isE1, isE2 bool
	}{{
		// empty case
	}, {
		errs: []error{nil, nil, nil},
	}, {
		errs:  []error{fmt.Errorf("not 1 or 2")},
		isErr: true,
	}, {
		errs:  []error{e1},
		isErr: true,
		isE1:  true,
	}, {
		errs:  []error{e2},
		isErr: true,
		isE2:  true,
	}, {
		errs:  []error{e1, e2},
		isErr: true,
		isE1:  true,
		isE2:  true,
	}}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("case_%02d", i+1), func(t *testing.T) {
			var got error
			for _, e := range tc.errs {
				got = Join(got, e)
			}
			if (got != nil) != tc.isErr {
				t.Errorf("got error %v, want error is %v", got, tc.isErr)
			}
			if errors.Is(got, e1) != tc.isE1 {
				t.Errorf("errors.Is(got, e1) is %v, want %v", errors.Is(got, e1), tc.isE1)
			}
			if errors.Is(got, e2) != tc.isE2 {
				t.Errorf("errors.Is(got, e2) is %v, want %v", errors.Is(got, e2), tc.isE2)
			}
		})
	}
}
