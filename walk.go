package errors

// Walk walks the error tree of e,
// calling f on each error found.
// If f returns a non-nil error,
// Walk stops and returns that error,
// unless the error is [ErrSkip],
// in which case the walk skips e's children
// but otherwise continues without error.
//
// The error tree of e consists of e itself
// and any children obtained by calling e's optional Unwrap method
// (whose return type may be error or []error),
// recursively.
//
// The nodes of this tree are visited in a preorder, depth-first traversal.
func Walk(e error, f func(error) error) error {
	if e == nil {
		return nil
	}

	if err := f(e); err != nil {
		if Is(err, ErrSkip) {
			return nil
		}
		return err
	}

	if u := Unwrap(e); u != nil {
		return Walk(u, f)
	}

	if multi, ok := e.(multiUnwrap); ok {
		uu := multi.Unwrap()
		for _, ue := range uu {
			if err := Walk(ue, f); err != nil {
				return err
			}
		}
	}

	return nil
}

// ErrSkip can be used by the callback to [Walk] to skip an error's subtree without aborting the walk.
var ErrSkip = New("skip")

type singleUnwrap interface {
	Unwrap() error
}

type multiUnwrap interface {
	Unwrap() []error
}
