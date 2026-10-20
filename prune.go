package errors

import "slices"

// Prune returns a copy of the error tree rooted at tree,
// minus the smallest subtrees matching any of the targets.
//
// So, for example, Prune(Join(a, b, c), b) returns Join(a, c).
// Prune(Wrap(a, x), a) returns nil.
//
// Normally, Prune(Wrap(a, x), b) returns a
// (when a is not b),
// losing the wrapping info x.
// But if Wrap(a, x) is a [Rewrapper] or a [MultiRewrapper],
// then the wrapping info is preserved.
func Prune(tree error, targets ...error) error {
	if tree == nil {
		return nil
	}
	if len(targets) == 0 {
		return tree
	}
	if slices.Contains(targets, tree) {
		return nil
	}

	switch t := tree.(type) {
	case Unwrapper:
		p := Prune(t.Unwrap(), targets...)
		if p == nil {
			return nil
		}

		switch t := tree.(type) {
		case Rewrapper:
			return t.Rewrap(p)

		case MultiRewrapper:
			return t.Rewrap(p)

		default:
			return p
		}

	case MultiUnwrapper:
		var (
			uu     = t.Unwrap()
			pruned = make([]error, 0, len(uu))
		)

		for _, ue := range uu {
			if ue == nil {
				continue
			}
			p := Prune(ue, targets...)
			if p == nil {
				continue
			}
			pruned = append(pruned, p)
		}

		if len(pruned) == 0 {
			return nil
		}

		switch t := tree.(type) {
		case Rewrapper:
			return t.Rewrap(Join(pruned...))

		case MultiRewrapper:
			return t.Rewrap(pruned...)

		default:
			return Join(pruned...)
		}

	default:
		return tree
	}
}
