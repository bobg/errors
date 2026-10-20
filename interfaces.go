package errors

// Unwrapper is an optional interface that an error may implement.
// The Unwrap method returns the nested error that this error wraps.
type Unwrapper interface {
	error
	Unwrap() error
}

// Unwrapper is an optional interface that an error may implement.
// The Unwrap method returns the list of errors that this error wraps.
type MultiUnwrapper interface {
	error
	Unwrap() []error
}

// Rewrapper is an optional interface that an [Unwrapper] may implement.
// The Rewrap method wraps the given error with the same information that the original error wrapped its nested error.
type Rewrapper interface {
	Unwrapper
	Rewrap(error) error
}

// MultiRewrapper is an optional interface that a [MultiUnwrapper] may implement.
// The Rewrap method wraps the given errors with the same information that the original error wrapped its nested errors.
type MultiRewrapper interface {
	MultiUnwrapper
	Rewrap(...error) error
}
