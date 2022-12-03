package hashing

import (
	"github.com/Edwing123/udem-chat-app/pkg/codes"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordHashingFail = codes.NewCode("password_hashing_fail")
)

// HashPassword creates the bcrypt hash of the password.
// On success it returns the hashed password and a nil error, and on failure,
// it returns a nil hashed password and non-nil error.
func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, 12)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

// VerifyPassword compares a bcrypt hashed password with its
// possible plaintext equivalent. Returns true on success, or false on failure.
func VerifyPassword(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(
		hashedPassword,
		password,
	)

	return err == nil
}
