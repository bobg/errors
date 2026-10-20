package errors

// Prune returns a copy of the error tree rooted at tree,
// minus the smallest subtrees matching any of the targets.
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
	switch t := tree.(type) {
	case Unwrapper:
		u := t.Unwrap()
		if isSmallestMatch(tree, targets, u) {
			return nil
		}
		if u == nil {
			return tree
		}
		u = Prune(u, targets...)
		if u == nil {
			return nil
		}

		switch t := tree.(type) {
		case Rewrapper:
			return t.Rewrap(u)

		case MultiRewrapper:
			return t.Rewrap(u)

		default:
			return u
		}

	case MultiUnwrapper:
		uu := t.Unwrap()
		if isSmallestMatch(tree, targets, uu...) {
			return nil
		}
		if len(uu) == 0 {
			return tree
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

		switch t := tree.(type) {
		case Rewrapper:
			return t.Rewrap(Join(newChildren...))

		case MultiRewrapper:
			return t.Rewrap(newChildren...)

		default:
			return Join(newChildren...)
		}

	default:
		if isSmallestMatch(tree, targets) {
			return nil
		}
		return tree
	}
}

func isSmallestMatch(tree error, targets []error, children ...error) bool {
	for _, target := range targets {
		if !Is(tree, target) {
			continue
		}
		childMatches := false
		for _, child := range children {
			if Is(child, target) {
				childMatches = true
				break
			}
		}
		if !childMatches {
			return true
		}
	}
	return false
}
