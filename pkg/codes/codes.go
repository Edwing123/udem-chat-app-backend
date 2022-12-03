package codes

func NewCode(description string) Code {
	return Code(description)
}

// Represents an error code.
type Code string

// It implements the error interface.
func (c Code) Error() string {
	return string(c)
}

var (
	// Server related.
	ErrServerInternal = NewCode("server_internal")

	// Client related.
	ErrCannotDecodeJSON = NewCode(("cannot_decode_json"))
	ErrAuthRequired     = NewCode("auth_required")

	// Validation related.
	ErrPasswordNotValid = NewCode("password_not_valid")
)
