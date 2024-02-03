package errors

// Walk walks the error tree of e,
// calling f on each error found.
// If f returns a non-nil error, Walk stops and returns that error.
//
// The error tree of e consists of e itself
// followed by any error or sequence of errors
// obtained by calling e's Unwrap method
// (whose return type may be error or []error),
// recursively.
func Walk(e error, f func(error) error) error {
	if e == nil {
		return nil
	}

	if err := f(e); err != nil {
		return err
	}

	if u := Unwrap(e); u != nil {
		return Walk(u, f)
	}

	if multi, ok := e.(interface{ Unwrap() []error }); ok {
		uu := multi.Unwrap()
		for _, ue := range uu {
			if err := Walk(ue, f); err != nil {
				return err
			}
		}
	}

	return nil
}
