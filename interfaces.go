package errors

// Unwrapper is an optional interface that an error may implement.
// The Unwrap method returns the nested error that this error wraps.
type Unwrapper interface {
	Unwrap() error
}

// Unwrapper is an optional interface that an error may implement.
// The Unwrap method returns the list of errors that this error wraps.
type MultiUnwrapper interface {
	Unwrap() []error
}

// Rewrapper is an optional interface that an [Unwrapper] or [MultiUnwrapper] may implement.
// The Rewrap method wraps the given errors with the same information that the original error wrapped its nested error.
// A type that wraps only a single nested error may implement Rewrap by calling [Join] on the given error list.
type Rewrapper interface {
	Rewrap(...error) error
}
