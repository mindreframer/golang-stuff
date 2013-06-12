package hood

const (
	ValidationErrorValueNotSet = (1<<16 + iota)
	ValidationErrorValueTooSmall
	ValidationErrorValueTooBig
	ValidationErrorValueTooShort
	ValidationErrorValueTooLong
	ValidationErrorValueNotMatch
)

// Validation error type
type ValidationError struct {
	kind  int
	field string
}

// NewValidationError returns a new validation error with the specified id and
// text. The id's purpose is to distinguish different validation error types.
// Built-in validation error ids start at 65536, so you should keep your custom
// ids under that value.
func NewValidationError(id int, field string) error {
	return &ValidationError{id, field}
}

func (e *ValidationError) Error() string {
	kindStr := ""
	switch e.kind {
	case ValidationErrorValueNotSet:
		kindStr = " not set"
	case ValidationErrorValueTooBig:
		kindStr = " too big"
	case ValidationErrorValueTooLong:
		kindStr = " too long"
	case ValidationErrorValueTooSmall:
		kindStr = " too small"
	case ValidationErrorValueTooShort:
		kindStr = " too short"
	case ValidationErrorValueNotMatch:
		kindStr = " not match"
	}
	return e.field + kindStr
}

func (e *ValidationError) Kind() int {
	return e.kind
}

func (e *ValidationError) Field() string {
	return e.field
}
