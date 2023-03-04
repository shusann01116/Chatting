package errors

type indexedError struct {
	cause error
	index int
}

// WithIndex attaches an array or slice index to the error.
// The data can be retreived using the Index() and Cause() method.
func WithIndex(err error, index int) error {
	return indexedError{cause: err, index: index}
}

func (e indexedError) Error() string {
	return e.cause.Error()
}

func (e indexedError) Cause() error {
	return e.cause
}

func (e indexedError) Index() int {
	return e.index
}
