package errors

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWalk(t *testing.T) {
	if err := Walk(nil, func(error) error { return nil }); err != nil {
		t.Errorf("Walk(nil) = %v, want nil", err)
	}

	testWalkHelper(t, ErrSkip)
	testWalkHelper(t, New("foo"))
	testWalkHelper(t, nil)
}

func testWalkHelper(t *testing.T, errToThrow error) {
	var (
		e1 = New("1")
		e2 = New("2")

		e1a = Unwrap(e1)
		e2a = Unwrap(e2)

		e3 = Join(e1, e2)
		e4 = Wrap(e3, "4")
		e5 = Wrap(e4, "5")

		got []int
	)

	err := Walk(e5, func(e error) error {
		switch e {
		case e1:
			if errToThrow != nil {
				return errToThrow
			}
			got = append(got, 1)
		case e2:
			got = append(got, 2)
		case e3:
			got = append(got, 3)
		case e4:
			got = append(got, 4)
		case e5:
			got = append(got, 5)
		case e1a:
			got = append(got, -1)
		case e2a:
			got = append(got, -2)
		default:
			return fmt.Errorf("unexpected error %v", e)
		}
		return nil
	})

	var want []int
	switch {
	case Is(errToThrow, ErrSkip):
		if err != nil {
			t.Fatal(err)
		}
		want = []int{5, 4, 3, 2, -2}

	case errToThrow != nil:
		if err == nil {
			t.Error("did not get expected error")
		}
		want = []int{5, 4, 3}

	default:
		if err != nil {
			t.Fatal(err)
		}
		want = []int{5, 4, 3, 1, -1, 2, -2}
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}
