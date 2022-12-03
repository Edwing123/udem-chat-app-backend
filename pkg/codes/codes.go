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
