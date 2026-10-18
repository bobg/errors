package errors

// Prune returns a copy of the error tree rooted at tree,
// minus the smallest subtrees matching (via [Is]) any of the targets.
//
// So, for example, Prune(Join(a, b, c), b) returns Join(a, c).
// Prune(Wrap(a, x), a) returns nil.
//
// Normally, Prune(Wrap(a, x), b) returns a
// (when !Is(a, b)),
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
	switch tree := tree.(type) {
	case Unwrapper:
		u := tree.Unwrap()
		if u == nil {
			return nil
		}
		u = Prune(u, targets...)
		if u == nil {
			return nil
		}

		switch tree := tree.(type) {
		case Rewrapper:
			return tree.Rewrap(u)

		case MultiRewrapper:
			return tree.Rewrap(u)

		default:
			return u
		}

	case MultiUnwrapper:
		uu := tree.Unwrap()
		if len(uu) == 0 {
			return nil
		}

		var newChildren []error
		for _, u := range uu {
			u = Prune(u, targets...)
			if u != nil {
				newChildren = append(newChildren, u)
			}
		}

		if len(newChildren) == 0 {
			return nil
		}

		switch tree := tree.(type) {
		case Rewrapper:
			return tree.Rewrap(Join(newChildren...))

		case MultiRewrapper:
			return tree.Rewrap(newChildren...)

		default:
			return Join(newChildren...)
		}

	default:
		for _, target := range targets {
			if Is(tree, target) {
				return nil
			}
		}
		return tree
	}
}
