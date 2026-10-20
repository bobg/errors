package errors

// Unwrapper is an optional interface that errors may implement.
// The Unwrap method returns the nested error that this error wraps.
type Unwrapper interface {
	Unwrap() error
}

// Unwrapper is an optional interface that errors may implement.
// The Unwrap method returns the list of errors that this error wraps.
type MultiUnwrapper interface {
	Unwrap() []error
}

// Rewrapper is an optional interface that errors may implement.
// The Rewrap method wraps the given error with the same information that the original error wrapped its nested error.
type Rewrapper interface {
	Rewrap(error) error
}

// MultiRewrapper is an optional interface that errors may implement.
// The Rewrap method wraps the given errors with the same information that the original error wrapped its nested errors.
type MultiRewrapper interface {
	Rewrap(...error) error
}
